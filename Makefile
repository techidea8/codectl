build:
	rsrc -ico app.ico -o app.syso
	go build -o spider-server.exe -ldflags "-w -s -H windowsgui"
	go build -o spider-server-console.exe