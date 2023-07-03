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

FROM debian:bookworm-slim
WORKDIR /fiuba-reviews

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait
COPY --from=builder /fiuba-reviews/target/release/generador_descripciones /usr/local/bin

RUN apt-get update && apt-get install wget -y \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*

RUN wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.1f-1ubuntu2.19_amd64.deb
RUN dpkg -i libssl1.1_1.1.1f-1ubuntu2.19_amd64.deb

CMD /wait && /usr/local/bin/generador_descripciones
