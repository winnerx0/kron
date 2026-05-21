FROM oven/bun:alpine

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

ARG VITE_API_URL
ARG VITE_POSTHOG_API_KEY
ARG VITE_POSTHOG_API_HOST

ENV VITE_API_URL=$VITE_API_URL
ENV VITE_POSTHOG_API_KEY=$VITE_POSTHOG_API_KEY
ENV VITE_POSTHOG_API_HOST=$VITE_POSTHOG_API_HOST

RUN bun run build

EXPOSE 3000

CMD ["bun", "run", "preview", "--", "--host", "--port", "3000"]
