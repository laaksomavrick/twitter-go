## Tests

The integration test suite validates the application functions correctly. The tests generally send a request to the API gateway and observe the response, verifying the payload has the correct values where appropriate and the status code of the response.

Because of this, it should be easy to at minimum maintain the API contract and essential behaviour while adding new functionality or refactoring existing functionality. A developer will know when they've broken the API contract, or caused an error.

#### Running the tests

First, start cassandra and rabbit locally:

`make up`

Then, migrate the database so cassandra has the proper keyspace and tables:

`make migrate`

Then, run the test suite:

`make test`