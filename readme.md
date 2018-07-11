# skrt

**skrt** (Simple Kube Reliability Toolkit) is a simple tool providing simple and opinionated SRE infrastructure based on Kubernetes.

## Usage

To install kif platform into a bare metal cluster execute the following command...

    skrt create platform --ingress-ip=1.2.3.4

...where `ingress-ip` indicates the IP of machine that should be used as Ingress load balancer.

### Dry run

If you are interested only in generating Helm charts, not in executing them, use `--dry-run` option:

    skrt create platform --dry-run --ingress-ip=1.2.3.4