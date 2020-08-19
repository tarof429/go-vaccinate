default: build

build:
	go mod download
	go build -o go-vaccine

test:
	(cd vaccinate; go test)

clean:
	rm -f go-vaccine

install:
	go install