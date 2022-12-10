run: stop
	docker-compose up 

run-detach: stop
	docker-compose up -d

stop:
	docker-compose down
