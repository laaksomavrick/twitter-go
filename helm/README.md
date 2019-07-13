# Helm

This directory and its subdirectories contain a set of helm charts, each of which configures a component of the kubernetes cluster required for the twitter-go backend.
These components fall into a few categories:

* Business logic/domain level services, (e.g the users service or the tweets service)
* Data storage (cassandra)
* Messaging (rabbitmq)
* Ingress/load balancing (traefik)

Helm was chosen because it grants us:

* access to a templating system for repetitive configuration (DRY)
* a set of community tailored charts for common infrastructure components (quick wins)
* the ability to version our deployments
* simple deployments and deployment updates
