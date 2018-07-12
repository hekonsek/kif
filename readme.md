# kif

**kif** (Kube Is Fine) is a small tool providing simple and opinionated SRE infrastructure based on Kubernetes.

When installed into Kubernetes cluster Kif provides the following pre-configured features:
- [Nginx Ingress](https://github.com/kubernetes/ingress-nginx)
- [ACME certificate manager (Jetstack Cert Manager)](https://github.com/jetstack/cert-manager)
- [Prometheus](https://prometheus.io)

## Usage

To install kif platform into a bare metal cluster execute the following command...

    kif create platform --ingress-ip=1.2.3.4 --cert-email=admin@example.com

Option `ingress-ip` indicates the IP of machine that should be used as Ingress load balancer. Option `cert-email` tells
what administrator e-mail address should be used in ACME registration process when generating new CA certificates.

### Dry run

If you are interested only in generating Helm charts, not in executing them, use `--dry-run` option:

    kif create platform --dry-run --ingress-ip=1.2.3.4 --cert-email=admin@example.com