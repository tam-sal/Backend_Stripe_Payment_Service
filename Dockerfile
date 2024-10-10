FROM golang:1.23-alpine3.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o server .
EXPOSE 4242
CMD ["./server"]