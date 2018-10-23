package main

import (
	"flag"
	"fmt"
	"os"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	EvictionKind = "Eviction"
)

var (
	namespace = ""
	verbose   = false
)

func main() {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	currentNamespace, _, err := kubeConfig.Namespace()
	// if err != nil {
	// 	panic(err.Error())
	// }

	flag.StringVar(&namespace, "namespace", currentNamespace, "namespace of the pod")
	flag.BoolVar(&verbose, "verbose", false, "show more details")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("USAGE: kubectl evict [--namespace NAMESPACE] POD_NAME")
		os.Exit(1)
	}
	podName := flag.Args()[0]

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		panic(err.Error())
	}

	if verbose {
		fmt.Println("Pod Name:", podName)
		fmt.Println("Namespace:", namespace)
		fmt.Println("Context:", rawConfig.CurrentContext)
	}

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	policyGroupVersion := "v1beta1"
	deleteOptions := &metav1.DeleteOptions{}

	eviction := &policyv1beta1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		DeleteOptions: deleteOptions,
	}
	err = clientset.PolicyV1beta1().Evictions(eviction.Namespace).Evict(eviction)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if verbose {
		fmt.Println("Done!")
	}

	os.Exit(0)
}
