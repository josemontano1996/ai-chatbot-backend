package config

import (
	"github.com/spf13/viper"
)

type Env struct {
	RedisAddress        string `mapstructure:"REDIS__ADDRESS"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	GeminiApiKey        string `mapstructure:"GEMINI_API_KEY"`
	OpenAiApiKey        string `mapstructure:"OPENAI_API_KEY"`
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
