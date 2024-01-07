FROM node:20-slim AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /fiuba-reviews
COPY . .

RUN pnpm install

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

CMD /wait && pnpm build && pnpm preview --host
