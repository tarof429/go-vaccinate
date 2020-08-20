default: build

build:
	@go mod download
	@mkdir dist
	@go build -o dist/simulator main.go

test:
	@cd vaccinate; go test

clean:
	@rm -r dist

install:
	@go install