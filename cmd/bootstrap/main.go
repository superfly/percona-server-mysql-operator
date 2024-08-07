package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

const (
	fullClusterCrashFile = "/var/lib/mysql/full-cluster-crash"
	manualRecoveryFile   = "/var/lib/mysql/sleep-forever"
	bootstrapFile        = "/var/lib/mysql/startup_bootstrap.lock"
)

func main() {

	// f, err := os.OpenFile(filepath.Join(mysql.DataMountPath, "bootstrap.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// FKS: Keep logs in stdout so we can understand why the container crashes
	// log.SetOutput(f)

	var recovering bool

	fullClusterCrash, err := fileExists(fullClusterCrashFile)
	if err == nil && fullClusterCrash {
		log.Printf("%s exists. Skipping bootstrap.", fullClusterCrashFile)
		recovering = true
	}

	manualRecovery, err := fileExists(manualRecoveryFile)
	if err == nil && manualRecovery {
		log.Printf("%s exists. Skipping bootstrap.", manualRecoveryFile)
		recovering = true
	}

	if !recovering {

		log.Printf("bootstrap sleeping for 15 seconds to allow the mysql server to start, SRV records to populate and other members to boot")
		time.Sleep(15 * time.Second)

		os.Remove(bootstrapFile)

		log.Printf("Starting bootstrap")
		clusterType := os.Getenv("CLUSTER_TYPE")
		switch clusterType {
		case "group-replication":
			if err := bootstrapGroupReplication(context.Background()); err != nil {
				log.Fatalf("bootstrap: bootstrap failed: %v", err)
			}
		case "async":
			if err := bootstrapAsyncReplication(context.Background()); err != nil {
				log.Fatalf("bootstrap failed: %v", err)
			}
			touchBootstrap()
		default:
			log.Fatalf("bootstrap invalid cluster type: %v", clusterType)
		}
	}
}

func touchBootstrap() {
	f, err := os.Create(bootstrapFile)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer f.Close()
}

func restartContainer() {
	log.Printf("bootstrap attempting to restart container by killing healthcheck_http process")

	processes, err := process.Processes()
	if err != nil {
		fmt.Printf("bootstrap error listing processes: %v\n", err)
		return
	}

	for _, p := range processes {
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}

		if strings.Contains(cmdline, "healthcheck_http") {
			log.Printf("bootstrap killing process %d: %s\n", p.Pid, cmdline)

			// Kill the process
			err = p.Kill()
			if err != nil {
				log.Printf("bootstrap error killing process %d: %v\n", p.Pid, err)
			}
		}
	}
}
