package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gitu/katastasi/pkg/core"
	"github.com/gitu/katastasi/pkg/serve"
	"github.com/spf13/viper"
	"log"

	// load all auth plugins!
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var version, commit, date = "unknown", "unknown", "unknown"

func main() {
	info := "\n" +
		" \n    __         __             __             _ \n   / /______ _/ /_____ ______/ /_____ ______(_)\n  / //_/ __ `/ __/ __ `/ ___/ __/ __ `/ ___/ / \n / ,< / /_/ / /_/ /_/ (__  ) /_/ /_/ (__  ) /  \n/_/|_|\\__,_/\\__/\\__,_/____/\\__/\\__,_/____/_/   \n\n" +
		"katastasi: \n" +
		"  version: " + version + "\n" +
		"  commit:  " + commit + "\n" +
		"  built:   " + date + ""
	fmt.Println(info)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/katastasi/")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("KATASTASI")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("prometheus.url", "http://localhost:9090")
	viper.SetDefault("cache.ttl", "1m")
	viper.SetDefault("queries", map[string]string{"one": "1"})
	viper.SetDefault("autoload", true)
	viper.SetDefault("autoload.kubernetes.config", "")
	viper.SetDefault("autoload.kubernetes.in_cluster", true)
	viper.SetDefault("autoload.namespaces.pages", []string{"changeme"})
	viper.SetDefault("autoload.namespaces.services", []string{"changeme"})

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	k := core.NewKatastasi(info)
	k.ReloadConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)

		k.ReloadConfig()
	})
	viper.WatchConfig()

	serve.StartServer(k)
}
