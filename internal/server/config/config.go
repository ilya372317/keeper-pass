package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const defaultTokenExpInHours = 6

type GRPCConfig struct {
	Host string
}

type JWTConfig struct {
	SecretKey        string
	TokenExpDuration time.Duration
}

type SQLConfig struct {
	Host              string
	Port              string
	DBName            string
	User              string
	Password          string
	Charset           string
	Timezone          string
	Collation         string
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ConnMaxLifetime   time.Duration
	MaxOpenConns      int
	MaxIdleConns      int
	InterpolateParams bool
	ParseTime         bool
}

type Config struct {
	GRPC   GRPCConfig
	JWT    JWTConfig
	MainDB SQLConfig
}

func New(configPath string) (Config, error) {
	var cnfg Config
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed load config file. file path [%s] is invalid: %w", configPath, err)
	}

	viper.SetDefault("grpc.host", ":3200")
	viper.SetDefault("jwt.secret_key", "")
	viper.SetDefault("jwt.token_exp", time.Hour*defaultTokenExpInHours)

	cnfg = Config{
		GRPC: GRPCConfig{Host: viper.GetString("grpc.host")},
		JWT: JWTConfig{
			SecretKey:        viper.GetString("jwt.secret_key"),
			TokenExpDuration: viper.GetDuration("jwt.token_exp"),
		},
		MainDB: SQLConfig{
			Host:              viper.GetString("main_db.host"),
			Port:              viper.GetString("main_db.port"),
			DBName:            viper.GetString("main_db.name"),
			User:              viper.GetString("main_db.user"),
			Password:          viper.GetString("main_db.password"),
			ReadTimeout:       viper.GetDuration("main_db.read_timeout"),
			WriteTimeout:      viper.GetDuration("main_db.write_timeout"),
			Timeout:           viper.GetDuration("main_db.timeout"),
			InterpolateParams: false,
			Charset:           "UTF-8",
			ParseTime:         true,
			Timezone:          "Europe/Moscow",
			Collation:         "",
		},
	}

	return cnfg, nil
}
