PG = postgresql
MEMORY = memory

all:
		@echo "Specify storage type: make postgresql or make memory"

test:
		@echo "Running tests...."
		@go test ./internal/config_parser
		@go test ./internal/tools

$(PG):
		docker-compose --profile postgresql up --build


$(MEMORY):
		docker-compose --profile memory up --build