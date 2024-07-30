package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func commandHandler(arg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("/opt/percona/healthcheck", arg)
		err := cmd.Run()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Command failed: %v\n", err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Command executed successfully")
		}
	}
}

func main() {
	pathsAndArgs := []string{"readiness", "liveness"}

	for _, arg := range pathsAndArgs {
		http.HandleFunc("/"+arg, commandHandler(arg))
	}

	log.Println("Starting server on port 5500")
	log.Fatal(http.ListenAndServe(":5500", nil))
}
