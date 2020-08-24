FROM golang:1.14-alpine3.12 AS build

WORKDIR /app
COPY ./ /app
RUN go mod download && go build

FROM alpine:3.12

ENV TZ Europe/Moscow
RUN apk add --update --no-cache tzdata  && cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=build /app/entr /app/makaba ./
#USER 1000:1000
ENTRYPOINT echo "/app/makaba" | /app/entr -r /app/makaba