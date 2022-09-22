FROM golang:1.19-alpine

WORKDIR /storage
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY api.go ./api.go
COPY structs.go ./structs.go
COPY testing.go ./testing.go
COPY main.go ./main.go
RUN go mod tidy

RUN go build -o /aplicacion
EXPOSE 10000
CMD [ "/aplicacion" ]

