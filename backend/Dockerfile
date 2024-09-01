FROM rust:1.80 as builder
WORKDIR /app
COPY . .
RUN cargo install sqlx-cli --no-default-features --features postgres

FROM rust:1.80.1-slim-bookworm
RUN apt-get update && apt-get install -y libpq-dev postgresql-client && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo /usr/local/cargo
COPY --from=builder /app /app
COPY start.sh /start.sh
RUN chmod +x /start.sh

ENV PATH="/usr/local/cargo/bin:${PATH}"
WORKDIR /app
CMD ["/start.sh"]
