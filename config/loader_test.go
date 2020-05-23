package config

import (
	"testing"

	t_assert "github.com/stretchr/testify/assert"
)

type Configuration struct {
	Environment string `yaml:"environment" default:"local"`
	Debug       struct {
		Enable bool `yaml:"enable" default:"false"`
	} `yaml:"debug"`

	Logger struct {
		LogLevel string `yaml:"loglevel" default:"info"`
		Encoding string `yaml:"encoding" default:"console"`
		Sentry   string `yaml:"sentry" default:"dsn"`
	} `yaml:"logger"`

	Core struct {
		AutoMigrate bool   `yaml:"-" default:"false"`
		Host        string `yaml:"host" default:"localhost:5432"`
		Username    string `yaml:"username" default:"golang"`
		Database    string `yaml:"database" default:"accounts"`
		Password    string `yaml:"password" default:""`
	} `yaml:"core"`

	Server struct {
		Listen          string `yaml:"listen" default:":5555"`
		UseTLS          bool   `yaml:"useTLS" default:"false"`
		CertificatePath string `yaml:"certificatePath" default:""`
		PrivateKeyPath  string `yaml:"privateKeyPath" default:""`
	} `yaml:"server"`
}

func TestLoader(t *testing.T) {
	assert := t_assert.New(t)

	conf := &Configuration{}
	err := Load(conf, "BAHN", "")
	assert.Nil(err)

	assert.Equal(conf.Server.Listen, ":5555")
}
