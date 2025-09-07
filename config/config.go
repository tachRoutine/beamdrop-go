package config

const PORT = 7777

type Config struct {
	PORT int
}

func GetConfig() Config {
	return Config{
		PORT: PORT,
	}
}
