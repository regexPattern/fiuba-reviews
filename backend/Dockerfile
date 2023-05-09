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

# copiar e instalar la aplicacion
COPY . .
RUN cargo install --path .

# actualizar el backend
FROM debian:bullseye-slim
RUN apt-get update

# copiar script utilizado para esperar a que se conecte la DB
COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

# copiar y correr el servidor
WORKDIR /app
COPY --from=builder /usr/local/cargo/bin/backend /usr/local/bin/backend
CMD /wait && backend
