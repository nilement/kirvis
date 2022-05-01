package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"

	"github.com/nilement/kubelet/config"
	"github.com/nilement/kubelet/experiment"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	//configURLTemplate = "https://10.96.0.1:443/api/v1/nodes/%s/proxy/configz"
	configURLTemplate = "http://localhost:8001/api/v1/nodes/%s/proxy/configz"
	configFile = "./k8s/experiments.yaml"
)

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

	restconfig, err := kubeConfigFromEnv(cfg)
	if err != nil {
		log.Fatal("kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		log.Fatal("error initing clienset: %w", err)
	}

	resp, err := http.Get(fmt.Sprintf(configURLTemplate, cfg.NodeName))
	if err != nil {
		log.Fatal("could not retrieve old configmap: %w", err)
	}

	// read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed reading resp to bytes: %w", err)
	}

	var kubeletConfig map[string]interface{}
	err = json.Unmarshal(jsonDataFromHttp, &kubeletConfig)
	if err != nil {
		log.Fatal("failed reading resp to bytes: %w", err)
	}

	//   "kind": "KubeletConfiguration",
	//  "apiVersion": "kubelet.config.k8s.io/v1beta1"

	newCM := kubeletConfig["kubeletconfig"].(map[string]interface{})

	newCM["kind"] = "KubeletConfiguration"
	newCM["apiVersion"] = "kubelet.config.k8s.io/v1beta1"
	newCM["eventRecordQPS"] = 10

	cmbytes, err := json.Marshal(newCM)
	//
	//var kconfig kubelet.KubeletConfiguration
	//err = json.Unmarshal(newCM, &kconfig)
	//if err != nil {
	//	log.Fatal("kubelet configuration parse: %w", err)
	//}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kube-system",
			Name: "chaos-kubelet-config",
		},
		Data: map[string]string{
			"kubelet": string(cmbytes),
		},
	}

	ctx := context.Background()

	patchCM, err := clientset.CoreV1().ConfigMaps("kube-system").Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("error patching configmap: %w", err)
	}

	selfNode, err := clientset.CoreV1().Nodes().Get(ctx, cfg.NodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatal("error getting node: %w", err)
	}

	selfNode.Spec.ConfigSource = createConfigSource(patchCM.Name)

	_, err = clientset.CoreV1().Nodes().Update(ctx, selfNode, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal("error patching node: %w", err)
	}
}

func createConfigSource(configMapName string) *corev1.NodeConfigSource {
	return &corev1.NodeConfigSource{
		ConfigMap: &corev1.ConfigMapNodeConfigSource{
			Name: configMapName,
			Namespace: "kube-system",
			KubeletConfigKey: "kubelet",
		},
	}
}

func kubeConfigFromEnv(cfg *config.Config) (*rest.Config, error) {
	kubepath := cfg.KubeconfigPath
	if kubepath == "" {
		return nil, nil
	}

	data, err := ioutil.ReadFile(kubepath)
	if err != nil {
		return nil, fmt.Errorf("reading kubeconfig at %s: %w", kubepath, err)
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return nil, fmt.Errorf("building rest config from kubeconfig at %s: %w", kubepath, err)
	}

	return restConfig, nil
}
