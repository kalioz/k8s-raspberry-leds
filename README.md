# k8s-raspberry-leds

This project is a software meant to be deployed on a kubernetes cluster on raspberry pi nodes equipped with a led hat (8 RGB pixels).
It will change the color of the pixel depending on the pods deployed on the node.

## Installation on the Kubernetes cluster
this project is meant for a k3s cluster, but should work perfectly fine with a k8s cluster.

to install it, build the image and deploy it on your cluster with helm

```console
cd helm && helm install cluster-monitoring-leds .
``` 

This will install a daemonset which will create one pod per node. tweak the chart if you want it to be deployed only on some nodes.

## Build the project

```console
docker build . -t cluster-monitoring-leds
```

Note that on windows the installation of `github.com/stianeikeland/go-rpio/v4` can fail. You need to set the environment variable `GOOS=linux` to make it work (or trust blindly your IDE, which I do)