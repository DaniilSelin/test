run:
	go run cmd/main.go

stop:
	sudo kill -SIGINT $(sudo lsof -ti:8080)