package goowechat

type Config struct {
	Appid        string `yaml:"appid"`
	Secret       string `yaml:"secret"`
	H5Url        string `yaml:"h5_url"`
	AuthorizeUrl string `yaml:"authorize_url"`
}
