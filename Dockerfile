FROM golang:1.24 AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /app/bin/wow-server ./server/cmd

RUN CGO_ENABLED=0 go build -o /app/bin/wow-client ./client/cmd

FROM alpine:3.14 AS wow-server
WORKDIR /app/
COPY --from=builder /app/bin/wow-server .
COPY quotes.txt .
CMD ["./wow-server"]

FROM alpine:3.14 AS wow-client
WORKDIR /app/
COPY --from=builder /app/bin/wow-client .
CMD ["./wow-client"]