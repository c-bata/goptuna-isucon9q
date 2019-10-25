package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/c-bata/goptuna/rdb"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func replaceEnv(path string, openconns, idleconns, lifetime int) error {
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
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func reload() error {
	return exec.Command("/bin/sh", "-c", "sudo systemctl restart isucari.golang.service").Run()
}

func objective(trial goptuna.Trial) (float64, error) {
	openconns, _ := trial.SuggestInt("OpenConns", 1, 32)
	idleconns, _ := trial.SuggestInt("IdleConns", 1, 32)
	lifetime, _ := trial.SuggestInt("LifetimeSeconds", 1, 32)

	err := replaceEnv("", openconns, idleconns, lifetime)
	if err != nil {
		return 0, err
	}
	err = reload()
	if err != nil {
		return 0, err
	}

	cmd := exec.Command("ls")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	err = cmd.Run()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n")
	line := lines[len(lines)-1]

	// {"pass":true,"score":2010,"campaign":0,"language":"Go","messages":[]}
	var result struct {
		Pass     bool     `json:"pass"`
		Score    int      `json:"score"`
		Campaign int      `json:"campaign"`
		Language string   `json:"language"`
		Messages []string `json:"messages"`
	}
	err = json.Unmarshal([]byte(line), &result)
	if err != nil {
		return 0, err
	}
	return float64(result.Score), nil
}

func main() {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()
	db.DB().SetMaxOpenConns(1)
	storage := rdb.NewStorage(db)

	study, err := goptuna.CreateStudy(
		"isucon9q",
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSampler(tpe.NewSampler()),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMaximize),
	)
	if err != nil {
		log.Fatal("failed to create study:", err)
	}

	err = study.Optimize(objective, 5)
	if err != nil {
		log.Print("optimize catch error:", err)
	}

	v, _ := study.GetBestValue()
	p, _ := study.GetBestParams()
	log.Printf("Best trial: %#v %#v\n", v, p)
}
