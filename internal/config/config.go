package config

type Config struct {
	ServerAddr   string
	BaseURL      string
	LenQuestions int
}

func NewConfig(addr string, baseURL string, lenQuestions int) Config {
	return Config{
		ServerAddr:   addr,
		BaseURL:      baseURL,
		LenQuestions: lenQuestions,
	}
}
