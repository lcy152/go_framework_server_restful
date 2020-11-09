package service

import (
	"sync"
	"tumor_server/db"

	redis "tumor_server/redis"

	uuid "github.com/satori/go.uuid"
)

var instance *Container
var once sync.Once

func GetContainerInstance() *Container {
	if instance == nil {
		panic("Container not an instance")
	}
	return instance
}

func NewContainerInstance(serverConfig *ServerConfig) *Container {
	once.Do(
		func() {
			instance = NewContainer(serverConfig)
		},
	)
	return instance
}

type ServerConfig struct {
	Host                     string
	Port                     int
	RedisIP                  string
	RedisPort                int
	DatabaseUrl              string
	ExpireTime               int
	LogPath                  string
	LogLevel                 int
	BlackListTime            int
	BlackListLimiteTime      int
	BlackListLimitCount      int
	ShortMessageInvalidTime  int
	ShortMessageSpaceTime    int
	ShortMessageExpireTime   int
	ShortMessageLimitedCount int
	RabbitMQIP               string
	RabbitMQVHost            string
	RabbitMQVUserName        string
	RabbitMQVPassword        string
}

type Container struct {
	DB           *db.Database
	Dispatch     *Dispatcher
	Config       *ServerConfig
	RedisService *redis.RedisService
	RabbitMQ     *RabbitMQ
}

func NewContainer(serverConfig *ServerConfig) *Container {
	container := new(Container)
	container.Config = serverConfig
	container.RedisService = redis.ConnectRedis(serverConfig.RedisIP, serverConfig.RedisPort, serverConfig.ExpireTime, serverConfig.BlackListLimiteTime)
	InitRedis(container.RedisService)
	url := "mongodb://datu_super_root:c74c112dc3130e35e9ac88c90d214555__strong@" + container.Config.DatabaseUrl + "/datu_data"
	container.DB = db.NewDatabase(url)
	container.Dispatch = NewDispatcher()
	MQURL := "amqp://" + container.Config.RabbitMQVUserName + ":" + container.Config.RabbitMQVPassword + "@" + container.Config.RabbitMQIP + "/" + container.Config.RabbitMQVHost
	container.RabbitMQ = NewRabbitMQ(container, MQURL)
	return container
}

func NewUUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}
