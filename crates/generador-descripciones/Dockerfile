FROM lukemathwalker/cargo-chef:latest-rust-1.70 AS chef
WORKDIR /fiuba-reviews

FROM chef as planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder
COPY --from=planner /fiuba-reviews/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

COPY . .
RUN cargo build --release --bin generador_descripciones

FROM debian:buster-slim AS runtime
WORKDIR /fiuba-reviews

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait
COPY --from=builder /fiuba-reviews/target/release/generador_descripciones /usr/local/bin
CMD /wait && generador_descripciones
