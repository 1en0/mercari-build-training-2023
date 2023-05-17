FROM alpine

ENV CGO_ENABLED=1

RUN apk add --no-cache go
RUN apk update && apk add --no-cache curl
RUN apk update && apk add --no-cache build-base

RUN go version
#RUN addgroup -S mercari && adduser -S trainee -G mercari
## RUN chown -R trainee:mercari /path/to/db
#
#USER trainee
#
#CMD ["go", "version"]

WORKDIR /app/
COPY go.mod ./
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/go-sample-app ./app

EXPOSE 8080

ENTRYPOINT ["./out/go-sample-app"]
