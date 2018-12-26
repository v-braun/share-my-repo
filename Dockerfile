
# build client
FROM node:8 AS build-env-client
WORKDIR /app
COPY conf.js ./
COPY gulpfile.js ./
COPY package.json ./
COPY ./client ./client/
RUN npm install
RUN npm run dist



FROM golang:1.8

WORKDIR /go/src/github.com/v-braun/share-my-repo
COPY . .
COPY --from=build-env-client /app/bin ./bin/

RUN go get -d -v ./...

EXPOSE 3001

CMD ["go", "run", "main.go"]

