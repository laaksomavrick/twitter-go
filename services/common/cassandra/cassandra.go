package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

type Client struct {
	cluster *gocql.ClusterConfig
	Session *gocql.Session
}

func NewClient(host string, keyspace string) (*Client, error) {
	log.Println("Connecting cassandra...")
	cluster := gocql.NewCluster("127.0.0.1")
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
