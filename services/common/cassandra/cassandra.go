package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

// TODO-7: handle disconnects from cassandra

type Client struct {
	cluster *gocql.ClusterConfig
	Session *gocql.Session
}

func NewClient(host string, keyspace string) (*Client, error) {
	log.Println("Connecting cassandra...")
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
