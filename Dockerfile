FROM --platform=linux/amd64 golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /apiserver ./pkg/cmd/api/

# Deploy the application binary into a lean image
FROM --platform=linux/amd64 gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /apiserver /apiserver

CMD mkdir -p /pkg/database/migrations

COPY ./pkg/database/migrations/* /pkg/database/migrations/

EXPOSE 3001

CMD ["/apiserver"]