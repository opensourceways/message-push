package utils

import (
	"fmt"
	"github.com/gocql/gocql"
	"message-push/common/cassandra"
	"testing"
)

func TestCassandra(t *testing.T) {
	//cluster := gocql.NewCluster("127.0.0.1") // Cassandra节点的IP地址
	//cluster.Keyspace = "message_center"      // 替换为你的Keyspace名称
	//cluster.Port = 9042
	//session, _ := cluster.CreateSession()
	cassandra.Init()
	queryData(cassandra.Session())
}

func queryData(session *gocql.Session) error {
	var id gocql.UUID
	var name string
	var age int

	query := "SELECT id, name, age FROM example_table"
	iter := session.Query(query).Iter()
	for iter.Scan(&id, &name, &age) {
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", id, name, age)
	}
	if err := iter.Close(); err != nil {
		return err
	}
	return nil
}
