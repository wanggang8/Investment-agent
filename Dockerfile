FROM node:22-alpine AS web-build
WORKDIR /src/web
COPY web/package*.json ./
RUN npm ci
COPY web ./
RUN npm run build

FROM golang:1.25-alpine AS go-build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/investment-agent-server ./cmd/server

FROM nginx:1.27-alpine
RUN apk add --no-cache ca-certificates curl tzdata
WORKDIR /app
COPY --from=go-build /out/investment-agent-server /usr/local/bin/investment-agent-server
COPY --from=web-build /src/web/dist /usr/share/nginx/html
COPY configs/config.docker.yaml /app/configs/config.docker.yaml
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY docker/entrypoint.sh /usr/local/bin/investment-agent-entrypoint
COPY docker/healthcheck.sh /usr/local/bin/investment-agent-healthcheck
RUN chmod +x /usr/local/bin/investment-agent-entrypoint /usr/local/bin/investment-agent-healthcheck
ENV INVESTMENT_AGENT_CONFIG=/app/configs/config.docker.yaml
EXPOSE 4173 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 CMD /usr/local/bin/investment-agent-healthcheck
ENTRYPOINT ["/usr/local/bin/investment-agent-entrypoint"]
