FROM golang:1.27 AS build
WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/server /app/server

EXPOSE 8080
CMD ["/app/server"]