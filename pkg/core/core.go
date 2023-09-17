package core

import (
	"github.com/gitu/katastasi/pkg/config"
	"github.com/jellydator/ttlcache/v3"
	papi "github.com/prometheus/client_golang/api"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"sync"
	"text/template"
)

type KataStatus struct {
	Info    string
	Healthy bool
	LastMsg string
}

type Katastasi struct {
	mu               sync.Mutex
	Config           *config.Config
	KataStatus       *KataStatus
	statusCache      *ttlcache.Cache[string, *config.ServiceStatus]
	prometheusClient papi.Client
	kubeWatcher      *kubeWatcher
}

func NewKatastasi(info string) *Katastasi {

	client, err := papi.NewClient(papi.Config{
		Address: viper.GetString("prometheus.url"),
	})
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	ret := &Katastasi{
		Config: &config.Config{
			Environments: make(map[string]*config.Environment),
			Queries:      make(map[string]*template.Template),
		},
		KataStatus: &KataStatus{
			Info:    info,
			Healthy: true,
			LastMsg: "",
		},
		prometheusClient: client,
	}

	duration := viper.GetDuration("cache.ttl")
	loader := ttlcache.LoaderFunc[string, *config.ServiceStatus](
		func(c *ttlcache.Cache[string, *config.ServiceStatus], key string) *ttlcache.Item[string, *config.ServiceStatus] {
			return c.Set(key, ret.loadServiceStatus(key), ttlcache.DefaultTTL)
		})
	ret.statusCache = ttlcache.New[string, *config.ServiceStatus](
		ttlcache.WithLoader[string, *config.ServiceStatus](loader),
		ttlcache.WithTTL[string, *config.ServiceStatus](duration),
	)

	go ret.statusCache.Start()
	return ret
}

func loadQueries(c *config.Config) {
	for key, query := range viper.GetStringMapString("queries") {
		c.AddQuery(key, query, "config-file")
	}
}

func (k *Katastasi) ReloadConfig() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.kubeWatcher != nil {
		slog.Debug("Stopping kube watcher")
		k.kubeWatcher.stop()
		k.kubeWatcher = nil
	}

	c := &config.Config{
		Environments: make(map[string]*config.Environment),
		Queries:      make(map[string]*template.Template),
	}

	loadQueries(c)
	loadPages(c)
	loadServices(c)

	if viper.GetBool("autoload.active") {
		slog.Debug("watching kubernetes for changes")
		k.kubeWatcher = newKubeWatcher(c)
		err := k.kubeWatcher.start()
		if err != nil {
			log.Fatal("Error starting kube watcher: " + err.Error())
		}
	}

	k.Config = c
}

func loadPages(c *config.Config) {
	var pages []config.StatusPage
	err := viper.UnmarshalKey("pages", &pages)
	if err != nil {
		log.Fatal("Error loading pages in config file: " + err.Error())
	}
	for _, page := range pages {
		c.SetStatusPage(&page)
	}
}

func loadServices(c *config.Config) {
	var services []config.Service
	err := viper.UnmarshalKey("services", &services)
	if err != nil {
		log.Fatal("Error loading services in config file: " + err.Error())
	}
	for _, service := range services {
		c.SetService(&service)
	}
}
