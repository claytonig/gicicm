# base image
FROM golang:alpine as build

RUN apk update && apk add \
    curl \
    make

RUN rm -rf /var/lib/apt/lists/*

RUN mkdir -p /srv/app
WORKDIR /srv/app

COPY . /srv/app

RUN make mod

RUN make buildl

# ---> build binaries
FROM alpine
COPY --from=build /srv/app/go-server /app/
WORKDIR /app

CMD ["./go-icm"]


