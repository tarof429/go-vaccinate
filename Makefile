default: build

build:
	go mod download
	go build -o go-triage

clean:
	rm -f go-triage

install:
	go install