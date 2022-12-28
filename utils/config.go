package common_utils

import (
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type Database struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Name     string `json:"name" yaml:"name"`
	Url      string `json:"source" yaml:"source"`
}

type Redis struct {
	Host          string        `json:"host" yaml:"host"`
	Port          string        `json:"port" yaml:"port"`
	User          string        `json:"user" yaml:"user"`
	Password      string        `json:"password" yaml:"password"`
	DB            int           `json:"db" yaml:"db"`
	PoolSize      int           `json:"poolSize" yaml:"poolSize"`
	PoolTimeout   time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
	CacheDuration time.Duration `json:"cacheDuration" yaml:"cacheDuration"`
}

type Cassandra struct {
	HostPort []string `json:"hostPort" yaml:"hostPort"`
	User     string   `json:"user" yaml:"user"`
	Password string   `json:"password" yaml:"password"`
	Keyspace string   `json:"keyspace" yaml:"keyspace"`
}

type Elastic struct {
	Addresses     []string    `json:"addresses"`
	Username      string      `json:"username"`
	Password      string      `json:"password"`
	APIKey        string      `json:"apiKey"`
	Header        http.Header // Global HTTP request header.
	EnableLogging bool        `json:"enableLogging"`
}

type KafkaConfig struct {
	Brokers    []string `json:"brokers"`
	GroupID    string   `json:"groupID"`
	InitTopics bool     `json:"initTopics"`
}

type TopicConfig struct {
	TopicName         string `json:"topicName"`
	Partitions        int    `json:"partitions"`
	ReplicationFactor int    `json:"replicationFactor"`
}

type Kafka struct {
	Host            string               `json:"host" yaml:"host"`
	Port            string               `json:"port" yaml:"port"`
	User            string               `json:"user" yaml:"user"`
	Password        string               `json:"password" yaml:"password"`
	IsRemote        bool                 `json:"isRemote" yaml:"isRemote"`
	KafkaConfig     KafkaConfig          `json:"kafkaConfig" yaml:"kafkaConfig"`
	TopicConfig     TopicConfig          `json:"topicConfig" yaml:"topicConfig"`
	PublisherConfig KafkaPublisherConfig `mapstructure:"publisherConfig" json:"publisherConfig" yaml:"publisherConfig"`
}

type KafkaPublisherConfig struct {
	Topic             string         `json:"topic" yaml:"topic"`
	TopicPrefix       string         `json:"topicPrefix" yaml:"topicPrefix"`
	Partitions        int            `json:"partitions" yaml:"partitions"`
	ReplicationFactor int            `json:"replicationFactor" yaml:"replicationFactor"`
	Headers           []kafka.Header `json:"headers" yaml:"headers"`
}

type Jaeger struct {
	ServiceName string `json:"serviceName" yaml:"serviceName"`
	Host        string `json:"host" yaml:"host"`
	Port        string `json:"port" yaml:"port"`
	Enable      bool   `json:"enable" yaml:"enable"`
	LogSpans    bool   `json:"logSpans" yaml:"logSpans"`
}

type S3 struct {
	Endpoint           string        `json:"endpoint" yaml:"endpoint"`
	SecretKey          string        `json:"secretKey" yaml:"secretKey"`
	AccessKey          string        `json:"accessKey" yaml:"accessKey"`
	Region             string        `json:"region" yaml:"region"`
	PublicBucketName   string        `json:"publicBucketName" yaml:"publicBucketName"`
	PrivateBucketName  string        `json:"privateBucketName" yaml:"privateBucketName"`
	PublicUrl          string        `json:"publicUrl" yaml:"publicUrl"`
	PreSignUrlDuration time.Duration `json:"preSignUrlDuration" yaml:"preSignUrlDuration"`
}

type Mail struct {
	Driver     string `json:"driver" yaml:"driver"`
	Host       string `json:"host" yaml:"host"`
	Port       string `json:"port" yaml:"port"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	Encryption bool   `json:"encryption" yaml:"encryption"`
}

type JWT struct {
	Secret  string        `json:"secret" yaml:"secret"`
	Expires time.Duration `json:"expires" yaml:"expires"`
}

type HTTP struct {
	Port string `json:"port" yaml:"port"`
}

type GRPC struct {
	Port string `json:"port" yaml:"port"`
}

type Metrics struct {
	Port          string   `json:"port" yaml:"port"`
	IgnoreLogUrls []string `json:"ignoreLogUrls" yaml:"ignoreLogUrls"`
}

type Server struct {
	Name    string `json:"name" yaml:"name"`
	Host    string `json:"host" yaml:"host"`
	GRPC    GRPC
	HTTP    HTTP
	Metrics Metrics
}

type Config struct {
	Server    Server
	Database  Database
	Cassandra Cassandra
	Elastic   Elastic
	Redis     Redis
	JWT       JWT
	Kafka     Kafka
	S3        S3
	Jaeger    Jaeger
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = GetPath()
	}
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.JWT.Expires = conf.JWT.Expires * time.Second
	conf.Database.Url = conf.Database.Driver + "://" + conf.Database.User + ":" + conf.Database.Password + "@" + conf.Database.Host + conf.Database.Port + "/" + conf.Database.Name + "?sslmode=disable"

	return conf, nil
}

func GetPath() string {
	dir := getSourcePath()
	return dir + "/../config.yaml"
}

func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
