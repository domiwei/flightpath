run: stop
	docker-compose up 

run-detach: stop
	docker-compose up -d

build:
	docker-compose build

stop:
	docker-compose down
