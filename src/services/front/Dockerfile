# Stage 1: Build the Angular application with Bun
FROM node:lts-slim AS base
WORKDIR /app

COPY . .
RUN npm i -g bun

FROM base AS all-deps
RUN bun install

FROM base AS prod-deps

RUN bun install --production

FROM all-deps AS build
# Set working directory
WORKDIR /app

# Copy package.json, bun.lockb, and other necessary files
COPY package.json bun.lockb ./

# Install dependencies
RUN bun install

# Copy the rest of the application files
COPY . .

# Build the Angular application
RUN bun run ng build --base-href ./ --deploy-url ./

# Stage 2: Serve the application with NGINX
FROM nginx:alpine-slim

# Remove the default NGINX website
RUN rm -rf /usr/share/nginx/html/*

# Copy the built Angular app from the build stage
COPY --from=build /app/dist/front/browser /usr/share/nginx/html

# Copy custom NGINX configuration file (optional)
# COPY nginx.conf /etc/nginx/nginx.conf
# RUN chmod -R 755 /usr/share/nginx/html
# Expose the port on which the app will run
EXPOSE 80

# Start NGINX
CMD ["nginx", "-g", "daemon off;"]
