FROM golang:1.20-buster
WORKDIR /
COPY . .
RUN go build -v -o /app
CMD ["/app"]
