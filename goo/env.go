package goo

type Env string

const (
	PRODUCTION  Env = "prod"
	SIM         Env = "sim"
	TEST        Env = "test"
	DEVELOPMENT Env = "dev"
)

func (env Env) String() string {
	return string(env)
}
