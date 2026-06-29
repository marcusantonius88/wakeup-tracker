FROM node:20-alpine AS deps
WORKDIR /app
COPY frontend/web-app/package.json frontend/web-app/package-lock.json* ./
RUN npm ci

FROM node:20-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY frontend/web-app/app ./app
COPY frontend/web-app/components ./components
COPY frontend/web-app/lib ./lib
COPY frontend/web-app/public ./public
COPY frontend/web-app/next.config.mjs ./next.config.mjs
COPY frontend/web-app/tsconfig.json ./tsconfig.json
COPY frontend/web-app/tailwind.config.ts ./tailwind.config.ts
COPY frontend/web-app/postcss.config.js ./postcss.config.js
COPY frontend/web-app/next-env.d.ts ./next-env.d.ts
COPY frontend/web-app/package.json ./package.json
RUN npm run build

FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app/.next ./.next
COPY --from=build /app/public ./public
COPY --from=build /app/package.json ./package.json
COPY --from=deps /app/node_modules ./node_modules
EXPOSE 3000
CMD ["npm", "start"]
