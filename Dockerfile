FROM golang:1.15-alpine as builder
WORKDIR /root/go/src/github.com/LILILIhuahuahua/ustc_tencent_game
COPY . /root/go/src/github.com/LILILIhuahuahua/ustc_tencent_game
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
#RUN apk add --no-cache --virtual .build-deps gcc musl-dev
#RUN git config --global url."".insteadOf ""
#RUN export GOPRIVATE=git.enjoymusic.ltd && go build -o bifrost-api main.go plugin.go
RUN go build -o dgs main.go

FROM alpine:latest
ARG DB_HOST
ARG DB_PORT
ENV ENV_HOST=$DB_HOST
ENV ENV_PORT=$DB_PORT
WORKDIR  /root/go/src/github.com/LILILIhuahuahua/ustc_tencent_game
COPY --from=builder  /root/go/src/github.com/LILILIhuahuahua/ustc_tencent_game/dgs .
EXPOSE 8888/udp
ENTRYPOINT ./dgs -DBProxyHost $ENV_HOST -DBProxyPort $ENV_PORT