FROM golang:1.18-alpine as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN cd cmd; CGO_ENABLED=0 go build -ldflags "-s -w" -installsuffix cgo -o /cmd/rating

FROM gcr.io/distroless/static
COPY --from=build /cmd/rating /
CMD ["/rating"]