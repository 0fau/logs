FROM --platform=linux/amd64 golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /apiserver ./pkg/cmd/api/

FROM --platform=linux/amd64 ubuntu AS build-release-stage

RUN apt-get update && \
    apt-get install ca-certificates -y && \
    apt-get clean

WORKDIR /

COPY --from=build-stage /apiserver /apiserver

CMD mkdir -p migrations

COPY migrations/* migrations/

EXPOSE 3001

CMD ["/apiserver"]