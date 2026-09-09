package goo_mongo

type Config struct {
	Name       string `yaml:"name" json:"name"`
	Addr       string `yaml:"addr" json:"addr"`
	User       string `yaml:"user" json:"user"`
	Password   string `yaml:"password" json:"password"`
	Database   string `yaml:"database" json:"database"`
	AuthSource string `yaml:"auth_source" json:"authSource"` // 认证库，默认 admin
	Timeout    int    `yaml:"timeout" json:"timeout"`         // Connect/Ping 超时秒数，默认 10
	AutoPing   bool   `yaml:"auto_ping" json:"autoPing"`
}
