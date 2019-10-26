package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/goptuna/rdb"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func reload() error {
	if err := exec.Command("/bin/sh", "-c", "sudo systemctl restart isucari.golang.service").Run(); err != nil {
		return fmt.Errorf("failed to reload app: %s", err)
	}
	if err := exec.Command("/bin/sh", "-c", "sudo systemctl restart nginx.service").Run(); err != nil {
		return fmt.Errorf("failed to reload nginx: %s", err)
	}
	if err := exec.Command("/bin/sh", "-c", "sudo systemctl restart mysql.service").Run(); err != nil {
		return fmt.Errorf("failed to relaod mysql: %s", err)
	}
	return nil
}

func bench() (int, error) {
	cmd := exec.Command("./bin/benchmarker")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n")
	line := lines[len(lines)-1]
	log.Println("line:", line)

	// {"pass":true,"score":2010,"campaign":0,"language":"Go","messages":[]}
	var result struct {
		Pass     bool     `json:"pass"`
		Score    int      `json:"score"`
		Campaign int      `json:"campaign"`
		Language string   `json:"language"`
		Messages []string `json:"messages"`
	}

	if err := json.Unmarshal([]byte(line), &result); err != nil {
		return 0, err
	}
	return result.Score, nil
}

func objective(trial goptuna.Trial) (float64, error) {
	// go application
	goMySQLOpenConns, _ := trial.SuggestInt("mysql_client_open_conns", 1, 32)
	goMySQLIdleConns, _ := trial.SuggestInt("mysql_client_idle_conns", 1, 32)
	goMySQLMaxLifetime, _ := trial.SuggestInt("mysql_client_max_lifetime", 1, 64)
	goMySQLHttpIdleConnsPerHost, _ := trial.SuggestInt("http_max_idle_conns_per_host", 1, 2048)
	campaign, _ := trial.SuggestInt("campaign", 0, 4)
	if err := replaceEnv(goMySQLOpenConns, goMySQLIdleConns, goMySQLMaxLifetime, goMySQLHttpIdleConnsPerHost, campaign); err != nil {
		return 0, err
	}

	// nginx
	nginxWorkerProcesses, _ := trial.SuggestInt("nginx_worker_processes", 1, 16)
	nginxWorkerConns, _ := trial.SuggestInt("nginx_worker_connections", 1, 4096)
	nginxKeepAliveTimeout, _ := trial.SuggestInt("nginx_keep_alive_timeout", 1, 100)
	nginxOpenFileCacheMax, _ := trial.SuggestInt("nginx_open_file_cache_max", 100, 10000)
	nginxOpenFileCacheInActive, _ := trial.SuggestInt("nginx_open_file_cache_inactive", 1, 64)
	nginxGzip, _ := trial.SuggestCategorical("nginx_gzip", []string{"on", "off"})
	if err := replaceNginxConf(
		nginxWorkerProcesses,
		nginxWorkerConns,
		nginxKeepAliveTimeout,
		nginxOpenFileCacheMax,
		nginxOpenFileCacheInActive,
		nginxGzip,
	); err != nil {
		return 0, err
	}

	// mysql
	innoDBBufferPoolSize, _ := trial.SuggestInt("innodb_buffer_pool_size", 10, 800)                                     // default 128MB
	innoDBLogBufferSize, _ := trial.SuggestInt("innodb_log_buffer_size", 1, 64)                                         // default 8MB or 16MB
	innoDBLogFileSize, _ := trial.SuggestInt("innodb_log_file_size", 10, 1024)                                          // default 48MB
	innoDBFlushLogAtTRXCommit, _ := trial.SuggestCategorical("innodb_flush_log_at_trx_commit", []string{"0", "1", "2"}) // default 1
	innodbFlushMethod, _ := trial.SuggestCategorical("innodb_flush_method", []string{
		"fsync",
		"littlesync",
		"nosync",
		"O_DIRECT",
		"O_DIRECT_NO_FSYNC",
	})
	if err := replaceMySQLConf(innoDBBufferPoolSize, innoDBLogBufferSize, innoDBLogFileSize, innoDBFlushLogAtTRXCommit, innodbFlushMethod); err != nil {
		return 0, err
	}

	if err := reload(); err != nil {
		return 0, err
	}
	score, err := bench()
	if err != nil {
		return 0, err
	}
	return float64(score), nil
}

func main() {
	flag.StringVar(&envPath, "envfile", "/home/isucon/env.sh", "filepath to env")
	flag.StringVar(&mysqlConfPath, "mysqlcnf", "/home/isucon/isucari/etc/mysql/mysqld.cnf", "filepath to mysql conf")
	flag.StringVar(&nginxConfPath, "nginxcnf", "/home/isucon/isucari/etc/nginx/nginx.conf", "filepath to nginx conf")
	flag.Parse()

	_ = os.Remove("db.sqlite3")
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()
	db.DB().SetMaxOpenConns(1)
	rdb.RunAutoMigrate(db)
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

	err = study.Optimize(objective, 250)
	if err != nil {
		log.Print("optimize catch error:", err)
	}

	v, _ := study.GetBestValue()
	p, _ := study.GetBestParams()
	log.Printf("Best trial: %#v %#v\n", v, p)
}
