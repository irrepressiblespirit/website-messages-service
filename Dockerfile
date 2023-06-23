FROM golang:1.17 as builder
RUN mkdir /build
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o target/application cmd/main.go

FROM scratch
COPY --from=builder /build/target /app/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY api/proto/ /app/api/proto/
ENV TZ=UTC
WORKDIR /app
CMD ["./application"]
