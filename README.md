# backend

### Prerequisites

Make sure you have the following tools installed to get *QUALITY OF LIFE*

- go v1.22.11
- wire (dependency injection tool)
```bash
go install github.com/google/wire/cmd/wire@latest
```
- Makefile (command alias tool)
- air (auto-reload, similar to nodemon)
```bash
go install github.com/air-verse/air@latest
```
> **Note:** install with go v1.23 or higher

### Project structures

- Each sub-project in `/cmd` has its own `/cmd/<name>/internal` directory which contains its private components such as handlers or router.
- `/internal` is a shared directory and is shared among sub-project such as config, models, or services.
- `.env` is also shared among sub-projects for the sake of simplicity. That means every sub-projects have the same set of envs.

### Collaboration

1. Split new branch
2. Copy `.env.example` to `.env` and fill out all the empty fields
3. Run auto-reloader
```bash
make dev-agent # agent
make dev-backend # backend
```
4. If `wire.go` in any sub-projects change, generate `wire_gen.go` by
```bash
make gen
```
5. Open pull request with title in the same fashion as commit message ex. `feat: blah blah blah`
