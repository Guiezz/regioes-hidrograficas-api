package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	ServerPort  string `mapstructure:"PORT"`
}

func LoadConfig() (config Config, err error) {
	// 1. Diz ao Viper para procurar um arquivo chamado ".env"
	viper.SetConfigFile(".env")

	viper.ReadInConfig()

	viper.BindEnv("DATABASE_URL") // <--- Adicione isso
	viper.BindEnv("PORT")         // <--- Adicione isso

	// 3. Define defaults (caso não tenha no .env nem nas variáveis do sistema)
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "user")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "regioes_db")
	viper.SetDefault("PORT", "8080")

	// 4. Lê variáveis de ambiente do sistema (sobrescreve o .env se houver colisão)
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	fmt.Printf("Config carregada. Porta: %s\n", config.ServerPort)
	return
}
