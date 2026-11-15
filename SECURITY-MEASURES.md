# Security measures: Sicheong

These are the security measures that exist in the code and configuration today. Each one addresses a threat listed in [THREAT-MODEL.md](THREAT-MODEL.md). Planned work is listed separately and is not counted as a measure.

| # | Measure | Threat | Where it lives | Status | Notes |
|---|---|---|---|---|---|
| M1 | Database values are passed as separate parameters, never as SQL text | 6 | `db/` | Implemented | Does not fix the author-filter bug (threat 9). |
| M2 | Sort and direction options come from a fixed list | 6 | `Params.Query()` in `db/db.go` | Implemented | Allow-list only; not general input cleaning. |
| M3 | Only the post's author can edit or delete it | 5 | `http/owner.go` | Implemented | Tests still to be added. |
| M4 | The author of an edited post comes from the existing post, not the form | 5 | Edit handlers in `http/` | Implemented | Check every write path, including the API. |
| M5 | Page templates escape user text | 7 | Templates in `http/` | Implemented | User text must not be marked as trusted HTML. |
| M6 | Login cookie is encrypted and marked `HttpOnly` and `SameSite` | 4, 8 | `http/http.go`, `http/human.go` | Partly implemented | Cookie holds the password, has no expiry, and its key changes on restart. |
| M7 | Public site works over HTTPS only | 12 | `fly.toml` | Implemented on Fly | Render and GKE to be checked. |
| M8 | CI checks formatting, tests, linting and Gosec on every change | 10, 11 | `.github/workflows/ci.yml` | Implemented | Does not block deploys; no `govulncheck` yet. |
| M9 | GitHub watches for vulnerable libraries and secrets, and runs code scanning | 11 | Repository settings; CodeQL workflow | Implemented | No `dependabot.yml` yet; alerts need regular review. |
| M10 | Database file is kept out of Git | Data copy | `.gitignore` | Implemented | Verified in `.gitignore`. |

## Planned, not yet in place

- Hashed passwords (threat 2).
- Removing default passwords (threat 1).
- Login rate limiting (threat 3).
- Session expiry and logout (threat 4).
- Deploys that wait for tests (threat 10).
