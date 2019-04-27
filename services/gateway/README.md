## API Gateway

The API gateway is the entry point to the backend. It defines the api surface for the system, assembles responses via rpc to one or many services, and handles global concerns such as authorization. Having an API gateway makes service discovery simple - only the gateway needs to be exposed to the public internet, and all rpc calls can be dispatched to the event bus (RabbitMQ), which will handle routing and response logic. However, it does introduce a single point of failure in the system - if the gateway goes down, the backend will no longer be able to respond to incoming requests.

#### API Schema

##### Users

- `POST /users`
    - Create a new user
