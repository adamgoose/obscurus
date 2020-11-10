# Obscurus

Obscurus is a simple web-based tool for sharing secrets via a short-lived
one-time-use link.

You can view a preview [here](https://minio.enge.me/dropshare/Screen-Recording-2020-11-10-01-32-06.mp4).

## Developing

Obscurus requires Vault in Kubernetes. The easiest way to spin up a dev
environment is by using [Tilt](https://tilt.dev). Prepare a local Kubernetes
cluster with something like [k3d](https://github.com/rancher/k3d), then simply
run `tilt up`. To access the web UI, run
`kubectl -n obscurus port-forward svc/obscurus 8000:80` and access
http://localhost:8000 in your browser.