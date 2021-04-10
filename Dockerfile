FROM golang AS builder
WORKDIR /build
COPY ["*.go", "go.mod", "./"]
RUN go mod download && CGO_ENABLED=0 go build -o server .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /build/server /
ENTRYPOINT ["/server"]