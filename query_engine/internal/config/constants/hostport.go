package constants

// HostPortConfig — адрес, на котором query_engine поднимает gRPC-сервер.
type HostPortConfig struct {
	HOST string
	PORT int
}

// InitHostPort читает адрес gRPC-сервера query_engine из окружения.
func InitHostPort() HostPortConfig {
	return HostPortConfig{
		HOST: getStringEnv("QUERY_ENGINE_GRPC_HOST"),
		PORT: int(getIntEnv("QUERY_ENGINE_GRPC_PORT")),
	}
}
