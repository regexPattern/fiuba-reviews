FROM lukemathwalker/cargo-chef:latest-rust-1.76 AS chef
WORKDIR /adaptador-datos

FROM chef AS planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder
COPY --from=planner /adaptador-datos/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

COPY . .
RUN cargo run --release inicializar

FROM postgres:15.2
COPY --from=builder /adaptador-datos/init.sql /docker-entrypoint-initdb.d
