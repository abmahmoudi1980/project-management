# Data Model: Project Member Role Assignment

**Feature**: 001-project-member-roles  
**Date**: 2026-02-26  
**Status**: ✅ Implemented

## Overview

This feature adds explicit project membership and predefined project roles while reusing existing `users` and `projects` tables.

## Implementation Notes

- **Migration**: `backend/migration/008_add_project_members.sql`
- **Run command**: `cd backend && go run ./cmd/migrate`
- **Role ordering**: Custom CASE sort (owner=1, manager=2, contributor=3, viewer=4)

## New Entities

### 1) ProjectRole

**Table**: `project_roles`

| Field | Type | Constraints | Description |
|---|---|---|---|
| id | UUID | PK, NOT NULL | Role identifier |
| name | VARCHAR(50) | NOT NULL, UNIQUE | Internal role key (e.g., `manager`) |
| display_name | VARCHAR(100) | NOT NULL | UI display label (Persian: "مدیر") |
| is_active | BOOLEAN | NOT NULL DEFAULT true | Whether role is assignable |
| created_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Update timestamp |

**Initial seed set**:
- `owner` (مالک)
- `manager` (مدیر)
- `contributor` (مشارکت‌کننده)
- `viewer` (بیننده)

**Validation Rules**:
- `name`: lowercase slug, unique, immutable after seed for MVP.
- Only rows with `is_active = true` can be assigned (enforced in service layer).

### 2) ProjectMember

**Table**: `project_members`

| Field | Type | Constraints | Description |
|---|---|---|---|
| project_id | UUID | PK(part), FK → `projects(id)` ON DELETE CASCADE | Project reference |
| user_id | UUID | PK(part), FK → `users(id)` ON DELETE CASCADE | User reference |
| project_role_id | UUID | NOT NULL, FK → `project_roles(id)` | Assigned predefined role |
| joined_at | TIMESTAMPTZ | NOT NULL DEFAULT NOW() | Membership creation time |
| added_by | UUID | NULL, FK → `users(id)` | Admin who added member |

**Primary Key**: (`project_id`, `user_id`)

**Indexes**:
- `idx_project_members_project_id` on (`project_id`)
- `idx_project_members_user_id` on (`user_id`)
- `idx_project_members_project_role_id` on (`project_role_id`)

**Validation Rules**:
- Membership must be unique per project/user pair.
- `project_role_id` must reference active predefined role at add time (enforced in service).
- `user_id` must reference active/eligible user (enforced in service).

## Existing Entity Dependencies

### User (existing)

Used fields:
- `id`
- `username`
- `email`
- `is_active`
- `role` (system role: `admin`/`user`)

Constraint usage:
- `is_active = true` required for eligible-user listing and add operation.

### Project (existing)

Used fields:
- `id`
- `title`

Constraint usage:
- Membership operations require existing project ID.

## Relationships

- One `Project` has many `ProjectMember` rows.
- One `User` can belong to many projects through `ProjectMember`.
- One `ProjectRole` can be assigned to many `ProjectMember` rows.

## State/Transition Rules

### Membership lifecycle

1. **Not Member** → **Member Added**
   - Trigger: admin creates membership with valid role.
   - Validation: project exists, user active, role active, not duplicate.
2. **Member Added** → **Role Changed** (optional endpoint)
   - Trigger: admin updates member role to another active predefined role.
3. **Member Added/Role Changed** → **Removed** (optional endpoint)
   - Trigger: admin removes member from project.

For MVP alignment with spec, required transition is **Not Member → Member Added with role**.

## Derived Query Model

### Eligible users for add-member flow

Query filters:
- `users.is_active = true`
- User ID not in `SELECT user_id FROM project_members WHERE project_id = $projectId`
- Ordered by username ASC

### Project member list

Return joined projection (from `ProjectMemberResponse`):
- user fields (`id`, `username`, `email`)
- role fields (`project_role_id`, `project_roles.name`, `project_roles.display_name`)
- membership metadata (`joined_at`, `added_by`)
- Ordered by `joined_at DESC`

### Active roles list

Query filters:
- `is_active = true`
- Custom ordering by hierarchy (owner → manager → contributor → viewer)

## API Response Formats

All endpoints return arrays directly (not wrapped in object):
- `GET /api/projects/{id}/members` → `ProjectMemberResponse[]`
- `GET /api/projects/{id}/members/eligible-users` → `UserSummary[]`
- `GET /api/project-roles` → `ProjectRole[]`
- `POST /api/projects/{id}/members` → `{ message: string }`

See `contracts/openapi.yaml` for full API specification.

## Files

- **Migration**: `backend/migration/008_add_project_members.sql`
- **Models**: `backend/models/project_member.go`
- **Repository**: `backend/repositories/project_member_repository.go`
- **Service**: `backend/services/project_member_service.go`
- **Handler**: `backend/handlers/project_member_handler.go`
- **Routes**: `backend/routes/routes.go` (lines 61-66)
- **Frontend API**: `frontend/src/lib/api/projectMembers.js`
- **Frontend UI**: `frontend/src/components/ProjectMembersPanel.svelte`
