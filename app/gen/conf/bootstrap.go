package conf

type Bootstrap struct {
	DbType   string
	Dsn      string
	Env      ENVDEF
	Addr     string
	LogFile  string
	LogLevel string
}
