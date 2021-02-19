FROM golang:alpine AS build

WORKDIR /build

COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

WORKDIR /app

ENV DATABASE_URL foo

COPY --from=build /build .

CMD ["/app/main"]