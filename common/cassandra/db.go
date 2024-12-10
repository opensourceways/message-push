package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func Init(cfg *Config) error {
	return nil
	cluster := gocql.NewCluster(cfg.Host) // Cassandra节点的IP地址
	cluster.Keyspace = cfg.KeySpace       // 替换为你的Keyspace名称
	cluster.Port = cfg.Port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Password: cfg.Pwd,
		Username: cfg.User,
	}
	sessionInstance, err := cluster.CreateSession()
	if err != nil {
		log.Println("cassandra init session err:", err)
	}
	session = sessionInstance
	//defer sessionInstance.Close()
	return nil
}

// DB returns the current database instance.
func Session() *gocql.Session {
	return session
}
