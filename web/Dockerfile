FROM node:18-alpine as dev
RUN apk add --no-cache bash
WORKDIR /app
COPY docker/entrypoint-dev.sh /usr/local/bin/entrypoint
CMD ["entrypoint"]

FROM node:18-alpine as build
WORKDIR /app
COPY package.json .
COPY package-lock.json . 
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.24-alpine as prod
COPY --from=build /app/build /usr/share/nginx/html
