FROM --platform=linux/amd64 golang:1.21 AS build-stage

WORKDIR /app

COPY ../../../go.mod go.sum ./
RUN CGO_ENABLED=0 GOOS=linux go mod download

COPY ../../.. ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /screenshot ./pkg/cmd/screenshot/

FROM --platform=linux/amd64 chromedp/headless-shell:latest AS build-release-stage

RUN apt-get update && \
    apt-get install dumb-init ca-certificates -y && \
    apt-get clean

ENTRYPOINT ["dumb-init", "--"]

WORKDIR /

COPY --from=build-stage /screenshot /screenshot

EXPOSE 3010

CMD ["/screenshot"]