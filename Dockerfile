FROM golang:alpine AS build-env
RUN apk --no-cache add git
ADD . /src
RUN cd /src && go build -o http-telegram-notify

FROM alpine
WORKDIR /app
COPY --from=build-env /src/http-telegram-notify /app/
ENTRYPOINT ./http-telegram-notify