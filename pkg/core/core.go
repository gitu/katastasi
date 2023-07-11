package core

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/gitu/katastasi/pkg/config"
	"github.com/jellydator/ttlcache/v3"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type KataStatus struct {
	Info    string
	Healthy bool
	LastMsg string
}

type Katastasi struct {
	QueryNamespaces []string
	DataConfig      *rest.Config
	mu              sync.Mutex
	Config          *config.Config
	KataStatus      *KataStatus
	StatusCache     *ttlcache.Cache[string, config.ServiceStatus]
}

func NewKatastasi(info string, queryNamespaces []string, dataConfig *rest.Config) *Katastasi {
	return &Katastasi{
		QueryNamespaces: queryNamespaces,
		DataConfig:      dataConfig,
		Config: &config.Config{
			Environments: make(map[string]config.Environment),
			Queries:      make(map[string]string),
		},
		KataStatus: &KataStatus{
			Info:    info,
			Healthy: true,
			LastMsg: "",
		},
	}
}

func (k *Katastasi) Start() {
	k.buildAndStartCache()
	err := k.reloadConfig()
	if err != nil {
		log.Printf("error loading config: %v", err)
	}
	k.startMonitor()
}

func (k *Katastasi) buildAndStartCache() {
	loader := ttlcache.LoaderFunc[string, config.ServiceStatus](
		func(c *ttlcache.Cache[string, config.ServiceStatus], key string) *ttlcache.Item[string, config.ServiceStatus] {
			return c.Set(key, k.loadServiceStatus(key), ttlcache.DefaultTTL)
		})
	k.StatusCache = ttlcache.New[string, config.ServiceStatus](
		ttlcache.WithLoader[string, config.ServiceStatus](loader),
		ttlcache.WithTTL[string, config.ServiceStatus](1*time.Minute),
	)
	go k.StatusCache.Start()
}

func (k *Katastasi) reloadConfig() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	ret := &config.Config{
		Environments: make(map[string]config.Environment),
	}
	clientSet, err := kubernetes.NewForConfig(k.DataConfig)
	if err != nil {
		return err
	}

	err = k.addQueriesToConfig(ret, clientSet)
	if err != nil {
		return err
	}

	err = k.addServicesToConfig(ret, clientSet)
	if err != nil {
		return err
	}

	err = k.addStatusPagesToConfig(ret, clientSet)
	if err != nil {
		return err
	}

	k.Config = ret

	return nil
}

func (k *Katastasi) addServicesToConfig(ret *config.Config, clientSet *kubernetes.Clientset) error {
	list, err := clientSet.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{
		LabelSelector: "katastasi.io=service",
	})
	if err != nil {
		return err
	}
	for _, configMap := range list.Items {
		annotations := configMap.GetAnnotations()
		service := config.Service{
			ID:          annotations["katastasi.io/service-id"],
			Name:        annotations["katastasi.io/name"],
			Contact:     annotations["katastasi.io/contact"],
			Owner:       annotations["katastasi.io/owner"],
			URL:         annotations["katastasi.io/url"],
			Environment: annotations["katastasi.io/env"],
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
			var components []config.ServiceComponent
			err = yaml.Unmarshal([]byte(componentData), &components)
			if err != nil {
				log.Printf("Error parsing service components %s for ConfigMap %s in %s: %s", service.Name, configMap.Name, configMap.Namespace, err.Error())
				continue
			}
			service.Components = components
		}

		err := ret.AddService(service)
		if err != nil {
			log.Printf("Error adding service %s for ConfigMap %s in %s: %s", service.Name, configMap.Name, configMap.Namespace, err.Error())
		}
	}
	return nil
}

