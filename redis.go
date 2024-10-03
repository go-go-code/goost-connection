package connection

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	logger "gitlab.com/gr50y.world/goost-logger"
)

var _redis *redisConnection

type redisConnection struct {
	Client redis.UniversalClient
	status bool
}

func (r *redisConnection) Close() {
	if !r.status {
		return
	}

	if err := r.Client.Close(); err != nil {
		logger.ErrorF("Closing Redis Error : %v", err)
	} else {
		logger.Info("Redis Connection Is Close")
	}

	r.status = false
}

func (r *redisConnection) SetAlive() {
	r.status = true
}

func (r *redisConnection) IsAlive() bool {
	return r.status
}

func NewRedisConnection() (client redis.UniversalClient) {
	if _redis == nil {
		initRedisConnection()
	}

	return _redis.Client
}

func initRedisConnection() {
	_redis = &redisConnection{}

	if _redis.IsAlive() {
		return
	}

	if enabled, ok := cfg["redis_enable"].(bool); !ok || !enabled {
		logger.Info("⚠️ Redis is Disabled ⚠️")
		return
	}

	username, _ := cfg["redis_username"].(string)
	password, _ := cfg["redis_password"].(string)

	opt := &redis.UniversalOptions{
		Username: username,
		Password: password,

		//连接池容量及闲置连接数量
		PoolSize: 50, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		// MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		// DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		// ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		// WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		// PoolTimeout:  5 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		// IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		// IdleTimeout:        10 * time.Second, //闲置超时，默认5分钟，-1表示取消闲置超时检查
		// MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		// MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		// MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		// MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },

		// ReadOnly = true，只择 Slave Node
		// ReadOnly = true 且 RouteByLatency = true 将从 slot 对应的 Master Node 和 Slave Node， 择策略为: 选择PING延迟最低的点
		// ReadOnly = true 且 RouteRandomly = true 将从 slot 对应的 Master Node 和 Slave Node 选择，选择策略为: 随机选择

		// ReadOnly:       true,
		// RouteRandomly:  true,
		// RouteByLatency: true,
	}

	host, ok := cfg["redis_host"].(string)
	if !ok || host == "" {
		host = "127.0.0.1"
	}

	port, ok := cfg["redis_port"].(string)
	if !ok || port == "" {
		port = "6379"
	}

	address := host + ":" + port

	opt.Addrs = strings.Split(address, ",")

	client := redis.NewUniversalClient(opt)

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		client.Close()
		logger.ErrorF("CONNECTing Redis ERROR , %v", err)
		return
	}

	_redis.Client = client
	_redis.SetAlive()
	add(_redis)

	logger.Info("Connecting Redis Success")

	return
}
