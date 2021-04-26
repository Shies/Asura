package conf

import (
	"flag"

	xtime "Asura/app/time"

	"github.com/BurntSushi/toml"
)

var (
	Conf Config
	path string
)

// common config
type Common struct {
	Version    string
	Debug      bool
	Env 	   string
	Family     string
}

// final config
type Config struct {
	*Common
	HttpServer 		*HTTPServer
	RPCServer2      *RPCServer2
	Zookeeper  		*Zookeeper
	Log 	   		*Log 		 	 `toml:"xlog"`
	HttpClient 		*HttpClient
	RpcClient  		*RpcClient
	Redis	   		*RedisCluster
	Memcache   		*MemcacheCluster
	KafkaConsumer   *KafkaConsumerList
	KafkaProducer   *KafkaProducerList
}

// Degrade config
type DegradeConfig struct {
	Memcache *struct{
		*Memcache
		Idle 	int32
		Timeout xtime.Duration	// socket io read/write timeout
	}
	Expire 	int32
}

// HTTPServer http server settings.
type HTTPServer struct {
	Addrs        []string
	MaxListen    int32
	Timeout		 xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

// HttpClient
type HttpClient struct {
	DialTimeout	xtime.Duration	`toml:"dialtimeout"`
	Timeout	    xtime.Duration	`toml:"timeout"`
	KeepAlive   xtime.Duration 	`toml:"keepalive"`
}

// RPCServer rpc server settings.
type RPCServer struct {
	Proto  string
	Addr   string
	Group  string
	Color  string
	Weight int // weight of rpc server and also means num of client connections.
}

// RPCServer2 net/rpc service discover server settings.
type RPCServer2 struct {
	DiscoverOff bool
	Token       string
	Servers     []*RPCServer
	Zookeeper   *Zookeeper
}

// grpc client config
type RpcClient struct {
	Proto   string
	Addrs   []string
	Retry 	int64		 // 连接并行度
	Times   int64		 // 每连接请求次数
}

// Zookeeper Server&Client settings.
type Zookeeper struct {
	Root    string
	Addrs   []string
}

type KafkaProducerList struct {
	Test *KafkaProducer
}

// KafkaProducer kafka producer settings.
type KafkaProducer struct {
	Zookeeper 	*Zookeeper
	Brokers   	[]string
	Cluster		string
	Sync      	bool // true: sync, false: async
}

type KafkaConsumerList struct {
	Test *KafkaConsumer
}

// KafkaConsumer kafka client settings.
type KafkaConsumer struct {
	Offset    string	// true: new, false: old
	GroupID   string
	Topic     []string
	Addrs     []string
	Monitor   *HTTPServer // Consumer Ping Addr
	Redis     *KafkaRedisConfig
}

// Config Ali-DTS config.
type KafkaRedisConfig struct {
	Key          string
	Secret       string
	Group        string
	Topic        string
	Action       string // shoule be "pub" or "sub" or "pubsub"
	Buffer       int
	Name         string // redis name, for trace
	Proto        string
	Addr         string
	Auth         string
	Active       int // pool
	Idle         int // pool
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
}

// Log
type Log struct {
	Dir	string	`toml:"dir"`
}

// Memcache
type Memcache struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// MemcacheCluster
type MemcacheCluster struct {
	Idle 	int32			`toml:"idle"`
	Timeout xtime.Duration	`toml:"timeout"`	// socket io read/write timeout
	Node0   *Memcache		`toml:"node0"`
	Node1	*Memcache	`toml:"node1"`
	Node2   *Memcache	`toml:"node2"`
	Node3   *Memcache	`toml:"node3"`
	Node4   *Memcache	`toml:"node4"`
	Node5   *Memcache	`toml:"node5"`
	Node6   *Memcache	`toml:"node6"`
	Node7   *Memcache	`toml:"node7"`
	Node8   *Memcache	`toml:"node8"`
	Node9   *Memcache	`toml:"node9"`
}

// RedisCluster
type RedisCluster struct {
	Idle   		 int			  `toml:"idle"`
	Active 		 int			  `toml:"active"`
	IdleTimeout  xtime.Duration	  `toml:"idletimeout"`
	DialTimeout  xtime.Duration	  `toml:"dialtimeout"`
	ReadTimeout  xtime.Duration	  `toml:"readtimeout"`
	WriteTimeout xtime.Duration	  `toml:"writetimeout"`
	Node1  		 *RedisSync	  	  `toml:"node1"`
	Node2  		 *RedisSync	  	  `toml:"node2"`
}

// RedisSync
type RedisSync struct {
	Master 	*Redis	`toml:"master"`
	Slave   *Redis	`toml:"slave"`
}

// Redis
type Redis struct {
	Proto string `toml:"proto"`
	Host  string `toml:"host"`
	Port  string `toml:"port"`
	Auth  string `toml:"auth"`
}

// UDP server & client
type UDPServer struct {
	Proto string	`toml:"proto"`
	Host  string	`toml:"host"`
	Port  int8		`toml:"port"`
	Role  string	`toml:"role"`
}

func init() {
	flag.StringVar(&path, "conf", "", "default config path")
}

func ParseConfig() (err error) {
	_, err = toml.DecodeFile(path, &Conf)
	return
}

func TestConfig() (err error) {
	_, err = toml.DecodeFile("./config.example.toml", &Conf)
	return
}