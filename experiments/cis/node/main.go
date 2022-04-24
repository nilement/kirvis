package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nilement/node/config"
	"github.com/nilement/node/experiment"
)

type Runner struct {
	log *logrus.Entry
}

func main() {
	log := logrus.WithFields(logrus.Fields{})
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

	completedExperiments := make([]experiment.Experiment, 0)
	for _, e := range experiments {
		err := e.Execute()
		if err != nil {
			log.Printf("error while executing: %s, %v", e.Key, err)
			for _, ce := range completedExperiments {
				err := ce.RestoreFile()
				if err != nil {
					log.Fatal("failed rollback!")
				}
				log.Infof("Completed rollback for: %s", ce.Key)
			}
			log.Fatalf("rollbacked %s", e.Key)
		}
		completedExperiments = append(completedExperiments, e)
	}

	r := &Runner{
		log: log,
	}

	err = r.wait(completedExperiments)
	if err != nil {
		log.Fatalf("failed restore: %v", err)
	}
}

func(r *Runner) rollbackExperiments(experiments []experiment.Experiment) error {
	for _, ce := range experiments {
		err := ce.RestoreFile()
		if err != nil {
			r.log.Errorf("Failed rollback for: %s", ce.Key)
			return err
		}
		r.log.Infof("Completed rollback for: %s", ce.Key)
	}

	return nil
}

func(r *Runner) wait(completedExperiments []experiment.Experiment) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	keys := make([]string, len(completedExperiments))
	for i, e := range completedExperiments {
		keys[i] = e.Key
	}
	misconf := strings.Join(keys, ", ")

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <- t.C:
			r.log.Infof("Active CIS Benchmark Worker node misconfigurations: %s", misconf)
		case <-sigs:
			r.log.Info("Experiments terminated. Rolling back")
			return r.rollbackExperiments(completedExperiments)
		}
	}
}
