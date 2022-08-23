FROM golang:1.18-alpine as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -installsuffix cgo -o brbarmex-rating-app .

FROM gcr.io/distroless/static
COPY --from=build /app/.  /

CMD ["/brbarmex-rating-app"]