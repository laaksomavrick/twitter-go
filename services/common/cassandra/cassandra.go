package cassandra

import (
	"github.com/gocql/gocql"
)

// TODO: handle disconnects from cassandra

// Client represents an active session (connection) to a particular cassandra cluster
type Client struct {
	cluster *gocql.ClusterConfig
	Session *gocql.Session
}

// NewClient synchronously constructs a cassandra.Client
func NewClient(host string, keyspace string) (*Client, error) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &Client{
		cluster: cluster,
		Session: session,
	}, nil
}
