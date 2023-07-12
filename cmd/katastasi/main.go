package main

import (
	"flag"
	"fmt"
	"github.com/gitu/katastasi/pkg/core"
	"github.com/gitu/katastasi/pkg/serve"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"

	// load all auth plugins!
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var outOfCluster = flag.Bool("out-of-cluster", false, "load data via kubeconfig")
var inCluster = flag.Bool("in-cluster", false, "load data via deployment")
var namespaces = flag.String("namespaces", "katastasi", "namespaces to load queries from, split by comma")
var prometheusUrl = flag.String("prometheus-url", "http://prometheus:9090", "prometheus url")

var version, commit, date = "unknown", "unknown", "unknown"

func main() {
	info := "\n" +
		" \n    __         __             __             _ \n   / /______ _/ /_____ ______/ /_____ ______(_)\n  / //_/ __ `/ __/ __ `/ ___/ __/ __ `/ ___/ / \n / ,< / /_/ / /_/ /_/ (__  ) /_/ /_/ (__  ) /  \n/_/|_|\\__,_/\\__/\\__,_/____/\\__/\\__,_/____/_/   \n\n" +
		"katastasi: \n" +
		"  version: " + version + "\n" +
		"  commit:  " + commit + "\n" +
		"  built:   " + date + ""
	fmt.Println(info)

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var config *rest.Config
	var err error
	if *outOfCluster {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	k := core.NewKatastasi(info, strings.Split(*namespaces, ","), config, *prometheusUrl)
	k.Start()

	serve.StartServer(k)

}
