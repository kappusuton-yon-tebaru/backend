# backend

### Prerequisites

Make sure you have the following tools installed to get quality of life

- go v1.22.11
- wire (dependency injection tool)
```bash
go install github.com/google/wire/cmd/wire@latest
```
- Makefile (command alias tool)

### Project structures

- Each sub-project in `/cmd` has its own `/cmd/<name>/internal` directory which contains its private components such as handlers or router.
- `/internal` is a shared directory and is shared among sub-project such as config, models, or services.
- `.env` is also shared among sub-projects for the sake of simplicity. That means every sub-projects have the same set of envs.

### Get started

1. Copy `.env.example` to `.env` and fill out all the empty fields
