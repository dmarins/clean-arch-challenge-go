DOCKERCOMPOSECMD=docker-compose
WIRECMD=wire
GOCMD=go

.PHONY: wire

wire:
	@command -v wire >/dev/null 2>&1 || $(GOCMD) install github.com/google/wire/cmd/wire@latest
	@cd cmd/ordersystem && $(WIRECMD)

dc-up:
	$(DOCKERCOMPOSECMD) up -d --force-recreate
	@echo "Waiting until Mysql be ready..."
	@until docker ps | grep mysql | grep "(healthy)"; do sleep 1; done
	@echo "Mysql is started."

dc-down:
	docker-compose down --remove-orphans

dc-restart: dc-down dc-up

db-init:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

db-drop:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

db-query:
	docker exec -it mysql mysql -uroot -proot orders

run:
	cd cmd/ordersystem/ && $(GOCMD) run main.go wire_gen.go