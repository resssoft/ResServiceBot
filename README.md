# ResServiceBot

## Dependencies
- golang 1.18 or higher

## Build backend
go build -o /bin/appname cmd/appname/main.go

## App control
Move file appname.service to ```/etc/systemd/system/``` and fix environments. <br />
Enabled service ```systemctl enable appname.service``` <br />
Start service ```systemctl start appname.service``` <br />
Check status ```systemctl -l status appname.service``` <br />
Reload systemd daemon after fixes ```systemctl daemon-reload``` <br />
Show log ```journalctl -u appname.service --no-pager | tail -10``` <br />
