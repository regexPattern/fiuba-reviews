ARG NODE_VERSION=20.10.0
FROM node:${NODE_VERSION}-slim AS base

LABEL fly_launch_runtime="SvelteKit"

WORKDIR /app

ENV NODE_ENV="production"

ARG PNPM_VERSION=9.4.0
RUN npm install -g pnpm@$PNPM_VERSION


FROM base AS build

RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y build-essential node-gyp pkg-config python-is-python3

COPY --link .npmrc package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --prod=false

COPY --link . .

RUN --mount=type=secret,id=DATABASE_URL                                    \
    --mount=type=secret,id=PUBLIC_TURNSTILE_SITE_KEY                       \
    --mount=type=secret,id=TURNSTILE_SECRET_KEY                            \
    DATABASE_URL=`cat /run/secrets/DATABASE_URL`                           \
    PUBLIC_TURNSTILE_SITE_KEY=`cat /run/secrets/PUBLIC_TURNSTILE_SITE_KEY` \
    TURNSTILE_SECRET_KEY=`cat /run/secrets/TURNSTILE_SECRET_KEY`           \
    pnpm run build

# RUN pnpm prune --prod
RUN ls


FROM base

COPY --from=build /app/build /app/build
COPY --from=build /app/node_modules /app/node_modules
COPY --from=build /app/package.json /app

EXPOSE 3000
CMD [ "node", "./build/index.js" ]
