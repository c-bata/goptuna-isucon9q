#!/bin/sh

set -ex

cd ~/isucari

cd ./webapp/go/
make all
cd ~/isucari

cd ./optimizer/
go build .
cd ~/isucari

# log rotate
if [ -f /var/log/nginx/access.log ]; then
     sudo mv /var/log/nginx/access.log /var/log/nginx/access.log.$(date +"%Y%m%d_%H%M%S")
fi
if [ -f /var/log/mysql/mysql-slow.log ]; then
    sudo mv /var/log/mysql/mysql-slow.log /var/log/mysql/mysql-slow.log.$(date +"%Y%m%d_%H%M%S")
fi

# service restart
sudo systemctl restart isucari.golang.service
sudo systemctl restart nginx.service
sudo systemctl restart mysql.service
echo "Success to restart!"

# mysql setting
# echo "set global slow_query_log = ON;" | sudo mysql -u root
# echo "set global slow_query_log_file = '/var/log/mysql/mysql-slow.log';" | sudo mysql -u root
# echo "set global long_query_time = 0;" | sudo mysql -u root
# echo "Success slow_query_log set!"

# init private dir
rm -rf /home/isucon/isucari/webapp/private
mkdir -p /home/isucon/isucari/webapp/private/qrcode

echo "Success to reload!"

