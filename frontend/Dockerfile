
# Step 1: Build the React app
FROM node:18-alpine AS build

# Set the working directory
WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the application code
COPY . ./

# Build the React app
RUN npm run build:prod

# Step 2: Serve the React app
FROM nginx:alpine

# Copy the build directory to nginx's html directory
COPY --from=build /app/build /data/www/

COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start nginx server
CMD ["nginx", "-g", "daemon off;"]
