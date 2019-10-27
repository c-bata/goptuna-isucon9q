package main

import (
	"html/template"
	"os"
)

const envfile = `
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=isucari
MYSQL_DBNAME=isucari
MYSQL_PASS=isucari

MYSQL_MAX_OPEN_CONNECTIONS={{.MaxOpenConns}}
MYSQL_MAX_IDLE_CONNECTIONS={{.MaxIdleConns}}
MYSQL_MAX_LIFETIME_SECONDS={{.MaxLifetimeSeconds}}

ISUCARI_CAMPAIGN={{.Campaign}}
`

var (
	envPath string

	envTemplate = template.Must(template.New("envfile").Parse(envfile))
)

type EnvfileContext struct {
	MaxOpenConns       int
	MaxIdleConns       int
	MaxLifetimeSeconds int
	Campaign           int
}

func replaceEnv(envfileCtx EnvfileContext) error {
	_ = os.Remove(envPath)
	f, err := os.Create(envPath)
	if err != nil {
		return err
	}
	return envTemplate.Execute(f, envfileCtx)
}
