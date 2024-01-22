FROM --platform=linux/amd64 golang:1.21 AS build-stage

WORKDIR /app

COPY ../../../go.mod go.sum ./
RUN CGO_ENABLED=0 GOOS=linux go mod download

COPY ../../.. ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./pkg/cmd/bot/

FROM --platform=linux/amd64 ubuntu AS build-release-stage

RUN apt-get update && \
    apt-get install ca-certificates google-chrome -y && \
    apt-get clean

WORKDIR /

COPY --from=build-stage /bot /bot

CMD ["/bot"]