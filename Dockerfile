
FROM golang:latest as builder

LABEL maintainer="Milan Miljus <miljusmilan44@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# This container exposes port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]