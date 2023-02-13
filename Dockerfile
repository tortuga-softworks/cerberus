FROM golang:1.20-alpine

WORKDIR /cerberus-src

# Downloading dependencies
COPY . .
RUN go mod download

# Building the application
RUN go build -o /cerberus

CMD [ "/cerberus" ]
