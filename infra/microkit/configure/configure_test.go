package configure_test

type Account struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Mode    string `json:"mode"`
	Port    int64  `json:"port"`
	Version string `json:"version"`
}

type Mysql struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     int64  `json:"port"`
}

type Log struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int64  `json:"max_size"`
	MaxAge     int64  `json:"max_age"`
	MaxBackips int64  `json:"max_backips"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Password string `json:"password"`
	Db       int64  `json:"db"`
	PoolSize int64  `json:"pool_size"`
}

type Email struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
	Port   int64  `json:"port"`
	Rename string `json:"rename"`
}

type Consul struct {
	Host           string `json:"host"`
	Port           int64  `json:"port"`
	Prefix         string `json:"prefix"`
	ConsulRegistry string `json:"consulRegistry"`
}

type Jaeger struct {
	ServiceName string `json:"serviceName"`
	Addr        string `json:"addr"`
}

type Prometheus struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
}

type Ratelimit struct {
	QPS int64 `json:"QPS"`
}

type Micro struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Address string `json:"address"`
}

type ConsulConfig struct {
	Account    Account    `json:"account"`
	Mysql      Mysql      `json:"mysql"`
	Log        Log        `json:"log"`
	Redis      Redis      `json:"redis"`
	Email      Email      `json:"email"`
	Consul     Consul     `json:"consul"`
	Jaeger     Jaeger     `json:"jaeger"`
	Prometheus Prometheus `json:"prometheus"`
	Ratelimit  Ratelimit  `json:"ratelimit"`
	Micro      Micro      `json:"micro"`
}
