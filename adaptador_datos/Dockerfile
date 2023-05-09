FROM lukemathwalker/cargo-chef:latest-rust-1.68 as chef
WORKDIR /app

# determinar las dependencias
FROM chef as planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

# hacer cache a las dependencias
FROM chef AS builder
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

# copiar y correr el adaptador
COPY . .
RUN cargo run --release

# copiar los datos generados para cuando se inicialize la DB
FROM postgres:15.2
COPY --from=builder /app/init.sql /docker-entrypoint-initdb.d
