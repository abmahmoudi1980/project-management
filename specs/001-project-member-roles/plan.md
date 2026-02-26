# Implementation Plan: Project Member Role Assignment

**Branch**: `001-project-member-roles` | **Date**: 2026-02-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-project-member-roles/spec.md`

## Summary

Enable admin users to add project members by selecting any eligible system user and assigning a predefined project role at add time. Implementation extends existing Go/Fiber layered backend with project-member APIs and persistence, plus frontend member-management UI updates in Svelte 5 to select users and roles from controlled data sources.

## Technical Context

**Language/Version**: Go 1.25 (backend), JavaScript ES modules with Svelte 5 + Vite 6 (frontend)  
**Primary Dependencies**: Fiber v2, pgx v5, UUID, existing auth middleware (`RequireAuth`, `RequireRole`)  
**Storage**: PostgreSQL (existing schema + new membership/role tables)  
**Testing**: Go integration/manual API verification via existing backend test style; manual frontend verification  
**Target Platform**: Linux-hosted web application, browser-based frontend  
**Project Type**: Web application (`backend/` + `frontend/`)  
**Performance Goals**: Member add/list operations return within ~300ms p95 for normal project sizes (<200 members/project)  
**Constraints**: Must keep strict handler→service→repository layering; no direct DB access in handlers/services; admin-only mutation operations; no breaking auth flow  
**Scale/Scope**: Organization-level PM system (10–500 active users; up to thousands of projects/tasks), one membership management flow per project

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`/home/abolfazl/apps/project-management/.specify/memory/constitution.md` currently contains placeholder template sections (`[PRINCIPLE_*]`) without enforceable rules.

### Pre-Phase 0 Gate Result

- **Constitution status**: No ratified project constitution rules detected.
- **Gate decision**: ✅ PASS (no enforceable constraints to violate).
- **Applied governance fallback**: Follow repository conventions from `AGENTS.md` (Go layered architecture, pgx, Svelte 5 runes, minimal scope changes).

### Post-Phase 1 Re-Check

- **Design compliance**: ✅ PASS.
- **Reason**: Generated design preserves layered backend structure, uses existing auth/role middleware pattern, and keeps scope focused on membership + predefined roles.

## Project Structure

### Documentation (this feature)

```text
specs/001-project-member-roles/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
├── checklists/
│   └── requirements.md
└── tasks.md             # Created by /speckit.tasks
```

### Source Code (repository root)

```text
backend/
├── handlers/
│   └── project_member_handler.go      # NEW
├── services/
│   └── project_member_service.go      # NEW
├── repositories/
│   └── project_member_repository.go   # NEW
├── models/
│   └── project_member.go              # NEW
├── migration/
│   └── 008_add_project_members.sql    # NEW
├── routes/
│   └── routes.go                      # MODIFY
└── main.go                            # MODIFY wiring

frontend/
└── src/
    ├── lib/
    │   └── api/
    │       └── projectMembers.js      # NEW
    └── components/
        └── ProjectMembers.svelte      # MODIFY or NEW based on existing UI
```

**Structure Decision**: Use the existing web app split (`backend/`, `frontend/`) and add one cohesive project-membership vertical slice in each backend layer plus minimal frontend integration points.

## Phase 0: Research Output

- Completed in [research.md](research.md).
- All technical choices and integration patterns documented with decisions and alternatives.
- No unresolved `NEEDS CLARIFICATION` items remain.

## Phase 1: Design Output

- Data model: [data-model.md](data-model.md)
- API contract: [contracts/openapi.yaml](contracts/openapi.yaml)
- Implementation guide: [quickstart.md](quickstart.md)

## Complexity Tracking

No constitution violations or additional complexity exceptions required.
