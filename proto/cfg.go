package proto

type ConfigSt struct {
	Mysql   MysqlSt   `json:"mysql" yaml:"mysql"`
	Redis   RedisSt   `json:"redis" yaml:"redis"`
	LogType LogTypeSt `json:"logger" yaml:"logger"`
	Kafka   KafkaSt   `json:"kafka" yaml:"kafka"`
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

type LogTypeSt struct {
	StdoutLevel   string `json:"stdout.level" yaml:"stdout.level"`
	StdoutEnabled bool   `json:"stdout.enabled" yaml:"stdout.enabled"`

	FileLevel    string `json:"file.level" yaml:"file.level"`
	FileEnabled  bool   `json:"file.enabled" yaml:"file.enabled"`
	FileFilename string `json:"file.filename" yaml:"file.filename"`
}

type KafkaSt struct {
	Hostports  string `json:"hostports" yaml:"hostports"`
	Concurrent int64  `json:"concurrent" yaml:"concurrent"`
	Enable     bool   `json:"enable" yaml:"enable"`
}
