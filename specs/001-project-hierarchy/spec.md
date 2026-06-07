# Feature Specification: Project Hierarchy (Parent / Child Projects)

**Feature Branch**: `001-project-hierarchy`
**Created**: 2026-06-07
**Status**: Draft
**Input**: User description: "When defining a project, we must be able to specify the parent project. In the project list, child projects should be displayed in a tree structure under the parent project. In the screen of each project, a list of child projects should be displayed."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create a Project with a Parent (Priority: P1)

A team lead wants to create a new sub-project that belongs under an existing project (for example, a "Mobile App" project created under a parent "Product Development" project, or a "Q3 Marketing Campaign" under "Marketing 2026"). When creating a new project, the user must be able to pick an existing project as the parent, or leave the parent blank to create a top-level project.

**Why this priority**: This is the foundational capability. Without the ability to attach a parent on creation, no hierarchy can be built.

**Independent Test**: Can be fully tested by creating a new project with a parent selected, then verifying the new project appears nested under that parent in the list and on the parent's detail page. Delivers the value of organizing related work under a common umbrella from day one.

**Acceptance Scenarios**:

1. **Given** an authenticated user on the "New Project" form and at least one existing project, **When** the user fills in the project details and selects an existing project from a "Parent Project" picker, **Then** the new project is saved with the selected parent and confirmed to the user.
2. **Given** an authenticated user on the "New Project" form, **When** the user fills in the project details and leaves the "Parent Project" field blank, **Then** the new project is saved as a top-level (root) project.
3. **Given** an authenticated user on the "New Project" form, **When** the user attempts to select the new project as its own parent (or any of its descendants), **Then** the system rejects the action with a clear error message and does not save.

---

### User Story 2 - Browse Projects in a Tree Structure (Priority: P1)

A user opens the main "Projects" page to see all projects in the workspace. Instead of a flat list, the projects are displayed as a hierarchical tree. Top-level projects are shown at the root; child projects are visually nested under their parent. The user can expand and collapse each parent node to manage how much of the tree is visible at once.

**Why this priority**: This is the primary way users discover and navigate the project portfolio. A flat list becomes unusable once projects are organized hierarchically.

**Independent Test**: Can be fully tested by creating a workspace with several projects arranged in a multi-level hierarchy, opening the project list, and verifying the tree renders with correct parent-child indentation, expand/collapse controls, and shows the expected number of root projects.

**Acceptance Scenarios**:

1. **Given** a workspace containing a mix of top-level projects and projects with parents, **When** the user opens the project list, **Then** top-level projects are shown at the root level and each child project is visually nested (indented) under its parent.
2. **Given** a project list with at least one parent project that has children, **When** the user clicks the expand control on a parent, **Then** the children of that parent become visible; clicking the collapse control hides them again.
3. **Given** a project list with at least three levels of nesting, **When** the user expands a parent at every level, **Then** the deepest descendants are reachable in the tree without leaving the page.
4. **Given** a project list, **When** the user views it, **Then** each project's name is a link to that project's detail page.

---

### User Story 3 - View Child Projects on a Project's Detail Page (Priority: P2)

A user opens the detail page of a specific project to learn more about it. In addition to the project's own information, the page shows a clearly labeled list of the project's direct child projects. Each child in the list links to that child's own detail page, so the user can drill into the hierarchy from any point.

**Why this priority**: This is the natural complement to the tree view: from the tree the user goes to a parent's page, and from the parent's page the user goes down to a child. Without this, navigation from parent to child requires going back to the list.

**Independent Test**: Can be fully tested by opening any project that has children and verifying that the child projects are listed with working links. Projects with no children show an appropriate empty state.

**Acceptance Scenarios**:

1. **Given** a project that has at least one direct child project, **When** the user opens that project's detail page, **Then** a "Child Projects" section is displayed that lists the direct children by name, and each name is a link to that child's detail page.
2. **Given** a project that has no child projects, **When** the user opens that project's detail page, **Then** the "Child Projects" section is displayed with an empty-state message (for example, "No sub-projects yet") rather than being hidden.
3. **Given** a project that has both direct children and grandchildren, **When** the user opens that project's detail page, **Then** only the direct children are listed; grandchildren are not shown at this level.

---

### User Story 4 - Re-parent an Existing Project (Priority: P3)

A user realizes that a project is filed under the wrong parent (or under no parent at all), and wants to move it. The user can edit the project and change its parent to a different existing project, or to none (making it a top-level project).

**Why this priority**: Re-parenting is a real-world need (reorganizations, scope changes), but is secondary to being able to assign a parent at creation time.

**Independent Test**: Can be fully tested by editing an existing project, changing its parent, saving, and then verifying the project appears under the new parent in both the tree view and the old/new parent's detail pages.

**Acceptance Scenarios**:

