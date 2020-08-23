default: build

build:
	@go mod download
	@mkdir -p dist
	@go build -o dist/simulator main.go

test:
	@cd vaccinate; go test

clean:
	@rm -rf dist

install:
	@go install