# install kube-prometheus
follow: https://prometheus-operator.dev/docs/prologue/quick-start/

    git clone https://github.com/prometheus-operator/kube-prometheus.git


## run katastasi locally
you will have to run a port forward to access the prometheus UI

    kubectl port-forward -n monitoring svc/prometheus-k8s 9090:9090



# install sample katastasi
you can install the sample katastasi deployment and config

    kubectl apply -f sample/sample-deployment.yml

# install sample katastasi config
you can install the sample katastasi config

    kubectl apply -f sample/sample-config.yml

# port forward to katastasi

    kubectl port-forward -n default svc/katastasi-service 31323:1323