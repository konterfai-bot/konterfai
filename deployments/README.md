# Deployments

This directory contains example configurations for running konterfAI with different reverse proxies.
At the moment there are examples for [nginx](nginx) and [traefik](traefik).
There are two examples providing full stacks with prometheus and grafana for monitoring 
([traefik-prometheus-grafana](traefik-prometheus-grafana) and [nginx-prometheus-grafana](nginx-prometheus-grafana)).

There are two folders in this directory that are used by the `traefik-prometheus-grafana` and the
`nginx-prometheus-grafana` examples. These folders contain the configuration files for the Prometheus and Grafana
services. The folders in question are `_prometheus-config` and `_grafana-provisioning`.

Default credentials for Grafana are `admin:kponterfai` for prometheus and grafana. 
You need to change these credentials in the docker-compose files in the deployment folders and in the 
`_prometheus-config` folder. The prometheus folder requjires a bcrypt hash of the password.

**Note:** These examples are not intended for copy & paste usage.
Make sure to read them carefully and adjust them to your needs.

**WARNING:** IMPROPER CONFIGURATION WILL HAVE NEGATIVE EFFECTS ON YOUR SEO
