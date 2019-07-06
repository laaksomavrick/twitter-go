## Ready service

This service is used as an init container during deployments to confirm RabbitMQ and Cassandra are ready to start accepting connections. My bash skills are lacking, so I opted to do this in Go :). As a consequence, it is nothing fancy.