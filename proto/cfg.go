package proto

type ConfigSt struct {
	Mysql       MysqlSt             `json:"mysql" yaml:"mysql"`
	Redis       RedisSt             `json:"redis" yaml:"redis"`
	LogType     LogTypeSt           `json:"logger" yaml:"logger"`
	KafkaWriter KafkaWriterConfigSt `json:"kafka_writer" yaml:"kafka_writer"`
	KafkaReader KafkaReaderConfigSt `json:"kafka_reader" yaml:"kafka_reader"`
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

type KafkaWriterConfigSt struct {
	Hostports  string `json:"hostports" yaml:"hostports"`
	Concurrent int64  `json:"concurrent" yaml:"concurrent"`
	Enable     bool   `json:"enable" yaml:"enable"`
}

type KafkaReaderConfigSt struct {
	Enable     bool   `json:"enable" yaml:"enable"`
	Hostports  string `json:"hostports" yaml:"hostports"`
	Topic      string `json:"topic" yaml:"topic"`
	GroupId    string `json:"group_id" yaml:"group_id"`
	MsgDecoder string `json:"msg_decoder" yaml:"msg_decoder"`
}
