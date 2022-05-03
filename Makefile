.PHONY: all build run gotool clean help

APP=crontab

build: clean
	go build -o ${APP} ./app/worker/main.go

master-run:
	go run -race ./app/master/main.go

worker-run:
	go run -race ./app/worker/main.go

clean:
	go clean