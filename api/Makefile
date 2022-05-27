-include .env

SHA=`git rev-parse HEAD`

check-prod:
	@if [ "$(GCLOUD_PROJECT)" = "vama-prod" ]; then \
	    read -p "YOU ARE BUILDING PROD. ENTER 'prod' TO ACKNOWLEDGE:" userinput; \
	    if [ "$$userinput" != "prod" ]; then \
			echo "Incorrect input."; exit 1; \
		fi \
	fi

build: check-prod
	go build -o api cmd/api/*
	go build cmd/test/main.go

run-dev:
	go run cmd/api/*

test:
	./cmd/test/runtests.sh $(entity)

health:
	rm -f -- health.zip
	cd health; zip ../health.zip *
	cp health.zip ./terraform/environments/dev/health.zip
	cp health.zip ./terraform/environments/staging/health.zip
	cp health.zip ./terraform/environments/prod/health.zip
	rm -f -- health.zip
.PHONY: build health