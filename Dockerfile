# builder
FROM golang:1.22.4-alpine as builder

WORKDIR /app

COPY . . 

ENV GO111MODULE=auto
ENV GOPROXY=http://goproxy.cn

RUN go mod tidy && \
	go build -o air-ticket && \
	chmod +x air-ticket 

# runner
FROM alpine:3.19

ARG SERVICE_NAME

WORKDIR /app
COPY --from=builder /app/air-ticket /usr/local/bin/air-ticket 
COPY --from=builder /app/config.yaml /app/conf/config.yaml

ENTRYPOINT ["air-ticket", "-f", "/app/conf/config.yaml"]

