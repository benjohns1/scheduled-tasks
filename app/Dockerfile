FROM node:12-alpine AS builder
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY *.js ./
COPY src src/
RUN npm run build

FROM node:12-alpine
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install --production
COPY --from=builder /usr/src/app/__sapper__/build __sapper__/build/
COPY static static/
EXPOSE 3000
CMD ["node", "__sapper__/build"]