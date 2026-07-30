package constants

type HostPortConfig struct {
	HOST string
	PORT int
}

func InitHostPort() HostPortConfig {
	return HostPortConfig{
		HOST: getStringEnv("DATACATALOGUE_GRPC_HOST"),
		PORT: int(getIntEnv("DATACATALOGUE_GRPC_PORT")),
	}
}
