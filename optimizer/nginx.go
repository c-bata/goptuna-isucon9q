package main

import (
	"html/template"
	"os"
)

const nginxconf = `
user www-data;
worker_processes {{.WorkerProcesses}};
pid /run/nginx.pid;
include /etc/nginx/modules-enabled/*.conf;

error_log  /var/log/nginx/error.log error;

events {
    worker_connections {{.WorkerConnections}};
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    server_tokens off;
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    gzip on;
    keepalive_timeout {{.KeepAliveTimeout}};
    open_file_cache max=1000 inactive=20s;

    client_max_body_size 10m;

    access_log /var/log/nginx/access.log;

    # TLS configuration
    ssl_protocols TLSv1.2;
    ssl_prefer_server_ciphers on;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384';

    include conf.d/*.conf;
    include sites-enabled/*.conf;
}
`

var (
	nginxConfPath string

	nginxTemplate = template.Must(template.New("nginxconf").Parse(nginxconf))
)

type NginxContext struct {
	WorkerProcesses   int
	WorkerConnections int
	KeepAliveTimeout  int
}

func replaceNginxConf(nginxCtx NginxContext) error {
	_ = os.Remove(nginxConfPath)
	f, err := os.Create(nginxConfPath)
	if err != nil {
		return err
	}
	return nginxTemplate.Execute(f, nginxCtx)
}
