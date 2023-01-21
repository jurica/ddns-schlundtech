FROM golang:alpine as builder
WORKDIR "/workspace"
COPY . .
RUN go build -o ddns-schlundtech

FROM alpine:latest
COPY --from=builder /workspace/ddns-schlundtech .
COPY --from=builder /workspace/ddns-schlundtech-template.toml ddns-schlundtech.toml
EXPOSE 8080
CMD ["./ddns-schlundtech"]