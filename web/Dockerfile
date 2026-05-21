FROM oven/bun:alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

ARG VITE_POSTHOG_API_KEY
ARG VITE_POSTHOG_API_HOST

ENV VITE_POSTHOG_API_KEY=$VITE_POSTHOG_API_KEY
ENV VITE_POSTHOG_API_HOST=$VITE_POSTHOG_API_HOST

RUN bun run build

FROM alpine:3

COPY --from=builder /app/dist /dist

ENTRYPOINT ["cp", "-r", "/dist/.", "/app/"]
