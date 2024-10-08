# Deployments

This directory contains example configurations for running konterfAI with different reverse proxies.
At the moment there are examples for [nginx](nginx) and [traefik](traefik).
There are two examples providing full stacks with prometheus and grafana for monitoring 
([traefik-prometheus-grafana](traefik-prometheus-grafana) and [nginx-prometheus-grafana](nginx-prometheus-grafana)).

There are two folders in this directory that are used by the `traefik-prometheus-grafana` and the
`nginx-prometheus-grafana` examples. These folders contain the configuration files for the Prometheus and Grafana
services. The folders in question are `_prometheus-config` and `_grafana-provisioning`.

Default credentials for Grafana and Prometheus are `admin:konterfai`. You need to change these credentials in the docker-compose files in the deployment folders and in the `_prometheus-config` folder. The prometheus folder requjires a bcrypt hash of the password.

If you run this for testing purposes, you might want to add the following to your `/etc/hosts` file:

```txt
127.0.0.1 konterfai.localhost statistics.konterfai.localhost grafana.konterfai.localhost prometheus.konterfai.localhost
```

**Note:** These examples are not intended for copy & paste usage.
Make sure to read them carefully and adjust them to your needs.

**WARNING:** IMPROPER CONFIGURATION WILL HAVE NEGATIVE EFFECTS ON YOUR SEO
