run: get
	@echo "------------start service----------------"
	@go run .

get:
	@echo "------------Install Dependencies-------------------"
	@go mod tidy
	@echo "------------Finish Install Dependencies------------"