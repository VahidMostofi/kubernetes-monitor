package main

// https://medium.com/swlh/clientset-module-for-in-cluster-and-out-cluster-3f0d80af79ed

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	discovery "github.com/gkarthiks/k8s-discovery"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	k8s *discovery.K8s
)

type Usage struct {
	PodName string
	CPU     int64
	Memory  int64
}

var measurementName string

func report(u Usage) {
	//format: https://docs.influxdata.com/influxdb/v2.0/reference/syntax/line-protocol/
	fmt.Println(measurementName + "," + "podName=" + u.PodName + " " + "cpu=" + strconv.FormatInt(u.CPU, 10) + "," + "memory=" + strconv.FormatInt(u.Memory, 10))
}

func main() {
	k8s, err := discovery.NewK8s()
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("unable to create k8s object: %w\n", err))
	}
	namespace := os.Getenv("METRICS_MONITOR_K8S_NAMESPACE")
	interval := os.Getenv("METRICS_MONITOR_INTERVAL") // interval in seconds
	if len(interval) == 0 {
		interval = "10"
	}
	intervalValue, err := strconv.Atoi(interval)
	if err != nil {
		panic(fmt.Sprintf("cant parse interval value: %s, %w", interval, err))
	}
	measurementName = os.Getenv("METRICS_MONITOR_MEASUREMENT_NAME")
	if len(measurementName) == 0 {
		measurementName = "resource_usage"
	}
	// version, _ := k8s.GetVersion()
	// fmt.Printf("Specified Namespace: %s\n", namespace)
	// fmt.Printf("K8s version: %s\n", version)
	for {
		podMetrics, err := k8s.MetricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, pdItem := range podMetrics.Items {
			parts := strings.Split(pdItem.GetObjectMeta().GetSelfLink(), "/")
			podName := parts[len(parts)-1]
			u := Usage{
				podName,
				0,
				0,
			}
			for _, c := range pdItem.Containers {
				u.CPU += c.Usage.Cpu().MilliValue()
				u.Memory += c.Usage.Memory().MilliValue()
			}
			report(u)
		}
		time.Sleep(time.Duration(intervalValue) * time.Second)
	}
}
