USER_ID=$USER

go mod init
go build -o encrypted-go-backup
mv encrypted-go-backup install/

sudo install -d -o $USER_ID -g $USER_ID -m 0774 -p /var/log/encrypted-go-backup/
sudo install -o $USER_ID -g $USER_ID -m 0774 -p -t /etc/encrypted-go-backup/ install/*

sudo cp system/* /etc/systemd/system/