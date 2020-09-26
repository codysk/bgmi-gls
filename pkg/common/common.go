package common

type GLSConfig struct {
	Port string
}

func Init(config *GLSConfig) {
	Config = config
}

var Config *GLSConfig
