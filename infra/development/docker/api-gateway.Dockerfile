FROM alpine:3.22.2
WORKDIR /app

COPY shared shared
COPY build build

ENTRYPOINT ["build/api-gateway"]
