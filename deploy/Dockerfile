FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /out/alertd ./cmd/alertd
FROM alpine:3.20
RUN adduser -D -u 10001 app
COPY --from=build /out/alertd /usr/local/bin/alertd
USER app
EXPOSE 8092
ENTRYPOINT ["/usr/local/bin/alertd"]
