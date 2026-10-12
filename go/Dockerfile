FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /webapi .

FROM alpine:3
COPY --from=build /webapi /webapi
EXPOSE 8080
ENTRYPOINT ["/webapi"]
