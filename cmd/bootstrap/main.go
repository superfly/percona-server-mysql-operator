package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/percona/percona-server-mysql-operator/pkg/mysql"
)

const (
	fullClusterCrashFile = "/var/lib/mysql/full-cluster-crash"
	manualRecoveryFile   = "/var/lib/mysql/sleep-forever"
)

func main() {
	http.HandleFunc("/startup", startupProbeHandler)

	log.Println("Starting HTTP server on port 8091")
	if err := http.ListenAndServe(":8091", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startupProbeHandler(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile(filepath.Join(mysql.DataMountPath, "bootstrap.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening file: %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	fullClusterCrash, err := fileExists(fullClusterCrashFile)
	if err == nil && fullClusterCrash {
		log.Printf("%s exists. exiting...", fullClusterCrashFile)
		http.Error(w, fmt.Sprintf("%s exists. exiting...", fullClusterCrashFile), http.StatusInternalServerError)
		return
	}

	manualRecovery, err := fileExists(manualRecoveryFile)
	if err == nil && manualRecovery {
		log.Printf("%s exists. exiting...", manualRecoveryFile)
		http.Error(w, fmt.Sprintf("%s exists. exiting...", manualRecoveryFile), http.StatusInternalServerError)
		return
	}

	clusterType := os.Getenv("CLUSTER_TYPE")
	switch clusterType {
	case "group-replication":
		if err := bootstrapGroupReplication(context.Background()); err != nil {
			log.Printf("bootstrap failed: %v", err)
			http.Error(w, fmt.Sprintf("bootstrap failed: %v", err), http.StatusInternalServerError)
			return
		}
	case "async":
		if err := bootstrapAsyncReplication(context.Background()); err != nil {
			log.Printf("bootstrap failed: %v", err)
			http.Error(w, fmt.Sprintf("bootstrap failed: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		errMsg := fmt.Sprintf("Invalid cluster type: %v", clusterType)
		log.Print(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Startup probe successful"))
}
