#!/bin/bash

function init () {
    if [ -z $NAMESPACE ]; then namespace=red; else namespace=$NAMESPACE; fi
    # kubectl apply -f $(dirname -- $(readlink -fn -- "$0"))/init-trench-a.yaml --namespace=$namespace
    # kubectl apply -f $(dirname -- $(readlink -fn -- "$0"))/init-trench-b.yaml --namespace=$namespace
    kubectl wait --for=condition=ready pod -l app=lb-fe-attractor-a-1 -n $namespace --timeout=180s
    kubectl wait --for=condition=ready pod -l app=proxy-trench-a -n $namespace --timeout=180s
    # helm install $(pwd)/../../Meridio/examples/target/helm/ --name-template target-a --create-namespace --namespace $namespace --set applicationName=target-a --set default.trench.name=trench-a --set default.conduit.name=conduit-a
}

function end () {
    :
    # if [ -z $NAMESPACE ]; then namespace=red; else namespace=$NAMESPACE; fi
    # kubectl delete -f $(dirname -- $(readlink -fn -- "$0"))/init-trench-a.yaml --namespace=$namespace
    # kubectl delete -f $(dirname -- $(readlink -fn -- "$0"))/init-trench-b.yaml --namespace=$namespace
    # kubectl delete trench trench-a --namespace=$namespace
    # kubectl delete trench trench-b --namespace=$namespace
}

function configuration_new_ip () {
    if [ -z $NAMESPACE ]; then namespace=red; else namespace=$NAMESPACE; fi
    kubectl apply -f $(dirname -- $(readlink -fn -- "$0"))/new-vip.yaml --namespace=$namespace
    kubectl apply -f $(dirname -- $(readlink -fn -- "$0"))/configuration-new-vip.yaml --namespace=$namespace
}

function configuration_new_ip_revert () {
    if [ -z $NAMESPACE ]; then namespace=red; else namespace=$NAMESPACE; fi
    kubectl apply -f $(dirname -- $(readlink -fn -- "$0"))/init-trench-a.yaml --namespace=$namespace
    kubectl delete -f $(dirname -- $(readlink -fn -- "$0"))/new-vip.yaml --namespace=$namespace
}

# Required to call the corresponding function
$1 $@