package config

import (
	"flag"
	"time"

	"quiz-app-be/internal/config/settings"
)

type (
	Config struct {
		SecretKey  string           `yaml:"secretKey" mapstructure:"SECRET_KEY"`
		DB         Psql             `yaml:"psql" mapstructure:"PSQL"`
		Aws        AwsClient        `yaml:"aws" mapstructure:"AWS"`
		HTTPServer HttpServerConfig `yaml:"httpServer" mapstructure:"HTTPSERVER"`
	}
	Psql struct {
		User         string `yaml:"user" mapstructure:"USER"`
		Pass         string `yaml:"pass" mapstructure:"PASS"`
		Host         string `yaml:"host" mapstructure:"HOST"`
		Port         string `yaml:"port" mapstructure:"PORT"`
		Dbname       string `yaml:"dbname" mapstructure:"DBNAME"`
		Schema       string `yaml:"schema" mapstructure:"SCHEMA"`
		Sslmode      string `yaml:"sslmode" mapstructure:"SSLMODE"`
		RootCert     string `yaml:"rootCert" mapstructure:"ROOTCERT"`
		MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"MAXIDLECONNS"`
		MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"MAXOPENCONNS"`
		OcSQLTrace   bool   `yaml:"ocsql_trace" mapstructure:"OCSQL_TRACE"`
	}
	AwsClient struct {
		AccessId         string `yaml:"accessId" mapstructure:"ACCESS_ID"`
		AccessKey        string `yaml:"accessKey" mapstructure:"ACCESS_KEY"`
		AccessToken      string `yaml:"accessToken" mapstructure:"ACCESS_TOKEN"`
		Endpoint         string `yaml:"endpoint" mapstructure:"ENDPOINT"`
		Region           string `yaml:"region" mapstructure:"REGION"`
		DisableSSL       bool   `yaml:"disableSSL" mapstructure:"DISABLE_SSL"`
		S3ForcePathStyle bool   `yaml:"s3ForcePathStyle" mapstructure:"S3_FORCE_PATH_STYLE"`
		Bucket           string `yaml:"bucket" mapstructure:"BUCKET"`
	}
	HttpServerConfig struct {
		Proto          string        `yaml:"proto" mapstructure:"PROTO"`
		Host           string        `yaml:"host" mapstructure:"HOST"`
		HostOut        string        `yaml:"hostOut" mapstructure:"HOSTOUT"`
		Port           int           `yaml:"port" mapstructure:"PORT"`
		ReadTimeout    time.Duration `yaml:"readTimeout" mapstructure:"READTIMEOUT"`
		WriteTimeout   time.Duration `yaml:"writeTimeout" mapstructure:"WRITETIMEOUT"`
		RequestTimeout time.Duration `yaml:"requestTimeout" mapstructure:"REQUESTTIMEOUT"`
		IdleTimeout    time.Duration `yaml:"idleTimeout" mapstructure:"IDLETIMEOUT"`
		MaxHeaderBytes int           `yaml:"maxHeaderBytes" mapstructure:"MAXHEADERBYTES"`
	}
)

func (c *Config) SetDefault() error {
	if c == nil {
		c = &Config{}
	}
	return nil
}

var config *Config

func Init() (*Config, error) {
	filePath := flag.String("c", "", "Path to configuration file")
	envPrefix := flag.String("e", "", "Environment variable prefix")
	flag.Parse()
	config = &Config{}
	if err := settings.GetDefault(config, *filePath, *envPrefix); err != nil {
		return nil, err
	}
	return config, nil
}
