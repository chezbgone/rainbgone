FROM node:22-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:22-alpine

ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

WORKDIR /app

COPY --from=build /app/build ./build
COPY --from=build /app/package.json /app/package-lock.json ./

RUN npm ci --omit=dev --ignore-scripts && npm cache clean --force

USER node

EXPOSE 3000

CMD ["node", "build"]