func (k *Katastasi) addQueriesToConfig(ret *config.Config, clientSet *kubernetes.Clientset) error {
	for _, queryNamespace := range k.QueryNamespaces {
		list, err := clientSet.CoreV1().ConfigMaps(queryNamespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: "katastasi.io=queries",
		})
		if err != nil {
			return err
		}
		for _, configMap := range list.Items {
			for name, query := range configMap.Data {
				if ret.Queries[name] != "" {
					log.Printf("Duplicate query name %s in %s (%s)", name, configMap.Name, configMap.Namespace)
					continue
				}
				ret.Queries[name] = query
			}
		}
	}
	return nil
}

func (k *Katastasi) addStatusPagesToConfig(ret *config.Config, clientSet *kubernetes.Clientset) error {
	list, err := clientSet.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{
		LabelSelector: "katastasi.io=page",
	})
	if err != nil {
		return err
	}
	for _, configMap := range list.Items {
		annotations := configMap.GetAnnotations()
		page := config.StatusPage{
			ID:          annotations["katastasi.io/page-id"],
			Name:        annotations["katastasi.io/name"],
			Contact:     annotations["katastasi.io/contact"],
			Owner:       annotations["katastasi.io/owner"],
			URL:         annotations["katastasi.io/url"],
			Environment: annotations["katastasi.io/env"],
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

		err := ret.AddStatusPage(page)
		if err != nil {
			log.Printf("Error adding status page %s for ConfigMap %s in %: %s", page.Name, configMap.Namespace, configMap.Name, err.Error())
		}
	}
	return nil
}

func (k *Katastasi) startMonitor() {
	go func() {
		for {
			err := backoff.Retry(func() error {
				err := k.reloadConfig()
				if err != nil {
					log.Printf("error reloading config: %v", err)
				}
				return err
			}, backoff.NewExponentialBackOff())
			if err != nil {
				log.Fatalf("permanent error reloading config: %v", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func (k *Katastasi) GetPageStatus(env string, page string) config.PageStatus {
	p := config.PageStatus{
		Status:     config.Unknown,
		LastUpdate: time.Now(),
		Services:   map[string]config.ServiceStatus{},
	}
	if _, f := k.Config.Environments[env]; !f {
		return p
	}
	if _, f := k.Config.Environments[env].StatusPages[page]; !f {
		return p
	}
	for _, service := range k.Config.Environments[env].StatusPages[page].Services {
		p.Services[service] = k.GetStatusOfService(env, service)
		if p.Services[service].Status.IsHigherThan(p.Status) {
			p.Status = p.Services[service].Status
		}
		if p.Services[service].LastUpdate.Before(p.LastUpdate) {
			p.LastUpdate = p.Services[service].LastUpdate
		}
	}
	return p
}

func (k *Katastasi) GetStatusOfService(env string, service string) config.ServiceStatus {
	s := k.StatusCache.Get(ToCacheKey(env, service))
	return s.Value()
}

func ToCacheKey(env string, service string) string {
	return env + "|" + service
}

func (k *Katastasi) loadServiceStatus(key string) config.ServiceStatus {
	keyval := strings.SplitN(key, "|", 2)
	env := keyval[0]
	service := keyval[1]

	ret := config.ServiceStatus{
		ID:         service,
		Status:     config.Unknown,
		LastUpdate: time.Now(),
		Components: make([]config.ComponentStatus, 0),
	}

	if _, f := k.Config.Environments[env]; !f {
		return ret
	}
	if _, f := k.Config.Environments[env].Services[service]; !f {
		return ret
	}

	for _, component := range k.Config.Environments[env].Services[service].Components {
		nc := k.funcName(component)
		if nc.Status.IsHigherThan(ret.Status) {
			ret.Status = nc.Status
		}
		ret.Components = append(ret.Components, nc)
	}

	return ret
}

func (k *Katastasi) funcName(component config.ServiceComponent) config.ComponentStatus {
	nc := config.ComponentStatus{
		Name:         component.Name,
		StatusString: component.Description,
		Status:       config.OK,
	}
	for _, cond := range component.Conditions {
		if rand.Int31n(30) > 20 {
			nc.Status = cond.Severity
		}
	}
	return nc
}
