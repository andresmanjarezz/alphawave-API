# .PHONY:
# .SILENT:
.DEFAULT_GOAL := run
.SHELL := /bin/bash

run: 
	docker compose up --remove-orphans app

lint: 
	golangci-lint run


 