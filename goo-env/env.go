package gooenv

type Env string

const (
	PROD       Env = "prod"
	PRODUCTION Env = "production"

	PRE Env = "pre"
	SIM Env = "sim"

	UAT  Env = "uat"
	TEST Env = "test"

	DEV         Env = "dev"
	DEVELOPMENT Env = "development"
)

var (
	envTags = map[Env]string{
		PROD:       "prod",
		PRODUCTION: "prod",

		PRE: "pre",
		SIM: "sim",

		UAT:  "uat",
		TEST: "test",

		DEV:         "dev",
		DEVELOPMENT: "dev",
	}
)

func (env Env) IsProd() bool {
	return env == PRODUCTION || env == PROD
}

func (env Env) IsPre() bool {
	return env == PRE || env == SIM
}

func (env Env) IsSim() bool {
	return env == PRE || env == SIM
}

func (env Env) IsUat() bool {
	return env == UAT || env == TEST
}

func (env Env) IsTest() bool {
	return env == UAT || env == TEST
}

func (env Env) IsDev() bool {
	return env == DEVELOPMENT || env == DEV
}

func (env Env) String() string {
	return string(env)
}

func (env Env) Tag() string {
	return envTags[env]
}
