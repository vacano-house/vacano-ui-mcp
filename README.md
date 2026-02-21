# vacano-ui-mcp

![Version](https://img.shields.io/badge/version-1.1.0-blue)

MCP server providing documentation for [vacano-ui](https://github.com/vacano-house/vacano-ui) React component library.

Clones the vacano-ui repository, parses markdown documentation, and exposes three MCP tools:
- **search_docs** — full-text search across component names, descriptions, and content
- **get_component_docs** — get full documentation for a specific component by name
- **list_components** — list all components, optionally filtered by category

## Categories

`form`, `data-display`, `feedback`, `layout`, `navigation`, `utility`, `guide`

Categories are automatically parsed from the VitePress sidebar config (`docs/.vitepress/config.ts`).

## Quick start

```bash
cp .env.example .env
make run
```

The server starts on port 3000 and exposes a single MCP endpoint at `/mcp`.

## Build

```bash
make build        # Build binary to ./bin/server
make docker-build # Build Docker image
```

## Docker

```bash
docker compose up -d
```

Runs on `127.0.0.1:3007` (mapped to container port 3000).

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `3000` | Server port |
| `GIT_REPO_URL` | `https://github.com/vacano-house/vacano-ui.git` | Git repository URL |
| `GIT_BRANCH` | `master` | Git branch |
| `GIT_SSH_KEY` | — | Optional SSH key for private repos |
| `DOCS_REFRESH_INTERVAL` | `5m` | Background refresh interval |
