## Build
FROM golang:1.19.3-bullseye AS build

WORKDIR /app

ADD . .

RUN make all

## Deploy
FROM gcr.io/distroless/base-debian11
WORKDIR /

COPY --from=build /app/bin/http-mirror /http-mirror

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/http-mirror"]
