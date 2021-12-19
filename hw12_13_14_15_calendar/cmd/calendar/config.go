package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Env  string
	File string
}

const (
	STORAGE_MEMORY = "memory"
	STORAGE_SQL    = "sql"
)

type StorageConf struct {
	Type string
	Dsn  string
}

func NewConfig() Config {
	return Config{}
}

// TODO
