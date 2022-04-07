package config

import "github.com/kelseyhightower/envconfig"

var Cfg Config

type Config struct {
	Debug  bool   `envconfig:"DEBUG" default:"true"`
	JwtKey []byte `envconfig:"JWT_KEY" default:"super secret"`
	// Postgres config
	DbName string `envconfig:"DB_NAME" default:"reddit"`
	DbHost string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DbPort string `envconfig:"DB_PORT" default:"55436"`
	DbUser string `envconfig:"DB_USER" default:"reddit"`
	DbPass string `envconfig:"DB_PASS" default:"reddit"`
	// Redis config
	RedisHost string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort string `envconfig:"REDIS_PORT" default:"63790"`
	RedisDb   int    `envconfig:"REDIS_DB" default:"0"`
	RedisPwd  string `envconfig:"REDIS_PWD" default:""`
	// Mongo config
	MongoHost string `envconfig:"MONGO_HOST" default:"127.0.0.1"`
	MongoPort string `envconfig:"MONGO_PORT" default:"27017"`
}

func Load() {
	envconfig.MustProcess("", &Cfg)
}
