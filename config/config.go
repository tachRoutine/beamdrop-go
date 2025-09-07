package config

const PORT = 55555

type Config struct {
	PORT	  int 
}

func GetConfig() Config {
	return Config{
		PORT: PORT,
	}
}