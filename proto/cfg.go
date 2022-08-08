package proto

type ConfigSt struct {
	Mysql MysqlSt `json:"mysql" yaml:"mysql"`
}

type MysqlSt struct {
	Hostport string `json:"hostport" yaml:"hostport"`
	Password string `json:"password" yaml:"password"`
	Username string `json:"username" yaml:"username"`
	Database string `json:"database" yaml:"database"`

	Poolsize int `json:"poolsize" yaml:"poolsize"`
	Idlesize int `json:"idlesize" yaml:"idlesize"`
}
