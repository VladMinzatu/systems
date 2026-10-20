# syntax=docker/dockerfile:1

FROM golang:1.27.1-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/tui-test .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/tui-test /tui-test
ENTRYPOINT ["/tui-test"]
