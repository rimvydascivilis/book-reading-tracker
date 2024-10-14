FROM node:20.10.0-alpine AS builder
WORKDIR /app

COPY package.json ./

RUN yarn install

COPY . .

RUN yarn build


# Build the final image
FROM nginx:stable-alpine AS app

ARG PORT=80
ARG BACKEND_URL=http://localhost:8000

# Set env vars for the entrypoint
ENV PORT=$PORT
ENV BACKEND_URL=$BACKEND_URL

WORKDIR /usr/share/nginx/html

COPY --from=builder app/build/ ./

COPY nginx/default.conf.template /etc/nginx/conf.d/default.conf.template

COPY nginx/entrypoint.sh /

RUN chmod +x /entrypoint.sh

EXPOSE $PORT

ENTRYPOINT ["/entrypoint.sh"]
