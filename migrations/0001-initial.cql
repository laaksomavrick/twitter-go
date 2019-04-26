DROP KEYSPACE IF EXISTS twtr;

CREATE KEYSPACE twtr WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};