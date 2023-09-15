# ResServiceBot

## Dependencies
- golang 1.18 or higher (check version in ubuntu ```apt-cache policy golang```)
- git ```apt install git```
- imagemagick ```apt install imagemagick libmagickwand-dev```

## Build backend
go mod tidy
go mod vendor
go build -o /bin/resbot

## App control
Move file appname.service to ```/etc/systemd/system/``` and fix environments. <br />
Enabled service ```systemctl enable resbot.service``` <br />
Start service ```systemctl start resbot.service``` <br />
Check status ```systemctl -l status resbot.service``` <br />
Reload systemd daemon after fixes ```systemctl daemon-reload``` <br />
Show log ```journalctl -u resbot.service --no-pager | tail -10``` <br />
