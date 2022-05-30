package goo

type Env string

const (
	PRODUCTION  Env = "production"
	SIM         Env = "sim"
	TEST        Env = "test"
	DEVELOPMENT Env = "development"
)

func (env Env) String() string {
	return string(env)
}
