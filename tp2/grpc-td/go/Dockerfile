FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./server && CGO_ENABLED=0 go build -o /out/client ./client

FROM alpine:3
COPY --from=build /out/ /usr/local/bin/
EXPOSE 50051
CMD ["server"]
