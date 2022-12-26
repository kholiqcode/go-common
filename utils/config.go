package common_utils

import (
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
	PoolSize      int           `json:"pool_size" yaml:"pool_size"`
	PoolTimeout   time.Duration `json:"pool_timeout" yaml:"pool_timeout"`
	CacheDuration time.Duration `json:"cache_duration" yaml:"cache_duration"`
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
	IsRemote        bool                 `json:"is_remote" yaml:"is_remote"`
	KafkaConfig     KafkaConfig          `json:"kafka_config" yaml:"kafka_config"`
	TopicConfig     TopicConfig          `json:"topic_config" yaml:"topic_config"`
	PublisherConfig KafkaPublisherConfig `json:"publisher_config" yaml:"publisher_config"`
}

type KafkaPublisherConfig struct {
	Topic             string         `json:"topic" yaml:"topic"`
	TopicPrefix       string         `json:"topicPrefix" yaml:"topicPrefix"`
	Partitions        int            `json:"partitions" yaml:"partitions"`
	ReplicationFactor int            `json:"replicationFactor" yaml:"replicationFactor"`
	Headers           []kafka.Header `json:"headers" yaml:"headers"`
}

type S3 struct {
	Endpoint           string        `json:"endpoint" yaml:"endpoint"`
	SecretKey          string        `json:"secret_key" yaml:"secret_key"`
	AccessKey          string        `json:"access_key" yaml:"access_key"`
	Region             string        `json:"region" yaml:"region"`
	PublicBucketName   string        `json:"public_bucket_name" yaml:"public_bucket_name"`
	PrivateBucketName  string        `json:"private_bucket_name" yaml:"private_bucket_name"`
	PublicUrl          string        `json:"public_url" yaml:"public_url"`
	PreSignUrlDuration time.Duration `json:"pre_sign_url_duration" yaml:"pre_sign_url_duration"`
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
	Port string `json:"port" yaml:"port"`
}

type Server struct {
	Name    string `json:"name" yaml:"name"`
	Host    string `json:"host" yaml:"host"`
	GRPC    GRPC
	HTTP    HTTP
	Metrics Metrics
}

type Config struct {
	Server   Server
	Database Database
	Redis    Redis
	JWT      JWT
	Kafka    Kafka
	S3       S3
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
