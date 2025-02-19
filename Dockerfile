FROM alpine:3.19
WORKDIR /app
COPY ./temporal-benchmark /app/temporal-benchmark
ENTRYPOINT ["/app/temporal-benchmark"] 