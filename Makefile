default: build

build:
	go mod download
	go build -o go-triage

test:
	(cd vaccinate; go test)

clean:
	rm -f go-triage

install:
	go install