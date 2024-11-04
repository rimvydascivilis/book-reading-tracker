FROM node:20.10.0-alpine AS builder
WORKDIR /app

COPY package.json ./

RUN yarn install

COPY . .

RUN yarn build


# Build the final image
FROM nginx:stable-alpine AS app

WORKDIR /usr/share/nginx/html

COPY --from=builder app/build/ ./

COPY nginx/default.conf.template /etc/nginx/conf.d/default.conf.template

COPY nginx/entrypoint.sh /

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
