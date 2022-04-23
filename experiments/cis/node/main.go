package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/nilement/komrade/config"
	"github.com/nilement/komrade/experiment"
)

func main() {
	if runtime.GOOS != "linux" {
		log.Fatalf("only Linux runtime is supported!")
	}

	if len(os.Args[1:]) == 0 {
		log.Fatalf("no experiments specified!")
	}

	cfg, err := config.ReadConfig("experiments.yaml")
	if err != nil {
		log.Fatal("error reading config: %w", err)
	}

	if len(cfg.Experiments) == 0 {
		log.Fatal("no experiments available!")
	}

	experiments := make([]experiment.Experiment, 0)
	for _, arg := range os.Args[1:] {
		exp, ok := cfg.ExperimentMap[arg]
		if !ok {
			log.Fatalf("Specified experiment key: %s is not supported", arg)
		}
		experiments = append(experiments, exp)
	}

	for _, e := range experiments {
		err := e.Backup()
		if err != nil {
			log.Fatal(err)
		}
	}

	completedExperiments := make([]experiment.Experiment, len(experiments))
	for _, e := range experiments {
		err := e.Execute()
		if err != nil {
			log.Printf("error while executing: %s, %v", e.Key, err)
			for _, ce := range completedExperiments {
				err := ce.RestoreFile()
				if err != nil {
					log.Fatal("failed rollback!")
				}
			}
			log.Fatalf("rollbacked %s", e.Key)
		}
		completedExperiments = append(completedExperiments, e)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		log.Println("Applied CIS Benchmark Worker node misconfigurations")
		time.Sleep(time.Second * 30)
		select {
		case <-sigs:
			log.Println("Experiments terminated. Rolling back")
			for _, ce := range completedExperiments {
				err := ce.RestoreFile()
				if err != nil {
					log.Fatal("Failed rollback!")
				}
			}
		}
	}
}
