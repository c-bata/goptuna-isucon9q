package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var nginxConfPath string

func replaceNginxConf(numOfWorkers int) error {
	_ = os.Remove(nginxConfPath)
	content := fmt.Sprintf(`
user www-data;
worker_processes %d;
pid /run/nginx.pid;
include /etc/nginx/modules-enabled/*.conf;

error_log  /var/log/nginx/error.log error;

events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    server_tokens off;
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 120;
    client_max_body_size 10m;

    access_log /var/log/nginx/access.log;

    # TLS configuration
    ssl_protocols TLSv1.2;
    ssl_prefer_server_ciphers on;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384';

    include conf.d/*.conf;
    include sites-enabled/*.conf;
}
`, numOfWorkers)
	err := ioutil.WriteFile(nginxConfPath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
