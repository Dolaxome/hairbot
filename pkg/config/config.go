package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	TelegramToken string `mapstructure:"TOKEN"`
	DBAuth        string `mapstructure:"DBAUTH"`

	Messages Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default   string `mapstructure:"default"`
	NoReq     string `mapstructure:"err-no-req"`
	NoAccess  string `mapstructure:"err-no-access"`
	NoCounter string `mapstructure:"err-count-missing"`
	Unlucky   string `mapstructure:"err-unlucky"`
}

type Responses struct {
	Start          string `mapstructure:"start"`
	UnknownCommand string `mapstructure:"unknown-command"`
	Welcome        string `mapstructure:"welcome"`
	NoAuth         string `mapstructure:"no-auth"`
	AddUser        string `mapstructure:"add-user"`
	MainMenu       string `mapstructure:"main-menu"`
	ChooseThick    string `mapstructure:"choose-thickness"`
	ChooseA        string `mapstructure:"choose-a"`
	ChooseCurl     string `mapstructure:"choose-curl"`
	ChooseDamage   string `mapstructure:"choose-damage"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := ParseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseEnv(cfg *Config) error {
	godotenv.Load()
	cfg.TelegramToken = os.Getenv("TOKEN")
	cfg.DBAuth = os.Getenv("DBAUTH")
	return nil
}
