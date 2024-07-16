package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	apiv1alpha1 "github.com/percona/percona-server-mysql-operator/api/v1alpha1"
	database "github.com/percona/percona-server-mysql-operator/cmd/db"
	mysqldb "github.com/percona/percona-server-mysql-operator/pkg/db"
	"github.com/percona/percona-server-mysql-operator/pkg/k8s"
	"github.com/percona/percona-server-mysql-operator/pkg/mysql"
)

func main() {
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/liveness", livenessHandler)
	http.HandleFunc("/replication", replicationHandler)

	log.Println("Starting HTTP server on port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch os.Getenv("CLUSTER_TYPE") {
	case "async":
		if err := checkReadinessAsync(ctx); err != nil {
			http.Error(w, fmt.Sprintf("readiness check failed: %v", err), http.StatusInternalServerError)
			return
		}
	case "group-replication":
		if err := checkReadinessGR(ctx); err != nil {
			http.Error(w, fmt.Sprintf("readiness check failed: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid or missing CLUSTER_TYPE environment variable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Readiness check passed"))
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch os.Getenv("CLUSTER_TYPE") {
	case "async":
		if err := checkLivenessAsync(ctx); err != nil {
			http.Error(w, fmt.Sprintf("liveness check failed: %v", err), http.StatusInternalServerError)
			return
		}
	case "group-replication":
		if err := checkLivenessGR(ctx); err != nil {
			http.Error(w, fmt.Sprintf("liveness check failed: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid or missing CLUSTER_TYPE environment variable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Liveness check passed"))
}

func replicationHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := checkReplication(ctx); err != nil {
		http.Error(w, fmt.Sprintf("replication check failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Replication check passed"))
}

func checkReadinessAsync(ctx context.Context) error {
	podIP, err := getPodIP()
	if err != nil {
		return errors.Wrap(err, "get pod IP")
	}

	monitorPass, err := getSecret(string(apiv1alpha1.UserMonitor))
	if err != nil {
		return errors.Wrapf(err, "get %s password", apiv1alpha1.UserMonitor)
	}

	db, err := database.NewDatabase(ctx, apiv1alpha1.UserMonitor, monitorPass, podIP, mysql.DefaultAdminPort)
	if err != nil {
		return errors.Wrap(err, "connect to db")
	}
	defer db.Close()

	readOnly, err := db.IsReadonly(ctx)
	if err != nil {
		return errors.Wrap(err, "check read only status")
	}

	// if isReplica is true, replication is active
	isReplica, err := db.IsReplica(ctx)
	if err != nil {
		return errors.Wrap(err, "check replica status")
	}

	if isReplica && !readOnly {
		return errors.New("replica is not read only")
	}

	return nil
}

func checkReadinessGR(ctx context.Context) error {
	podIP, err := getPodIP()
	if err != nil {
		return errors.Wrap(err, "get pod IP")
	}

	monitorPass, err := getSecret(string(apiv1alpha1.UserMonitor))
	if err != nil {
		return errors.Wrapf(err, "get %s password", apiv1alpha1.UserMonitor)
	}

	db, err := database.NewDatabase(ctx, apiv1alpha1.UserMonitor, monitorPass, podIP, mysql.DefaultAdminPort)
	if err != nil {
		return errors.Wrap(err, "connect to db")
	}
	defer db.Close()

	fqdn, err := getPodFQDN(os.Getenv("SERVICE_NAME"))
	if err != nil {
		return errors.Wrap(err, "get pod hostname")
	}

	state, err := db.GetMemberState(ctx, fqdn)
	if err != nil {
		return errors.Wrap(err, "get member state")
	}

	if state != mysqldb.MemberStateOnline {
		return errors.Errorf("Member state: %s", state)
	}

	return nil
}

func checkLivenessAsync(ctx context.Context) error {
	podIP, err := getPodIP()
	if err != nil {
		return errors.Wrap(err, "get pod IP")
	}

	monitorPass, err := getSecret(string(apiv1alpha1.UserMonitor))
	if err != nil {
		return errors.Wrapf(err, "get %s password", apiv1alpha1.UserMonitor)
	}

	db, err := database.NewDatabase(ctx, apiv1alpha1.UserMonitor, monitorPass, podIP, mysql.DefaultAdminPort)
	if err != nil {
		return errors.Wrap(err, "connect to db")
	}
	defer db.Close()

	return db.DumbQuery(ctx)
}

func checkLivenessGR(ctx context.Context) error {
	podIP, err := getPodIP()
	if err != nil {
		return errors.Wrap(err, "get pod IP")
	}

	monitorPass, err := getSecret(string(apiv1alpha1.UserMonitor))
	if err != nil {
		return errors.Wrapf(err, "get %s password", apiv1alpha1.UserMonitor)
	}

	db, err := database.NewDatabase(ctx, apiv1alpha1.UserMonitor, monitorPass, podIP, mysql.DefaultAdminPort)
	if err != nil {
		return errors.Wrap(err, "connect to db")
	}
	defer db.Close()

	in, err := db.CheckIfInPrimaryPartition(ctx)
	if err != nil {
		return errors.Wrap(err, "check if member in primary partition")
	}

	log.Printf("in primary partition: %t", in)

	if !in {
		return errors.New("possible split brain!")
	}

	return nil
}

func checkReplication(ctx context.Context) error {
	podIP, err := getPodIP()
	if err != nil {
		return errors.Wrap(err, "get pod IP")
	}

	monitorPass, err := getSecret(string(apiv1alpha1.UserMonitor))
	if err != nil {
		return errors.Wrapf(err, "get %s password", apiv1alpha1.UserMonitor)
	}

	db, err := database.NewDatabase(ctx, apiv1alpha1.UserMonitor, monitorPass, podIP, mysql.DefaultAdminPort)
	if err != nil {
		return errors.Wrap(err, "connect to db")
	}
	defer db.Close()

	// if isReplica is true, replication is active
	isReplica, err := db.IsReplica(ctx)
	if err != nil {
		return errors.Wrap(err, "check replica status")
	}

	if !isReplica {
		return errors.New("replication is not active")
	}

	return nil
}

func getSecret(username string) (string, error) {
	path := filepath.Join(mysql.CredsMountPath, username)
	sBytes, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "read %s", path)
	}

	return strings.TrimSpace(string(sBytes)), nil
}

func getPodHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", errors.Wrap(err, "get hostname")
	}

	return hostname, nil
}

func getPodIP() (string, error) {
	hostname, err := getPodHostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return "", errors.Wrapf(err, "lookup %s", hostname)
	}

	return addrs[1], nil
}

func getPodFQDN(svcName string) (string, error) {
	hostname, err := getPodHostname()
	if err != nil {
		return "", err
	}

	namespace, err := k8s.DefaultAPINamespace()
	if err != nil {
		return "", errors.Wrap(err, "get namespace")
	}

	return fmt.Sprintf("%s.%s.%s", hostname, svcName, namespace), nil
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "os stat")
	}
	return true, nil
}
