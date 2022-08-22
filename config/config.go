package config

import "github.com/spf13/viper"

type IConfig interface {
	DatabaseName() string
	AppName() string
	AppVersion() string
	DabaseDriverName() string
	Port() string
}

type configuration struct {
	dataBaseSource string
	appName        string
	appVersionName string
	databaseDriver string
	port           string
}

func (c *configuration) Port() string {
	return c.port
}

func (c *configuration) DabaseDriverName() string {
	return c.databaseDriver
}

func (c *configuration) DatabaseName() string {
	return c.dataBaseSource
}

func (c *configuration) AppName() string {
	return c.appName
}

func (c *configuration) AppVersion() string {
	return c.appVersionName
}

func NewConfiguration() IConfig {

	viper.SetConfigFile("env.yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &configuration{
		dataBaseSource: viper.GetString("database.source"),
		appName:        viper.GetString("app.name"),
		appVersionName: viper.GetString("app.version"),
		databaseDriver: viper.GetString("database.driver"),
		port:           viper.GetString("server.port"),
	}
}
