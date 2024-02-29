package goo_db

type Config struct {
	Name        string   `yaml:"name" json:"name"`
	Driver      string   `yaml:"driver" json:"driver"`
	Master      string   `yaml:"master" json:"master"`
	Slaves      []string `yaml:"slaves" json:"slaves"`
	MaxIdle     int      `yaml:"max_idle" json:"max_idle"`
	MaxOpen     int      `yaml:"max_open" json:"max_open"`
	MaxLifetime int      `yaml:"max_lifetime" json:"max_lifetime"` // 最大生命周期，单位秒
	AutoPing    bool     `yaml:"auto_ping" json:"auto_ping"`
	LogModel    bool     `yaml:"log_model" json:"log_model"`
	LogFilepath string   `yaml:"log_filepath" json:"log_filepath"`
}
