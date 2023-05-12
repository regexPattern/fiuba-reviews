FROM lukemathwalker/cargo-chef:latest-rust-1.68 AS chef
WORKDIR /fiuba-reviews

FROM chef as planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder
COPY --from=planner /fiuba-reviews/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

COPY . .
RUN cargo run --release

FROM postgres:15.2
COPY --from=builder /fiuba-reviews/init.sql /docker-entrypoint-initdb.d
