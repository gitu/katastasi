package core

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

type monitor struct {
	client *kubernetes.Clientset
}

func newMonitor(client *kubernetes.Clientset) *monitor {
	return &monitor{
		client: client,
	}
}

func (m *monitor) startMonitor() {
	fmt.Println("Starting monitor...")
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		return m.client.CoreV1().ConfigMaps("*").Watch(context.Background(), metav1.ListOptions{LabelSelector: "katastasi.io/statuspage-queries=true", TimeoutSeconds: &timeOut})
	}

	watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})

	go func() {
		c := watcher.ResultChan()
		for {
			e := <-c
			cm := e.Object.(*corev1.ConfigMap)

			println("configmap: " + cm.Name)
		}
	}()

	list, err := m.client.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{LabelSelector: "katastasi.io/statuspage-queries=true"})
	if err != nil {
		panic(err.Error())
	}
	for _, cm := range list.Items {
		println("configmap: " + cm.Name)
	}
}
