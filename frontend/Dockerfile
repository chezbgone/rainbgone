FROM node:22-alpine AS dev

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

EXPOSE 5173

CMD ["sh", "-c", "npm ci && exec npm run dev -- --host 0.0.0.0"]

FROM dev AS build
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
