# ai-test-assets-service

**Bounded context:** Assets
**Team:** Platform Guild

## Responsibility

> [One sentence — what this service is solely responsible for]

## Platform standards

This service follows the platform constitution:
https://github.com/Snip123/ai-test-platform-standards

See [CLAUDE.md](CLAUDE.md) for Claude Code context.

## Getting started

```bash
# Start platform dependencies
cd ../ai-test-platform-standards && docker compose up -d

# Copy and fill in local env
cp .env.example .env

# Run the service
make run

# Run tests
make test

# Run BDD scenarios
make test-bdd
```

## Project structure

See [ADR-0018](https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/adr/0018-standard-go-service-structure.md) for the full structure rationale.

```
cmd/server/    HTTP server entrypoint
cmd/migrate/   Migration job entrypoint
internal/      All application code (not importable externally)
docs/api/      OpenAPI spec (source of truth)
docs/features/ BDD Gherkin specs
```

## Docs

### This service
- [API spec](docs/api/openapi.yaml)
- [Architecture decisions](docs/adr/README.md)
- [Feature specs](docs/features/README.md)
- [Claude Code context](CLAUDE.md)

### Platform (authoritative)
- [Platform standards](https://github.com/Snip123/ai-test-platform-standards)
- [Platform ADRs](https://github.com/Snip123/ai-test-platform-standards/tree/main/docs/adr)
- [Ubiquitous language](https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/domain/ubiquitous-language.md)
- [Bounded contexts](https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/domain/bounded-contexts.md)
