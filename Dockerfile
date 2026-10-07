FROM golang:1.27.1-bookworm AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd

RUN CGO_ENABLED=0 go build -o /out/import-channels ./cmd/import-channels

FROM scratch

COPY --from=builder /out/import-channels /import-channels

USER 65532:65532
ENTRYPOINT ["/import-channels"]
