FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/weather-api .
FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates
COPY --from=build /out/weather-api /usr/local/bin/weather-api
USER app
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/usr/local/bin/weather-api"]
