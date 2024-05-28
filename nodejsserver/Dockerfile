# Use the official Node.js image based on Alpine for a smaller footprint
FROM node:latest-alpine

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
COPY package*.json ./
RUN npm install

# Bundle app source
COPY . .

# Expose the port
EXPOSE 8080

# Start the app
CMD ["npm", "run dev"]
