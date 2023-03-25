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
# Build backend
COPY . .
RUN cargo install --path .

# Run backend
FROM debian:bullseye-slim
RUN apt-get update
WORKDIR /app
COPY --from=builder /usr/local/cargo/bin/backend /usr/local/bin/backend
CMD ["backend"]
