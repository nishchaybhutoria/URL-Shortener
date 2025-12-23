.PHONY: build run stop logs clean reset

build:
	docker-compose build

run:
	docker-compose up -d --build

stop:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down --remove-orphans

reset:
	docker-compose down -v --remove-orphans
	sudo rm -rf ./data/*