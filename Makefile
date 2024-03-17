export DB_PATH = ./db

$(eval $(file):;@:)

check-migration:
  @[ "${file}" ] || ( echo "\x1b[31;1mERROR: 'file' name is not set\x1b[0m"; exit 1 )

migration-create: check-migration
  goose -dir $(DB_PATH)/migrations -s create $(file) go

run: get
	@echo "------------start service----------------"
	@go run .

get:
	@echo "------------Install Dependencies-------------------"
	@go mod tidy
	@echo "------------Finish Install Dependencies------------"