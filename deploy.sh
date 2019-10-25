#!/bin/sh

set -ex

USER=${ISUCON_USER:-vagrant}
PEM=${PEM:-~/.ssh/isucon9.pem}
IPADDR=${IPADDR:-127.0.0.1}
PORT=${PORT:-22}

BRANCH=master
if [ $# -eq 1 ]; then
  BRANCH=$1
fi

ssh -i $PEM ${USER}@${IPADDR} -p ${PORT} <<EOF

sudo su - isucon
cd ~/isucari  # HOME環境変数はlocalが使われるので注意

# git branch 切り替え
git fetch origin -p
git checkout origin/$BRANCH
git branch
git log --graph --all --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative | head -n 10

# go build

cd ./webapp/go/
make all
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

EOF

#curl -XPOST https://isucon9.catatsuy.org/initialize -H 'Content-Type: application/json' -d '{"payment_service_url":"https://payment.isucon9q.catatsuy.org","shipment_service_url":"https://shipment.isucon9q.catatsuy.org"}'

echo "Success to deploy!"
echo "Done!"
