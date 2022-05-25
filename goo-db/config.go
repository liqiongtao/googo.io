package goo_db

type Config struct {
	Name        string   `yaml:"name"`
	Driver      string   `yaml:"driver"`
	Master      string   `yaml:"master"`
	Slaves      []string `yaml:"slaves"`
	MaxIdle     int      `yaml:"max_idle"`
	MaxOpen     int      `yaml:"max_open"`
	AutoPing    bool     `yaml:"auto_ping"`
	LogModel    bool     `yaml:"log_model"`
	LogFilepath string   `yaml:"log_filepath"`
}
