package configs

import "strings"

type Env string

func (e Env) IsProd() bool {
	return strings.ToLower(string(e)) == "production" || strings.ToLower(string(e)) == "prod"
}

func (e Env) IsDev() bool {
	return e == "development"
}

func (e Env) IsLocal() bool {
	return e == "local"
}

func (e Env) String() string {
	return string(e)
}