1. **Given** an existing project, **When** the user edits it and selects a different existing project as its parent, **Then** the project is saved with the new parent and the old parent no longer lists it as a child.
2. **Given** an existing child project, **When** the user edits it and clears the "Parent Project" field, **Then** the project becomes a top-level (root) project.
3. **Given** an existing project, **When** the user attempts to re-parent it under one of its own descendants, **Then** the system rejects the change with a clear error message and leaves the parent unchanged.

---

### Edge Cases

- What happens when a user tries to delete a project that has child projects? The system MUST prevent deletion and inform the user that child projects exist; the user must first delete or re-parent the children.
- What happens when a project with a parent is created but the parent is later moved or changed? The child's parent reference is not implicitly modified; the child keeps the parent it was assigned at the time.
- How does the system handle a project that becomes orphaned (its parent is removed via an out-of-band operation)? The system MUST treat such a project as a top-level project on next read so that no data becomes inaccessible.
- What happens when the project tree becomes very deep (many levels)? There is no artificial depth limit; the tree view supports expand/collapse at every level, and the detail page lists only direct children regardless of total depth.
- What happens when a user opens a project list and there are no projects at all? The page shows a clear empty state and offers a way to create the first project.
- What happens when a user opens a project's detail page and the project itself has been deleted? The system returns a clear "project not found" state rather than crashing.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow a user to select an existing project as the parent of a new project at creation time, and MUST also allow leaving the parent unset to create a top-level (root) project.
- **FR-002**: The system MUST allow a user to change the parent of an existing project, including setting it to no parent (making it a top-level project).
- **FR-003**: The system MUST display the project list as a hierarchical tree in which top-level projects are shown at the root and child projects are nested under their parent in a visually distinct (indented) way.
- **FR-004**: The system MUST provide expand and collapse controls on each parent node in the tree so the user can control how much of the hierarchy is visible at once.
- **FR-005**: The system MUST display, on every project's detail page, a "Child Projects" section that lists the project's direct child projects by name, with each name linking to that child's detail page; if there are no direct children, the section MUST show a clear empty state.
- **FR-006**: The system MUST prevent a user from creating a circular reference in the project hierarchy (a project cannot be its own ancestor or descendant) and MUST show a clear error message when such an action is attempted.
- **FR-007**: The system MUST prevent deletion of a project that has one or more child projects and MUST inform the user that the children must be removed or re-parented first.
- **FR-008**: The system MUST treat any project whose parent no longer exists (orphaned project) as a top-level project so the data remains accessible.
- **FR-009**: The system MUST show each project's name as a working link to its own detail page wherever the project is listed (tree view, child list, and any other project listing).
- **FR-010**: The system MUST persist the parent project assignment across sessions so that the hierarchy is stable and survives reloads and restarts.

### Key Entities *(include if feature involves data)*

- **Project**: Represents a unit of work in the system. Attributes include a unique identifier, a name, a description, a reference to an optional parent Project, and standard lifecycle information (creation date, status, etc.). A Project may have zero or many child Projects and at most one parent Project.
- **Project Hierarchy**: The implicit, ordered, parent-to-child relationship between Projects. The hierarchy is a forest (a set of trees) because multiple top-level Projects are allowed, and every Project has at most one parent.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can create a new sub-project by selecting a parent and filling in the required fields in under 60 seconds.
- **SC-002**: The project list view renders a tree of at least three levels of nesting and at least 50 projects in under 2 seconds on a standard workstation.
- **SC-003**: 100% of projects that have a parent appear nested under that parent in the project list view, and 100% of projects that have children appear in the "Child Projects" section of their parent's detail page.
- **SC-004**: The system rejects 100% of attempts to create or introduce a circular reference in the hierarchy and surfaces a clear, user-readable error message in each case.
- **SC-005**: The system blocks 100% of attempts to delete a project that has child projects and explains to the user what must be done first.
- **SC-006**: A user can expand any parent node in the tree to reveal its children and collapse it again to hide them, with the entire action completing in under 200 milliseconds per node.
- **SC-007**: A user can navigate from the project list to any project's detail page, and from any project's detail page to any of its children's detail pages, in no more than two clicks per level of the hierarchy.

## Assumptions

- A project has at most one parent (the hierarchy is a forest of trees, not a directed acyclic graph).
- Any user who has permission to create or edit a project is also allowed to set or change its parent to any other project; finer-grained "who can be a parent of whom" rules are out of scope for this feature.
- There is no enforced maximum depth of nesting; the system supports arbitrarily deep hierarchies.
- "Top-level project" means a project whose parent is not set; this is the default for new projects unless the user explicitly picks a parent.
- Re-parenting a project does not automatically change anything about its children; children keep their existing parents.
- Orphan handling (FR-008) is a defensive measure for unexpected data states; normal deletion is gated by FR-007.
