FROM docker.io/library/golang:1.23 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main .

FROM scratch

COPY --from=build /app/main /app/main

CMD ["/app/main"]