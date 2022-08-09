package proto

type ConfigSt struct {
	Mysql MysqlSt `json:"mysql" yaml:"mysql"`
	Redis RedisSt `json:"redis" yaml:"redis"`
}

type MysqlSt struct {
	Hostport string `json:"hostport" yaml:"hostport"`
	Password string `json:"password" yaml:"password"`
	Username string `json:"username" yaml:"username"`
	Database string `json:"database" yaml:"database"`

	Poolsize int `json:"poolsize" yaml:"poolsize"`
	Idlesize int `json:"idlesize" yaml:"idlesize"`
}

type RedisSt struct {
	Hostport string `json:"hostport" yaml:"hostport"`
	Auth     string `json:"auth" yaml:"auth"`
	Poolsize int    `json:"poolsize" yaml:"poolsize"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
}
