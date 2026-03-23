# ML_Elec Edge Agent - Dockerfile multi-arch
# Basé sur node:20-alpine pour une image légère
# Supporte linux/amd64 et linux/arm64 (Raspberry Pi)

# ==============================================================================
# BUILD STAGE
# ==============================================================================
FROM node:20-alpine AS builder

# Installer les dépendances système pour DuckDB et node-opcua
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    linux-headers \
    openssl-dev

# Working directory
WORKDIR /app

# Installer pnpm
RUN corepack enable && corepack prepare pnpm@9 --activate

# Copier les fichiers de configuration
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml* turbo.json ./
COPY packages/shared/package.json ./packages/shared/
COPY packages/edge-agent/package.json ./packages/edge-agent/
COPY packages/ui/package.json ./packages/ui/
COPY apps/web/package.json ./apps/web/

# Installer les dépendances
RUN pnpm install --frozen-lockfile --prefer-offline

# Copier le code source
COPY . .

# Build des packages
RUN pnpm run build

# ==============================================================================
# PRODUCTION STAGE
# ==============================================================================
FROM node:20-alpine AS production

# Labels
LABEL org.opencontainers.image.title="ML_Elec Edge Agent"
LABEL org.opencontainers.image.description="Plateforme de maintenance prédictive pour l'électrotechnique"
LABEL org.opencontainers.image.source="https://github.com/ml-elec/ml-elec"
LABEL org.opencontainers.image.licenses="MIT"

# Installer les dépendances système minimales pour DuckDB
RUN apk add --no-cache \
    libgcc \
    libstdc++ \
    openssl

# Créer un utilisateur non-root pour la sécurité
RUN addgroup -g 1001 -S ml-elec && \
    adduser -S ml-elec -u 1001 -G ml-elec

# Working directory
WORKDIR /app

# Copier les artifacts buildés depuis le stage builder
COPY --from=builder --chown=ml-elec:ml-elec /app/packages/edge-agent/dist ./dist
COPY --from=builder --chown=ml-elec:ml-elec /app/packages/edge-agent/node_modules ./node_modules
COPY --from=builder --chown=ml-elec:ml-elec /app/packages/shared/dist ./packages/shared/dist
COPY --from=builder --chown=ml-elec:ml-elec /app/packages/shared/node_modules ./packages/shared/node_modules

# Copier le package.json pour les dépendances
COPY --from=builder --chown=ml-elec:ml-elec /app/packages/edge-agent/package.json ./package.json
COPY --from=builder --chown=ml-elec:ml-elec /app/packages/shared/package.json ./packages/shared/package.json

# Créer le dossier de données
RUN mkdir -p /data/duckdb && chown -R ml-elec:ml-elec /data

# Changer vers l'utilisateur non-root
USER ml-elec

# Exposer le port
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3001/health', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# Commande de démarrage
CMD ["node", "dist/index.js"]
