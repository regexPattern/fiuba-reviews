FROM lukemathwalker/cargo-chef:latest-rust-1.68 as chef
WORKDIR /app

# Determine dependencies
FROM chef as planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

# Cache dependencies
FROM chef AS builder
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json
# Build and run data generator
COPY . .
RUN cargo run --release

# Copy generated data
FROM postgres:15.2-alpine
COPY --from=builder /app/init.sql /docker-entrypoint-initdb.d
