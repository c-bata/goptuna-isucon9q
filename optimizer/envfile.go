package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var envPath string

func replaceEnv(openconns, idleconns, lifetime int) error {
	_ = os.Remove(envPath)
	content := fmt.Sprintf(`
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=isucari
MYSQL_DBNAME=isucari
MYSQL_PASS=isucari

MYSQL_MAX_OPEN_CONNECTIONS=%d
MYSQL_MAX_IDLE_CONNECTIONS=%d
MYSQL_MAX_LIFETIME_SECONDS=%d
`, openconns, idleconns, lifetime)
	err := ioutil.WriteFile(envPath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
