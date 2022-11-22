FROM golang:1.18-alpine AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /src/bin/go-user-service .


FROM alpine:3.9 AS go-user-service
RUN apk add ca-certificates
COPY --from=build /src/bin/go-user-service /bin/go-user-service

EXPOSE 80
WORKDIR ~
CMD ["go-user-service"]