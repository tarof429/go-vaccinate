# Build the code
FROM golang:1.19 as builder

RUN mkdir -p /src

WORKDIR /src

COPY . .

RUN make clean build

# Start a new stage to provide the runtime environment
FROM alpine:latest

# Install a library needed by code executables
RUN apk add --no-cache libc6-compat

WORKDIR /root

COPY --from=builder /src/profiles/gastrorhino /root/.vaccinate
COPY --from=builder /src/dist/simulator /root

CMD ["./simulator", "--terminal"]