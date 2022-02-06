package config

import (
	"github.com/rome314/3-commas-ws-forwarder/internal/sender"
	"github.com/rome314/3-commas-ws-forwarder/pkg/connections"
	"github.com/spf13/viper"
)

func init() {
	viperInit()
	cfg = &Config{
		Debug: viper.GetBool("DEBUG"),
		Redis: connections.RedisConfig{
			Address:  viper.GetString("REDIS_ADDRESS"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Db:       viper.GetInt("REDIS_DB"),
		},
		App: sender.Config{
			Topic:         viper.GetString("APP_TOPIC"),
			Url:           viper.GetString("WEBHOOK_URL"),
			ConsumerGroup: viper.GetString("APP_CONSUMER_GROUP"),
		},
		Api: ApiConfig{
			Key:    viper.GetString("API_KEY"),
			Secret: viper.GetString("API_SECRET"),
		},
	}
}

func viperInit() {

	viper.SetDefault("APP_TOPIC", "events")
	viper.SetDefault("APP_CONSUMER_GROUP", "events")

	viper.AutomaticEnv()

}
