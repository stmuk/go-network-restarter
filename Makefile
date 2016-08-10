all:
	go build pinger.go 
	sudo chown root pinger
	sudo chmod 4755 pinger
	./pinger
