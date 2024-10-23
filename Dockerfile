FROM golang:bookworm AS builder

WORKDIR /go/src/
COPY ./ /go/src/
RUN mkdir build
RUN go build -o /go/src/build/salt-exporter

FROM debian:bullseye-slim AS runner
WORKDIR /app/salt-exporter/
COPY --from=builder /go/src/build/salt-exporter ./
CMD ["/app/salt-exporter/salt-exporter"]
