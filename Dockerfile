FROM golang:1.20.0-alpine3.17 as base
RUN apk update
WORKDIR /src/urlshortener
COPY go.mod go.sum ./
COPY . .
RUN go build -o urlshortener ./

FROM alpine:3.17 as binary
WORKDIR /app
COPY --from=base /src/urlshortener/ .
EXPOSE 3000
CMD [ "./urlshortener" ]