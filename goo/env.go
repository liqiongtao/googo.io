package goo

type Env string

const (
	PRODUCTION  Env = "production"
	SIM         Env = "sim"
	TEST        Env = "test"
	DEVELOPMENT Env = "development"
)

var (
	envTags = map[Env]string{
		PRODUCTION:  "prod",
		SIM:         "sim",
		TEST:        "test",
		DEVELOPMENT: "dev",
	}
)

func (env Env) String() string {
	return string(env)
}

func (env Env) Tag() string {
	return envTags[env]
}
