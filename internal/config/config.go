package config

import (
	"github.com/spf13/viper"
)

type Env struct {
	ServerPort string `mapstructure:"SERVER_PORT"`

	RedisAddress  string `mapstructure:"REDIS__ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	GeminiApiKey    string `mapstructure:"GEMINI_API_KEY"`
	GeminiMaxTokens uint32 `mapstructure:"GEMINI_MAX_TOKENS"`

	OpenAiApiKey    string `mapstructure:"OPENAI_API_KEY"`
	OpenAiMaxTokens uint32 `mapstructure:"OPENAI_MAX_TOKENS"`

	MaxCompletionTokens uint32 `mapstructure:"MAX_COMPLETION_TOKENS"`
}

func LoadEnv(path string, fileName string) (env Env, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&env)
	return
}
