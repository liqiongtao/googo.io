package goo_clickhouse

type Config struct {
	Driver       string `yaml:"driver"`
	Addr         string `yaml:"addr"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	ReadTimeout  int32  `yaml:"read_timeout"`
	WriteTimeout int32  `yaml:"write_timeout"`
	AltHosts     string `yaml:"alt_hosts"`
	AutoPing     bool   `yaml:"auto_ping"`
	Debug        bool   `yaml:"debug"`
}
