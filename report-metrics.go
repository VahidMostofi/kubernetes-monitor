package main

// https://medium.com/swlh/clientset-module-for-in-cluster-and-out-cluster-3f0d80af79ed

import (
	"fmt"
	"encoding/json"
	"context"
	discovery "github.com/gkarthiks/k8s-discovery"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	k8s *discovery.K8s
)

func main() {
	k8s, _ = discovery.NewK8s()
	namespace, _ := k8s.GetNamespace()
	version, _ := k8s.GetVersion()
	fmt.Printf("Specified Namespace: %s\n", namespace)
	fmt.Printf("K8s version: %s\n", version)
	podMetrics, err := k8s.MetricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, lm := range podMetrics.Items{
		
		b,_ := json.Marshal(lm)
		fmt.Println(string(b))
		for _, c := range lm.Containers{
			fmt.Println("cpu",c.Usage.Cpu().ToDec())
		}
		break
	}
}