package main

import (
	"context"
	"log"
	"os"
	"time"
)

const (
	fullClusterCrashFile = "/var/lib/mysql/full-cluster-crash"
	manualRecoveryFile   = "/var/lib/mysql/sleep-forever"
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

		log.Printf("Sleeping for 30 seconds to allow the mysql server to start and SRV records to populate")
		time.Sleep(30 * time.Second)

		log.Printf("Starting bootstrap")
		clusterType := os.Getenv("CLUSTER_TYPE")
		switch clusterType {
		case "group-replication":
			if err := bootstrapGroupReplication(context.Background()); err != nil {
				log.Fatalf("bootstrap failed: %v", err)
			}
		case "async":
			if err := bootstrapAsyncReplication(context.Background()); err != nil {
				log.Fatalf("bootstrap failed: %v", err)
			}
		default:
			log.Fatalf("Invalid cluster type: %v", clusterType)
		}
	}

	// FKS: Since we run bootstrap as a sidecar instead of a startup probe,
	//      keep the bootstrap container running to avoid the pod being restarted.
	for {
		log.Printf("Bootstrap sleeping for 1 hour")
		time.Sleep(1 * time.Hour)
	}
}
