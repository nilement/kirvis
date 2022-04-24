package main

import (
	"fmt"
	"github.com/nilement/apiserver/backups"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"

	"github.com/nilement/apiserver/config"
	"github.com/nilement/apiserver/experiment"
)

const (
	apiServerFile = "/manifests/kube-apiserver.yaml"
	//apiServerFile = "./k8s/kube-apiserver.yaml"
	configFile = "experiments.yaml"
	//configFile = "./k8s/experiments.yaml"
	apiserverProcess = "kube-apiserver"
)

type Runner struct {
	log *logrus.Entry
}

func main() {
	log := logrus.WithFields(logrus.Fields{})
	log.Info("starting")
	if len(os.Args[1:]) == 0 {
		log.Fatalf("no experiments specified!")
	}

	cfg, err := config.ReadConfig(configFile)
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

	log.Infof("reading apiserver: %s", apiServerFile)
	apiserver, err := config.ReadAPIServer(apiServerFile)
	if err != nil {
		log.Fatalf("error reading api-server: %v", err)
	}

	log.Infof("backing up apiserver: %s", apiServerFile)
	err = backups.BackupFile(apiServerFile)
	if err != nil {
		log.Fatalf("error backing up api-server: %v", err)
	}

	commands := apiserver.Spec.Containers[0].Command
	log.Infof("found commands: %s", commands[0])
	for _, e := range experiments {
		commands, err = e.Execute(commands)
		if err != nil {
			log.Errorf("rolling back due to error: %v", err)
			err = backups.RestoreFile(apiServerFile)
			if err != nil {
				log.Fatalf("rollback failed: %v", err)
			}
			return
		}
	}

	apiserver.Spec.Containers[0].Command = commands
	err = config.WriteAPIServer(apiServerFile, apiserver)
	if err != nil {
		log.Fatalf("failed to output kube-apiserver.yaml: %v", err)
	}

	//args, err := getKubeAPIServerArgs()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, e := range experiments {
	//	applied := e.CheckIfApplied(args)
	//	if !applied {
	//		log.Errorf("Experiment %s not applied, rolling back", e.Key)
	//		err = backups.RestoreFile(apiServerFile)
	//		if err != nil {
	//			log.Fatalf("rollback failed: %v", err)
	//		}
	//		return
	//	}
	//}
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
			r.log.Infof("Active API server misconfigurations %s", misconf)
		case <-sigs:
			r.log.Info("Experiments terminated. Rolling back")
		return r.rollbackExperiments(completedExperiments)
		}
	}
}

func getKubeAPIServerArgs() ([]string, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range procs {
		name, _ := p.Name()
		if name == apiserverProcess {
			cmds, _ := p.CmdlineSlice()
			return cmds, nil
		}
	}
	return nil, fmt.Errorf("kube-apiserver process not found")
}