# Quickstart: Project Member Role Assignment

**Feature**: 001-project-member-roles  
**Branch**: `001-project-member-roles`  
**Date**: 2026-02-26

## Goal

Allow admins to add any eligible system user to a project and assign a predefined project role.

## 1) Apply database migration

1. Add migration file: `backend/migration/008_add_project_members.sql`
2. Migration contents should:
   - create `project_roles`
   - seed default active roles (`owner`, `manager`, `contributor`, `viewer`)
   - create `project_members` with unique (`project_id`, `user_id`)
   - add indexes for project/user/role lookups
3. Run migration:

```bash
cd backend
go run ./cmd/migrate
```

## 2) Implement backend layers

1. **Models**
   - Add `backend/models/project_member.go`
   - Add DTOs for add-member request and list responses

2. **Repository**
   - Add `backend/repositories/project_member_repository.go`
   - Implement queries for:
     - list members by project
     - list eligible users by project
     - list active project roles
     - add project member with role

3. **Service**
   - Add `backend/services/project_member_service.go`
   - Validate:
     - project exists
     - caller is admin for add operation
     - user is eligible and active
     - role is predefined and active
     - duplicate memberships are rejected

4. **Handler + Routes**
   - Add `backend/handlers/project_member_handler.go`
   - Register routes in `backend/routes/routes.go`:
     - `GET /api/projects/:projectId/members`
     - `POST /api/projects/:projectId/members`
     - `GET /api/projects/:projectId/members/eligible-users`
     - `GET /api/project-roles`

5. **Wire dependencies**
   - Update `backend/main.go` to initialize and inject new handler/service/repository.

## 3) Implement frontend integration

1. Add API client methods in `frontend/src/lib/api/projectMembers.js`.
2. Update project member management UI component to:
   - fetch eligible users
   - fetch predefined roles
   - submit add-member request with `user_id` and `project_role_id`
3. Render assigned member role in project member list.

## 4) Verify behavior manually

1. Login as admin.
2. Open a project’s member management view.
3. Add an active non-member user and choose role `manager`.
4. Confirm member appears with correct role.
5. Retry adding same user; confirm duplicate is blocked.
6. Try invalid role ID; confirm validation error.

## 5) API smoke checks (example)

```bash
# List predefined roles
curl -i -X GET http://localhost:3000/api/project-roles

# List eligible users for project
curl -i -X GET http://localhost:3000/api/projects/<project-id>/members/eligible-users

# Add member with role
curl -i -X POST http://localhost:3000/api/projects/<project-id>/members \
  -H "Content-Type: application/json" \
  -d '{"user_id":"<user-id>","project_role_id":"<role-id>"}'
```
