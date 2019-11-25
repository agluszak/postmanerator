FROM golang:1.13.4

WORKDIR /app

COPY . .

RUN go build -i -v -o postmanerator .

ENTRYPOINT ["./postmanerator"]