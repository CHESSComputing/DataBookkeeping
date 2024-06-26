FROM golang:latest as go-builder
MAINTAINER Valentin Kuznetsov vkuznet@gmail.com

# build procedure
ENV PROJECT=DataBookkeeping
ENV WDIR=/data
WORKDIR $WDIR
RUN mkdir -p /build
RUN git clone https://github.com/CHESSComputing/$PROJECT
RUN cd $PROJECT && CGO_ENABLED=1 make && cp srv /build && cp -r static /build

# add default database file
RUN apt-get update && apt-get install sqlite3 && sqlite3 /build/dbs.db "VACUUM;"

# build final image for given image
# FROM alpine as final
# FROM gcr.io/distroless/static as final
# for gibc library we will use debian:stretch
FROM debian:stable-slim
RUN mkdir -p /data
COPY --from=go-builder /build/srv /data
COPY --from=go-builder /build/dbs.db /data
COPY --from=go-builder /build/static /data/static
LABEL org.opencontainers.image.description="FOXDEN DataBookkeeping service"
LABEL org.opencontainers.image.source=https://github.com/chesscomputing/databookkeeping
LABEL org.opencontainers.image.licenses=MIT
WORKDIR /data
