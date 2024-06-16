package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type flags struct {
	serverConfigPath   string
	mailServerPassword string
	TgBotKey           string
}

type Config struct {
	Env         string `yaml:"env" envDefault:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	TgBotKey    string
	HTTPServer  `yaml:"http_server"`
	MailServer  `yaml:"mail_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" envDefault:"localhost:2001"`
	Timeout     time.Duration `yaml:"timeout" envDefault:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" envDefault:"60s"`
}

type MailServer struct {
	Host     string `yaml:"host" envDefault:"smtp.gmail.com"`
	Port     int    `yaml:"port" envDefault:"587"`
	Sender   string `yaml:"sender"`
	Password string
}

func MustLoad() *Config {
	e := fetchFlags()

	if e.serverConfigPath == "" {
		panic("server config path is empty")
	}

	if e.TgBotKey == "" {
		panic("telegram bot api key is empty")
	}

	if e.mailServerPassword == "" {
		panic("mail server password is empty")
	}

	return MustLoadByPath(e)
}

func MustLoadByPath(e flags) *Config {
	if _, err := os.Stat(e.serverConfigPath); os.IsNotExist(err) {
		panic("config file does not exist: " + e.serverConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(e.serverConfigPath, &cfg); err != nil {
		panic("failed to read config" + err.Error())
	}

	cfg.TgBotKey = e.TgBotKey
	cfg.MailServer.Password = e.mailServerPassword

	return &cfg
}

func fetchFlags() flags {
	var serverConfigPath string
	var tgBotKey string
	var mailServerPassword string

	flag.StringVar(&serverConfigPath, "config", "", "path to server config file")
	flag.StringVar(&tgBotKey, "tg-bot-key", "", "telegram bot api key")
	flag.StringVar(&mailServerPassword, "mail-passw", "", "mail server password")
	flag.Parse()

	if serverConfigPath == "" {
		serverConfigPath = os.Getenv("CONFIG_PATH")
	}

	if tgBotKey == "" {
		tgBotKey = os.Getenv("TG_BOT_KEY")
	}

	if mailServerPassword == "" {
		mailServerPassword = os.Getenv("MAIL_PASSW")
	}

	return flags{
		serverConfigPath:   serverConfigPath,
		TgBotKey:           tgBotKey,
		mailServerPassword: mailServerPassword,
	}
}
