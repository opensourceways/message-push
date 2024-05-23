package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func Init(cfg *Config) error {
	cluster := gocql.NewCluster(cfg.Host) // Cassandra节点的IP地址
	cluster.Keyspace = cfg.KeySpace       // 替换为你的Keyspace名称
	cluster.Port = cfg.Port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Password: cfg.Pwd,
		Username: cfg.User,
	}
	sessionInstance, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	session = sessionInstance
	//defer sessionInstance.Close()
	return nil
}

// DB returns the current database instance.
func Session() *gocql.Session {
	return session
}

func createTable(session *gocql.Session) error {
	query := `
        CREATE TABLE IF NOT EXISTS example_table (
            id UUID PRIMARY KEY,
            name TEXT,
            age INT
        )
    `
	return session.Query(query).Exec()
}

func insertData(session *gocql.Session) error {
	id := gocql.TimeUUID()
	query := "INSERT INTO example_table (id, name, age) VALUES (?, ?, ?)"
	return session.Query(query, id, "John", 30).Exec()
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
