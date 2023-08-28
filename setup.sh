#!/bin/sh
USER_ID=$USER

go mod download
go build -o encrypted-go-backup
mv encrypted-go-backup install/

sudo install -d -o $USER_ID -g $USER_ID -m 0774 -p /var/log/encrypted-go-backup/
sudo install -d -o $USER_ID -g $USER_ID -m 0774 -p /etc/encrypted-go-backup/ 

cp install/* /etc/encrypted-go-backup/

sed "s/%i/$USER/" system/timer.template > system/encrypted-go-backup.timer

sudo cp system/encrypted-go-backup* /etc/systemd/system/

rm system/encrypted-go-backup.timer

systemctl start encrypted-go-backup@$USER_ID.service
systemctl enable encrypted-go-backup.timer