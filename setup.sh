USER_ID=$USER

sudo install -d -o $USER_ID -g $USER_ID -m 0774 -p /var/log/backitup/
sudo install -d -o $USER_ID -g $USER_ID -m 0774 -p /etc/backitup/

go build .

sudo cp system/* /etc/systemd/system/