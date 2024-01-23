package goo_es

type Config struct {
	Addresses []string `json:"addresses" yaml:"addresses"`
	User      string   `json:"user" yaml:"user"`
	Password  string   `json:"password" yaml:"password"`
	EnableLog bool     `json:"enable_log" yaml:"enable_log"`
}
