USER_ID=$USER

go mod init
go build -o backitup
mv backitup install/

sudo install -d -o $USER_ID -g $USER_ID -m 0774 -p /var/log/backitup/
sudo install -o $USER_ID -g $USER_ID -m 0774 -p -t /etc/backitup/ install/*

sudo cp system/* /etc/systemd/system/