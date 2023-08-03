.PHONY:
.SILENT:
.DEFAULT_GOAL := run

run: 
	docker-compose up --remove-orphans app

lint: 
	golangci-lint run


 