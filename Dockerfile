FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/careeros ./cmd/careeros
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/career-migrate ./cmd/career-migrate

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/careeros /app/careeros
COPY --from=build /out/career-migrate /app/career-migrate
COPY database/migrations /app/database/migrations
COPY web /app/web
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/careeros"]
