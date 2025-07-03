FROM alpine:latest

WORKDIR /app

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Install Vite globally
RUN npm install -g vite

COPY package.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]