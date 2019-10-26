package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var envPath string

func replaceEnv(openconns, idleconns, lifetime, httpIdleConnsPerHost, campaign int) error {
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

HTTP_MAX_IDLE_CONNS_PER_HOST=%d
ISUCARI_CAMPAIGN=%d
`, openconns, idleconns, lifetime, httpIdleConnsPerHost, campaign)
	err := ioutil.WriteFile(envPath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
