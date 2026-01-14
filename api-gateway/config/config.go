package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type ServicesConfig struct {
	AuthURL           string
	UserURL           string
	MovieURL          string
	UploadURL         string
	StreamingURL      string
	RecommendationURL string
	NotificationURL   string
}

type JWTConfig struct {
	SecretKey string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.env", "development")
	viper.SetDefault("services.auth_url", "localhost:50051")
	viper.SetDefault("services.user_url", "localhost:50052")
	viper.SetDefault("services.movie_url", "localhost:50053")
	viper.SetDefault("services.upload_url", "localhost:50054")
	viper.SetDefault("services.streaming_url", "localhost:50055")
	viper.SetDefault("services.recommendation_url", "localhost:50056")
	viper.SetDefault("services.notification_url", "localhost:50057")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
			Env:  viper.GetString("server.env"),
		},
		Services: ServicesConfig{
			AuthURL:           viper.GetString("services.auth_url"),
			UserURL:           viper.GetString("services.user_url"),
			MovieURL:          viper.GetString("services.movie_url"),
			UploadURL:         viper.GetString("services.upload_url"),
			StreamingURL:      viper.GetString("services.streaming_url"),
			RecommendationURL: viper.GetString("services.recommendation_url"),
			NotificationURL:   viper.GetString("services.notification_url"),
		},
		JWT: JWTConfig{
			SecretKey: viper.GetString("jwt.secret_key"),
		},
	}

	return config, nil
}
