#!/bin/bash

systemctl stop goswmon

mkdir -p /var/www/goswmon/templates
cp ./goswmon /var/www/goswmon/
cp ./templates/* /var/www/goswmon/templates/
# cp config.yaml /var/www/goswmon/
cp ./goswmon.service /etc/systemd/system/
chown -R www-data:www-data /var/www/goswmon

systemctl daemon-reload
systemctl start goswmon
