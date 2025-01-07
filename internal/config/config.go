package config

type Config struct {
	Telegram struct {
		Token string
	}
	Convert struct {
		StoragePath string `mapstructure:"storage_path"`
	}
}
