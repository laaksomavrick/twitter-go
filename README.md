# twitter-go

Twitter-go is an application api intended to back a minimal feature set of twitter. Its development serves as a fun learning exercise to explore an event driven microservice architecture, using Go. I've been curious about Go for a variety of reasons (performance, native type system, binaries, pragmatic ecosystem) and want to explore it and its ecosystem.

Moreover, microservice backends are becoming more ubiquitous due to the organizational benefits they offer (independent deployability, independent scalability, fault tolerance), and an event driven architecture is a common way of developing loosely coupled services.

Further, writing the infrastructure for managing and deploying a microservice backend likely will and has proven to be a fantastic learning exercise for crystalizing knowledge I've acquired about Kubernetes, Docker and Helm.

## What is the app?

- A user makes post
- A user has a viewable list of their own posts
- A user can subscribe to other users
- A user has an activity feed of those they follow's posts (chronological order)

### This app needs to provide the ability:

- To create a user
- To login a user
- To get a list of a user's posts
- To follow other users
- To retrieve an aggregated activity feed of posts from those they follow

### Ergo, service breakdown:

- API gateway (Entry point into the backend; maps http to n rpc calls)
- User service (CRUD for users; user authorization)
- Tweet service (Adding to user tweet list; "my tweets")
- Follower service (Managing user - user follows/followers relationships; "my followers/follower count")
- Feed service (Aggregating user activity feed; "my feed")

### Service philosophy:

- Services are responsible solely for their domain (biz logic, tables)
- Services will publish events about their domain for other services to subscribe to as required
- Services that require data not belonging to their domain will embrace denormalization and eventual consistency
- Services should be written as dumb as possible and avoid pre-emptive abstractions; YAGNI

### Technical choices:

- Go for gateway and application code
- RabbitMQ for a message bus (rpc, pub/sub)
- Cassandra for a NoSQL datastore
- Docker, Kubernetes and Helm for deployment
- Some bash scripts for convenience, glueing things

## How do I get it running on k8s?

It's been a while (a few months) since I initially set up and deployed this project. However, I'll do my best to recount the steps required at a high level.

Requirements:
- An adequately provisioned k8s cluster, running on GCP (I believe my set up for dev had 3 nodes and 8gb of RAM; the default options for new projects).
- `kubectl` on your local machine, properly configured to talk to your k8s cluster.
- `Docker` on your local machine.
- `cqlsh` on your local machine.

First step is installing helm onto your cluster. The helm documentation will be more helpful than myself for this, but the gist is installing helm onto your local machine (e.g through `brew install helm`), and then running `scripts/setup-k8s.sh`, which will create a new namespace (`twtr-dev`) on your cluster and give tiller (the server side part of helm) permissions to manage this namespace. This, while not totally necessary, is important: in a production application, we probably don't want to give tiller free reign over an entire cluster. By giving tiller permissions only to the namespace where our application will reside, we remove a lot of surface area for security vulnerabilities or developer mistakes.

After that, the hard part is over. You'll now need to modify some of my scripts in the `scripts` directory and charts in the `helm` directory to replace my google project id with yours (`:%s/precise-clock-244301/YOUR_PROJECT_ID`). This is so we can build and upload docker images to google's container registry, and pull those images down later when we deploy the application. Once you've modified all references to the google project id, you can run `make docker-build-push` and grab a coffee while we build all the go services and push them up.

Now, coffee in hand, you should be able to run `make helm-install`. In retrospective, there is no container to create the Cassandra keyspace and run the migrations as a job, so you'll likely have to `kubectl port-forward twtr-dev-cassandra-0 9200:9200 9042:9042` and run `make migrate` so Cassandra has a keyspace and the proper tables. In a production set up, we'd have an initContainer configured to check for any new keyspaces/migrations and run them prior to any services being allowed to run in k8s. If things bork because of this, run `make helm-purge` once you've run the migrations and then `make helm-install` again. Since Cassandra has a pvc configured, the keyspace + tables will be there the second time around, and you can ignore it.

Following this, you should see the following output from `kubectl get pods`

```
tiller-deploy-74d85979b6-xq579        1/1       Running   0          20d
twtr-dev-cassandra-0                  1/1       Running   0          5d16h
twtr-dev-feeds-5b64b5f899-d9g2k       1/1       Running   0          5d16h
twtr-dev-followers-748b8d7b45-xlhsz   1/1       Running   0          5d16h
twtr-dev-gateway-6fc7d454d6-9fqtm     1/1       Running   0          5d16h
twtr-dev-rabbitmq-0                   1/1       Running   0          5d16h
twtr-dev-traefik-75b84ddd85-n65bc     1/1       Running   0          5d16h
twtr-dev-tweets-6476944c94-5d7jk      1/1       Running   0          5d16h
twtr-dev-users-64955776fb-ff6jb       1/1       Running   0          5d16h
```

If you want to actually send requests to the api, check what the ip of your ingress controller is via `kubectl get services | grep traefik` and observe the address of the `EXTERNAL IP`

```
twtr-dev-traefik             LoadBalancer   10.48.13.194   35.231.22.185   80:30004/TCP,443:31478/TCP   5d16h
```

Then, modify your `/etc/hosts` to point `twtr-dev.com` and `traefik.dashboard.com` (if you care about it) to that ip address:

```
/etc/hosts

127.0.0.1   localhost
255.255.255.255 broadcasthost
::1             localhost
127.0.0.1 192.168.0.143
35.231.22.185 traefik.dashboard.com
35.231.22.185 twtr-dev.com
```

To confirm things are working, run `curl -i http://twtr-dev.com/healthz` and you should see an `ok` response in the response body. Congratulations!

## How do I get it running locally for development?

This is much easier than configuring it for production. Assuming you have Docker and `cqlsh` on your local machine, run the following commands:

- `make up`
- `make migrate`
- `make run`

To run the integration tests, run `make test`.

You should be able to run `curl -i http://localhost:3002/healthz` to confirm things are working.
