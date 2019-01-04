# Builder image
FROM golang:alpine as builder

RUN apk add --no-cache git build-base

RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get -d -v ./...
RUN go build -o alexandria -v .

# Server image
FROM alpine

RUN adduser -S -D -H -h /app appuser
USER appuser

WORKDIR /app
COPY --from=builder /build/alexandria .
COPY --from=builder /build/static static
COPY --from=builder /build/templates templates
CMD ["./alexandria", "-prepared-keyword-store=/tmp/alexandria.db"]
