FROM golang:1.20-alpine

WORKDIR /cerberus-src

# Downloading dependencies
COPY . .
RUN go mod download

# Building the application
RUN go build -o /cerberus ./cmd

# Using a new base image to run the binary
FROM alpine:latest  

WORKDIR /root/

COPY --from=0 /cerberus /cerberus

CMD [ "/cerberus" ]
