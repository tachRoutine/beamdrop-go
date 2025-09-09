package config

const PORT = 7777

type Config struct {
	PORT int
}

type Flags struct {
	SharedDir string
	NoQR      bool
	Help      bool
}

func GetConfig() Config {
	return Config{
		PORT: PORT,
	}
}
