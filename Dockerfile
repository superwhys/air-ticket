# builder
FROM golang:alpine3.18 as builder

WORKDIR /app

COPY ./* /app/

ENV GO111MODULE=auto
ENV GOPROXY=http://goproxy.cn

RUN go mod tidy && \
	go build -o air-ticket && \
	chmod +x air-ticket 

# runner
FROM registry.cn-shenzhen.aliyuncs.com/fdz/golang-alpine-chrome:v1.0.0

ARG SERVICE_NAME

WORKDIR /app
COPY --from=builder /app/air-ticket /usr/local/bin/air-ticket 
COPY --from=builder /app/config.yaml /app/config.yaml

ENTRYPOINT ["air-ticket", "-f", "/app/config.yaml"]

