FROM node:20-bullseye as builder

WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:stable-bullseye

COPY --from=builder /usr/src/app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

CMD [ "nginx", "-g", "daemon off;" ]