package main

import (
	"flag"
	"fmt"
	"github.com/gitu/katastasi/pkg/serve"

	// load all auth plugins!
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var outOfCluster = flag.Bool("out-of-cluster", false, "load data via kubeconfig")
var namespaces = flag.String("namespaces", "katastasi", "namespaces to load configmaps from")
var filterNamespaces = flag.String("crdNamespaces", "default,katastasi", "namespaces to filter fetching CRDs by")

var version, commit, date = "unknown", "unknown", "unknown"

func main() {
	info := "\n" +
		" \n    __         __             __             _ \n   / /______ _/ /_____ ______/ /_____ ______(_)\n  / //_/ __ `/ __/ __ `/ ___/ __/ __ `/ ___/ / \n / ,< / /_/ / /_/ /_/ (__  ) /_/ /_/ (__  ) /  \n/_/|_|\\__,_/\\__/\\__,_/____/\\__/\\__,_/____/_/   \n\n" +
		"katastasi: \n" +
		"  version: " + version + "\n" +
		"  commit:  " + commit + "\n" +
		"  built:   " + date + ""
	fmt.Println(info)
	flag.Parse()
	serve.StartServer(info, *outOfCluster)
}
