# ai-test-assets-service — Claude Code Context

> Generated from ai-test-platform-standards/templates/service-CLAUDE.md
> Replace all [bracketed] placeholders before committing.

---

## Platform standards (authoritative)

This service follows the platform constitution. Key references — load these on demand:

- Ubiquitous language: https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/domain/ubiquitous-language.md
- Context map: https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/domain/bounded-contexts.md
- ADR template: https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/adr/0000-template.md
- Platform ADRs: https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/adr/
- BDD conventions: https://github.com/Snip123/ai-test-platform-standards/blob/main/docs/features/cross-context/

> If you cannot fetch these URLs, ask the developer to paste the relevant sections.

---

## This service

**Bounded context:** Assets
**Owned by:** Platform Guild
**Responsibility:** [One sentence — what this service is solely responsible for]

### Local docs (load on demand)
- Local ADRs: @docs/adr/
- Local feature specs: @docs/features/
- Local domain supplement: @docs/domain/ubiquitous-language-supplement.md (if exists)

---

## Workflow rules (service-level)

### Before implementing any feature
1. Check @docs/features/ for an existing `.feature` file
2. If none exists — **STOP. Do not write any code.** Write a Gherkin draft, present it to the developer, and wait for explicit confirmation before touching any `.go` file. This applies regardless of how the request was phrased.
3. Verify all terms used are in the platform ubiquitous language

### Before any architectural decision
1. Check @docs/adr/ for existing local decisions
2. Check ai-test-platform-standards ADRs for platform-level decisions
3. If making a new significant decision — draft an ADR at `docs/adr/NNNN-title.md`
4. Service-local ADRs: NNNN starts at 0001 for this repo
5. Platform-wide decisions → flag for a PR to ai-test-platform-standards

### Domain language
- ONLY use terms from the platform ubiquitous language
- Service-local terms must be in @docs/domain/ubiquitous-language-supplement.md
- Never invent terms; never use synonyms

---

## Service-specific rules

> Add any rules unique to this bounded context below.
> Examples: data formats, external API constraints, security rules, performance budgets.

### [Example: monetary values]
- All monetary amounts are stored and transmitted as integers in minor currency units (e.g. cents)
- Never use floats for money

### [Example: external constraints]
- [Third-party API name]: rate limit is X req/s — always use the retry wrapper in `src/lib/[name]`

---

## Build & test

```bash
# Install dependencies
[insert command]

# Run tests
[insert command]

# Run BDD specs
[insert command]

# Lint
[insert command]

# Start locally
[insert command]
```

---

## ADR index (this service)

| ID | Title | Status |
|----|-------|--------|
| — | No local ADRs yet | — |
