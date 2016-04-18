
all: go-swagger
	@echo "Build swaggerly..."; \
	go get; \
	go build


# Checkout correct branch of go-swagger
go-swagger:
	@echo "Get go-swagger and switch to swaggerly-versioning-extension branch..."; \
	go get github.com/zxchris/go-swagger; \
	cd ../go-swagger; \
	git checkout feature/swaggerly-versioning-extension
