package core

import (
	"context"
	"github.com/gitu/katastasi/pkg/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	toolsWatch "k8s.io/client-go/tools/watch"
	"k8s.io/client-go/util/homedir"
	"log/slog"
	"path/filepath"
	"strings"
)

type kubeWatcher struct {
	c        *config.Config
	watchers []*toolsWatch.RetryWatcher
}

func newKubeWatcher(c *config.Config) *kubeWatcher {
	return &kubeWatcher{
		c: c,
	}
}

func (kw *kubeWatcher) start() error {
	clientSet, err := kubernetes.NewForConfig(kubeClientConfig())
	if err != nil {
		return err
	}

	for _, namespace := range viper.GetStringSlice("autoload.namespaces.services") {
		watcher, err := watchConfigMaps(context.Background(), clientSet, namespace, "katastasi.io=service", kw.addServiceToConfig)
		if err == nil {
			kw.watchers = append(kw.watchers, watcher)
		}

	}
	for _, namespace := range viper.GetStringSlice("autoload.namespaces.pages") {
		watcher, err := watchConfigMaps(context.Background(), clientSet, namespace, "katastasi.io=page", kw.addStatusPagesToConfig)
		if err == nil {
			kw.watchers = append(kw.watchers, watcher)
		}
	}
	return nil
}

func (kw *kubeWatcher) stop() {
	for _, watcher := range kw.watchers {
		watcher.Stop()
		<-watcher.Done()
	}
}

func watchConfigMaps(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, labelSelector string, configMapHandler func(event watch.EventType, configMap *corev1.ConfigMap)) (*toolsWatch.RetryWatcher, error) {

	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		return clientSet.CoreV1().ConfigMaps(namespace).Watch(ctx, metav1.ListOptions{
			TimeoutSeconds: &timeOut,
			LabelSelector:  labelSelector,
		})
	}
	watcher, err := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	if err != nil {
		slog.Error("Error creating watcher for ConfigMaps", "namespace", namespace, "error", err.Error())
		return watcher, err
	}
	go func() {
		for event := range watcher.ResultChan() {
			configMap := event.Object.(*corev1.ConfigMap)
			configMapHandler(event.Type, configMap)
		}
		slog.Debug("Watcher for ConfigMaps stopped", "namespace", namespace)
	}()
	return watcher, nil
}

func (kw *kubeWatcher) addServiceToConfig(event watch.EventType, configMap *corev1.ConfigMap) {
	annotations := configMap.GetAnnotations()
	service := config.Service{
		ID:          annotations["katastasi.io/service-id"],
		Name:        annotations["katastasi.io/name"],
		Contact:     annotations["katastasi.io/contact"],
		Owner:       annotations["katastasi.io/owner"],
		URL:         annotations["katastasi.io/url"],
		Environment: annotations["katastasi.io/environment"],
	}
	if service.Name == "" {
		service.Name = configMap.Name
	}
	if service.ID == "" {
		service.ID = service.Name
	}
	if service.Environment == "" {
		service.Environment = "default"
	}
	if componentData, found := configMap.Data["components"]; found {
		var components []*config.ServiceComponent
		err := yaml.Unmarshal([]byte(componentData), &components)
		if err != nil {
			slog.Error("Error parsing service components",
				"service", service.ID,
				"component", configMap.Name,
				"namespace", configMap.Namespace,
				"error", err.Error())
			return
		}

		for _, c := range components {
			if c.Parameters == nil {
				c.Parameters = make(map[string]string)
			}
			if _, f := c.Parameters["Namespace"]; !f {
				c.Parameters["Namespace"] = configMap.Namespace
			}
		}
		service.Components = components
	}
	if event == watch.Added || event == watch.Modified {
		kw.c.SetService(&service)
	} else if event == watch.Deleted {
		kw.c.RemoveService(&service)
	}
}

func (kw *kubeWatcher) addStatusPagesToConfig(event watch.EventType, configMap *corev1.ConfigMap) {
	annotations := configMap.GetAnnotations()
	page := config.StatusPage{
		ID:          annotations["katastasi.io/page-id"],
		Name:        annotations["katastasi.io/name"],
		Contact:     annotations["katastasi.io/contact"],
		Owner:       annotations["katastasi.io/owner"],
		URL:         annotations["katastasi.io/url"],
		Environment: annotations["katastasi.io/environment"],
	}
	if page.Name == "" {
		page.Name = configMap.Name
	}
	if page.ID == "" {
		page.ID = page.Name
	}
	if page.Environment == "" {
		page.Environment = "default"
	}

	s := configMap.Data["services"]
	page.Services = strings.Split(s, ",")

	if event == watch.Added || event == watch.Modified {
		kw.c.SetStatusPage(&page)
	} else if event == watch.Deleted {
		kw.c.RemoveStatusPage(&page)
	}
}

func kubeClientConfig() *rest.Config {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}
	if viper.GetString("autoload.kubernetes.kubeconfig") != "" {
		kubeconfig = viper.GetString("autoload.kubernetes.kubeconfig")
	}

	var dataConfig *rest.Config
	var err error
	if viper.GetBool("autoload.kubernetes.in_cluster") == false {
		dataConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		dataConfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	return dataConfig
}
