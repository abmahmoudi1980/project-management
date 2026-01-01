# Quickstart Guide: Dashboard Implementation

**Feature**: 005-dashboard  
**Date**: 2026-01-01  
**Branch**: `005-dashboard`  

---

## Overview

This guide provides step-by-step instructions for implementing the project manager dashboard feature. The dashboard displays statistics, recent projects, user tasks, and upcoming meetings in a single view.

**Estimated Time**: 12-16 hours

---

## Prerequisites

### Required Knowledge
- Go programming with Fiber v2 framework
- Svelte 5 with runes syntax
- PostgreSQL and SQL
- REST API design
- JWT authentication

### Development Environment
- Go 1.21+ installed
- Node.js 18+ and npm installed
- PostgreSQL 14+ running locally or accessible
- VS Code or preferred editor
- Git for version control

### Existing Codebase
Ensure you're familiar with:
- Backend layered architecture (handlers → services → repositories)
- Svelte 5 component patterns with `$props()`, `$state()`, `$derived()`
- Existing authentication flow (JWT cookies)
- Database migration process using `go run ./cmd/migrate`

---

## Implementation Order

Follow this sequence for logical dependency flow:

1. **Database Migration** (30 min)
2. **Backend Models** (45 min)
3. **Backend Repositories** (2 hours)
4. **Backend Services** (2.5 hours)
5. **Backend Handlers & Routes** (1.5 hours)
6. **Frontend API Client** (45 min)
7. **Frontend Components** (4 hours)
8. **Frontend Dashboard Page** (2 hours)
9. **Testing & Refinement** (2 hours)

---

## Step-by-Step Implementation

### Step 1: Database Migration

**File**: `backend/migration/005_add_dashboard_meetings.sql`

**Action**: Create migration file with tables and indexes

```bash
cd backend/migration
touch 005_add_dashboard_meetings.sql
```

**Content**: Copy SQL from [data-model.md](data-model.md#database-migration)

**Key Tables**:
- `meetings`: Stores meeting information
- `meeting_attendees`: Junction table for attendees

**Run Migration**:
```bash
cd backend
go run ./cmd/migrate
```

**Verify**:
```sql
-- Connect to PostgreSQL
psql -U your_user -d your_database

-- Check tables exist
\dt meetings
\dt meeting_attendees

-- Check indexes
\di idx_meetings_*
\di idx_meeting_attendees_*
```

**Troubleshooting**:
- If migration fails, check PostgreSQL connection in `backend/config/database.go`
- Ensure UUID extension enabled: `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
- Check for naming conflicts with existing tables

---

### Step 2: Backend Models

**File**: `backend/models/meeting.go`

**Action**: Define Meeting and MeetingAttendee structs

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Meeting struct {
    ID              uuid.UUID  `json:"id"`
    Title           string     `json:"title"`
    Description     *string    `json:"description"`
    MeetingDate     time.Time  `json:"meeting_date"`
    DurationMinutes int        `json:"duration_minutes"`
    ProjectID       *uuid.UUID `json:"project_id"`
    CreatedBy       uuid.UUID  `json:"created_by"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}

type MeetingAttendee struct {
    MeetingID      uuid.UUID `json:"meeting_id"`
    UserID         uuid.UUID `json:"user_id"`
    ResponseStatus string    `json:"response_status"`
    AddedAt        time.Time `json:"added_at"`
}

// Response DTOs
type MeetingWithAttendees struct {
    Meeting
    Attendees      []UserSummary `json:"attendees"`
    TotalAttendees int           `json:"total_attendees"`
}

type UserSummary struct {
    ID       uuid.UUID `json:"id"`
    FullName string    `json:"full_name"`
}
```

**Add Dashboard DTOs**:

**File**: `backend/models/dashboard.go`

```go
package models

type DashboardResponse struct {
    Statistics     DashboardStatistics `json:"statistics"`
    RecentProjects []ProjectCard       `json:"recent_projects"`
    UserTasks      []TaskSummary       `json:"user_tasks"`
    NextMeeting    *MeetingWithAttendees `json:"next_meeting"`
}

type DashboardStatistics struct {
    ActiveProjects     StatValue `json:"active_projects"`
    PendingTasks       StatValue `json:"pending_tasks"`
    TeamMembers        StatValue `json:"team_members"`
    UpcomingDeadlines  StatValue `json:"upcoming_deadlines"`
}

type StatValue struct {
    Current  int `json:"current"`
    Previous int `json:"previous"`
    Change   int `json:"change"`
}

type ProjectCard struct {
    ID           uuid.UUID     `json:"id"`
    Name         string        `json:"name"`
    Client       string        `json:"client"`
    Status       string        `json:"status"`
    Progress     int           `json:"progress"`
    DueDate      *time.Time    `json:"due_date"`
    UpdatedAt    time.Time     `json:"updated_at"`
    TeamMembers  []UserSummary `json:"team_members"`
    TotalMembers int           `json:"total_members"`
}

type TaskSummary struct {
    ID            uuid.UUID  `json:"id"`
    Title         string     `json:"title"`
    ProjectName   string     `json:"project_name"`
    ProjectID     uuid.UUID  `json:"project_id"`
    Priority      int        `json:"priority"`
    PriorityLabel string     `json:"priority_label"`
    DueDate       *time.Time `json:"due_date"`
    Status        string     `json:"status"`
    CreatedAt     time.Time  `json:"created_at"`
}
```

**Note**: Use pointer types (`*string`, `*uuid.UUID`, `*time.Time`) for nullable fields

---

### Step 3: Backend Repositories

**File**: `backend/repositories/meeting_repository.go`

**Action**: Create repository for meeting data access

```go
package repositories

import (
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "your-module-path/models"
)

type MeetingRepository struct {
    db *pgxpool.Pool
}

func NewMeetingRepository(db *pgxpool.Pool) *MeetingRepository {
    return &MeetingRepository{db: db}
}

// GetNextMeetingForUser retrieves next upcoming meeting for user
func (r *MeetingRepository) GetNextMeetingForUser(ctx context.Context, userID uuid.UUID) (*models.MeetingWithAttendees, error) {
    // Query for next meeting where user is attendee or creator
    // Join with meeting_attendees and users tables
    // Filter: meeting_date > NOW() AND meeting_date <= NOW() + 7 days
    // Order by meeting_date ASC LIMIT 1
    // Implementation details in data-model.md
}

// CreateMeeting inserts new meeting
func (r *MeetingRepository) CreateMeeting(ctx context.Context, meeting *models.Meeting) error {
    // INSERT INTO meetings ...
}

// AddAttendees adds users to meeting
func (r *MeetingRepository) AddAttendees(ctx context.Context, meetingID uuid.UUID, userIDs []uuid.UUID) error {
    // INSERT INTO meeting_attendees ...
}

// More methods as needed...
```

**File**: `backend/repositories/dashboard_repository.go`

**Action**: Create repository for dashboard-specific queries

```go
package repositories

import (
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "your-module-path/models"
)

type DashboardRepository struct {
    db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
    return &DashboardRepository{db: db}
}

// GetStatistics calculates dashboard statistics
func (r *DashboardRepository) GetStatistics(ctx context.Context, userID uuid.UUID, userRole string) (*models.DashboardStatistics, error) {
    // Calculate current and previous (7 days ago) counts
    // Filter by user's visible projects based on role
    // Return StatValue structs with current, previous, change
}

// GetRecentProjects fetches up to 4 recent projects
func (r *DashboardRepository) GetRecentProjects(ctx context.Context, userID uuid.UUID, userRole string, limit int) ([]models.ProjectCard, error) {
    // SELECT projects with JOIN to tasks for progress calculation
    // JOIN to get team members (DISTINCT users from tasks.assigned_user_id)
    // Filter by user role and project visibility
    // ORDER BY updated_at DESC LIMIT limit
}

// GetUserTasks fetches user's top priority tasks
func (r *DashboardRepository) GetUserTasks(ctx context.Context, userID uuid.UUID, limit int) ([]models.TaskSummary, error) {
    // SELECT tasks WHERE assigned_user_id = userID AND status != 'done'
    // JOIN with projects to get project_name
    // ORDER BY priority DESC, due_date ASC, created_at ASC
    // LIMIT limit
}
```

**Key Implementation Notes**:
- Use `pgx` connection pool (`*pgxpool.Pool`)
- Handle NULL values with pointer types
- Use context for cancellation and timeouts
- Return appropriate errors (wrap with context when needed)
- Apply role-based filtering in SQL WHERE clauses

---

### Step 4: Backend Services

**File**: `backend/services/meeting_service.go`

**Action**: Create business logic layer for meetings

```go
package services

import (
    "context"
    "errors"
    "time"
    "github.com/google/uuid"
    "your-module-path/models"
    "your-module-path/repositories"
)

type MeetingService struct {
    meetingRepo *repositories.MeetingRepository
    userRepo    *repositories.UserRepository
}

func NewMeetingService(meetingRepo *repositories.MeetingRepository, userRepo *repositories.UserRepository) *MeetingService {
    return &MeetingService{
        meetingRepo: meetingRepo,
        userRepo:    userRepo,
    }
}

// GetNextMeetingForUser retrieves next upcoming meeting
func (s *MeetingService) GetNextMeetingForUser(ctx context.Context, userID uuid.UUID) (*models.MeetingWithAttendees, error) {
    return s.meetingRepo.GetNextMeetingForUser(ctx, userID)
}

// CreateMeeting creates a new meeting with validation
func (s *MeetingService) CreateMeeting(ctx context.Context, userID uuid.UUID, input *CreateMeetingInput) (*models.MeetingWithAttendees, error) {
    // Validate input
    if err := validateMeetingInput(input); err != nil {
        return nil, err
    }
    
    // Check meeting_date is in future
    if input.MeetingDate.Before(time.Now()) {
        return nil, errors.New("meeting date must be in the future")
    }
    
    // Verify all attendee IDs exist
    for _, attendeeID := range input.AttendeeIDs {
        if _, err := s.userRepo.GetByID(ctx, attendeeID); err != nil {
            return nil, errors.New("invalid attendee ID")
        }
    }
    
    // Create meeting
    meeting := &models.Meeting{
        Title:           input.Title,
        Description:     input.Description,
        MeetingDate:     input.MeetingDate,
        DurationMinutes: input.DurationMinutes,
        ProjectID:       input.ProjectID,
        CreatedBy:       userID,
    }
    
    if err := s.meetingRepo.CreateMeeting(ctx, meeting); err != nil {
        return nil, err
    }
    
    // Add attendees (including creator with 'accepted' status)
    // ...
    
    return s.meetingRepo.GetMeetingByID(ctx, meeting.ID)
}
```

**File**: `backend/services/dashboard_service.go`

**Action**: Create service for dashboard data aggregation

```go
package services

import (
    "context"
    "github.com/google/uuid"
    "your-module-path/models"
    "your-module-path/repositories"
)

type DashboardService struct {
    dashboardRepo *repositories.DashboardRepository
    meetingRepo   *repositories.MeetingRepository
}

func NewDashboardService(dashboardRepo *repositories.DashboardRepository, meetingRepo *repositories.MeetingRepository) *DashboardService {
    return &DashboardService{
        dashboardRepo: dashboardRepo,
        meetingRepo:   meetingRepo,
    }
}

// GetDashboardData retrieves all dashboard information
func (s *DashboardService) GetDashboardData(ctx context.Context, userID uuid.UUID, userRole string) (*models.DashboardResponse, error) {
    // Fetch statistics
    stats, err := s.dashboardRepo.GetStatistics(ctx, userID, userRole)
    if err != nil {
        return nil, err
    }
    
    // Fetch recent projects (limit 4)
    projects, err := s.dashboardRepo.GetRecentProjects(ctx, userID, userRole, 4)
    if err != nil {
        return nil, err
    }
    
    // Fetch user tasks (limit 5)
    tasks, err := s.dashboardRepo.GetUserTasks(ctx, userID, 5)
    if err != nil {
        return nil, err
    }
    
    // Fetch next meeting (nullable)
    meeting, err := s.meetingRepo.GetNextMeetingForUser(ctx, userID)
    if err != nil && err != pgx.ErrNoRows {
        return nil, err
    }
    
    return &models.DashboardResponse{
        Statistics:     *stats,
        RecentProjects: projects,
        UserTasks:      tasks,
        NextMeeting:    meeting,
    }, nil
}
```

**Validation Functions**:
```go
func validateMeetingInput(input *CreateMeetingInput) error {
    if len(input.Title) == 0 || len(input.Title) > 200 {
        return errors.New("title must be 1-200 characters")
    }
    if input.DurationMinutes <= 0 || input.DurationMinutes > 1440 {
        return errors.New("duration must be 1-1440 minutes")
    }
    if len(input.AttendeeIDs) == 0 {
        return errors.New("at least one attendee required")
    }
    return nil
}
```

---

### Step 5: Backend Handlers & Routes

**File**: `backend/handlers/dashboard_handler.go`

**Action**: Create HTTP handler for dashboard endpoint

```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "your-module-path/models"
    "your-module-path/services"
)

type DashboardHandler struct {
    dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
    return &DashboardHandler{
        dashboardService: dashboardService,
    }
}

// GetDashboard handles GET /api/dashboard
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
    // Extract user from JWT (set by auth middleware)
    user := c.Locals("user").(*models.User)
    
    // Fetch dashboard data
    data, err := h.dashboardService.GetDashboardData(c.Context(), user.ID, user.Role)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "Internal Server Error",
            "message": "Failed to fetch dashboard data",
        })
    }
    
    return c.JSON(data)
}
```

**File**: `backend/handlers/meeting_handler.go`

**Action**: Create HTTP handler for meeting endpoints

```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "your-module-path/models"
    "your-module-path/services"
)

type MeetingHandler struct {
    meetingService *services.MeetingService
}

func NewMeetingHandler(meetingService *services.MeetingService) *MeetingHandler {
    return &MeetingHandler{
        meetingService: meetingService,
    }
}

// GetNextMeeting handles GET /api/meetings/next
func (h *MeetingHandler) GetNextMeeting(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    meeting, err := h.meetingService.GetNextMeetingForUser(c.Context(), user.ID)
    if err != nil {
        if err == pgx.ErrNoRows {
            return c.SendStatus(fiber.StatusNoContent)
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "Internal Server Error",
            "message": "Failed to fetch meeting",
        })
    }
    
    return c.JSON(meeting)
}

// CreateMeeting handles POST /api/meetings
func (h *MeetingHandler) CreateMeeting(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    var input services.CreateMeetingInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   "Bad Request",
            "message": "Invalid request body",
        })
    }
    
    meeting, err := h.meetingService.CreateMeeting(c.Context(), user.ID, &input)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   "Bad Request",
            "message": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(meeting)
}
```

**File**: `backend/routes/routes.go`

**Action**: Register new routes

```go
// Add to existing SetupRoutes function
func SetupRoutes(app *fiber.App, /* ...dependencies... */) {
    api := app.Group("/api")
    
    // Existing routes...
    
    // Dashboard routes (protected)
    dashboardHandler := handlers.NewDashboardHandler(dashboardService)
    api.Get("/dashboard", authMiddleware, dashboardHandler.GetDashboard)
    
    // Meeting routes (protected)
    meetingHandler := handlers.NewMeetingHandler(meetingService)
    meetings := api.Group("/meetings", authMiddleware)
    meetings.Get("/next", meetingHandler.GetNextMeeting)
    meetings.Post("/", meetingHandler.CreateMeeting)
    meetings.Get("/", meetingHandler.ListMeetings)
    meetings.Get("/:id", meetingHandler.GetMeeting)
}
```

**Update**: `backend/main.go` to wire dependencies

```go
// Initialize repositories
dashboardRepo := repositories.NewDashboardRepository(db)
meetingRepo := repositories.NewMeetingRepository(db)

// Initialize services
dashboardService := services.NewDashboardService(dashboardRepo, meetingRepo)
meetingService := services.NewMeetingService(meetingRepo, userRepo)

// Pass to routes.SetupRoutes()
```

---

### Step 6: Frontend API Client

**File**: `frontend/src/lib/api/dashboard.js`

**Action**: Create API client functions

```javascript
import { API_BASE_URL } from './config';

/**
 * Fetch all dashboard data
 * @returns {Promise<DashboardResponse>}
 */
export async function getDashboardData() {
  const response = await fetch(`${API_BASE_URL}/dashboard`, {
    credentials: 'include', // Include cookies (JWT)
    headers: {
      'Content-Type': 'application/json',
    },
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch dashboard data');
  }
  
  return await response.json();
}

/**
 * Mark task as complete
 * @param {string} taskId - UUID of task
 * @returns {Promise<Task>}
 */
export async function completeTask(taskId) {
  const response = await fetch(`${API_BASE_URL}/tasks/${taskId}`, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ status: 'done' }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to complete task');
  }
  
  return await response.json();
}
```

**File**: `frontend/src/lib/api/meetings.js`

```javascript
import { API_BASE_URL } from './config';

/**
 * Fetch next upcoming meeting
 * @returns {Promise<Meeting|null>}
 */
export async function getNextMeeting() {
  const response = await fetch(`${API_BASE_URL}/meetings/next`, {
    credentials: 'include',
  });
  
  if (response.status === 204) {
    return null; // No meetings
  }
  
  if (!response.ok) {
    throw new Error('Failed to fetch meeting');
  }
  
  return await response.json();
}

/**
 * Create new meeting
 * @param {CreateMeetingInput} data
 * @returns {Promise<Meeting>}
 */
export async function createMeeting(data) {
  const response = await fetch(`${API_BASE_URL}/meetings`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to create meeting');
  }
  
  return await response.json();
}
```

---

### Step 7: Frontend Components

Create reusable Svelte 5 components following existing patterns.

**File**: `frontend/src/components/Avatar.svelte`

```svelte
<script>
  // Extract initials and generate color from user name/ID
  let { user, size = 'md' } = $props();
  
  const sizeClasses = {
    sm: 'w-7 h-7 text-xs',
    md: 'w-8 h-8 text-sm',
    lg: 'w-10 h-10 text-base',
  };
  
  const colors = [
    'bg-indigo-500',
    'bg-purple-500',
    'bg-blue-500',
    'bg-pink-500',
    'bg-orange-500',
    'bg-green-500',
    'bg-red-500',
    'bg-teal-500',
  ];
  
  function getInitials(fullName) {
    const parts = fullName.trim().split(' ');
    if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }
  
  function getColorClass(userId) {
    const hash = userId.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return colors[hash % colors.length];
  }
  
  const initials = $derived(getInitials(user.full_name));
  const colorClass = $derived(getColorClass(user.id));
</script>

<div class="rounded-full {sizeClasses[size]} {colorClass} flex items-center justify-center text-white font-bold border-2 border-white">
  {initials}
</div>
```

**File**: `frontend/src/components/StatCard.svelte`

```svelte
<script>
  let { title, value, change, icon, iconColor = 'blue' } = $props();
  
  const changeColor = $derived(change > 0 ? 'text-green-600 bg-green-50' : change < 0 ? 'text-red-600 bg-red-50' : 'text-slate-600 bg-slate-100');
  const changePrefix = $derived(change > 0 ? '+' : '');
</script>

<div class="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md transition-shadow duration-200">
  <div class="flex items-start justify-between mb-4">
    <div class="p-3 rounded-xl bg-{iconColor}-50">
      <i data-lucide={icon} class="w-6 h-6 text-{iconColor}-600"></i>
    </div>
    <span class="text-xs font-medium px-2 py-1 rounded-full {changeColor}">
      {changePrefix}{change}
    </span>
  </div>
  <div>
    <h3 class="text-2xl font-bold text-slate-800">{value}</h3>
    <p class="text-sm text-slate-500 font-medium">{title}</p>
  </div>
</div>
```

**File**: `frontend/src/components/ProjectCard.svelte`

```svelte
<script>
  import Avatar from './Avatar.svelte';
  import { formatJalaliDate } from '../lib/utils';
  
  let { project, onclick } = $props();
  
  const statusColors = {
    'Planning': 'bg-slate-100 text-slate-700',
    'In Progress': 'bg-blue-100 text-blue-700',
    'On Track': 'bg-green-100 text-green-700',
    'Review': 'bg-purple-100 text-purple-700',
  };
</script>

<div
  class="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm hover:border-indigo-100 transition-colors group cursor-pointer"
  onclick={onclick}
  role="button"
  tabindex="0"
>
  <div class="flex justify-between items-start mb-4">
    <span class="text-xs font-semibold px-2.5 py-1 rounded-md {statusColors[project.status]}">
      {project.status}
    </span>
  </div>
  
  <h3 class="font-bold text-slate-800 text-lg mb-1 group-hover:text-indigo-600 transition-colors line-clamp-2">
    {project.name}
  </h3>
  <p class="text-sm text-slate-500 mb-4">{project.client}</p>
  
  <div class="mb-4">
    <div class="flex justify-between text-xs font-medium text-slate-500 mb-1">
      <span>پیشرفت</span>
      <span>{project.progress}%</span>
    </div>
    <div class="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
      <div class="bg-indigo-600 h-2 rounded-full" style="width: {project.progress}%"></div>
    </div>
  </div>
  
  <div class="flex items-center justify-between pt-4 border-t border-slate-50">
    <div class="flex -space-x-2">
      {#each project.team_members.slice(0, 3) as member}
        <Avatar user={member} size="sm" />
      {/each}
      {#if project.total_members > 3}
        <div class="w-7 h-7 rounded-full bg-slate-100 border-2 border-white flex items-center justify-center text-xs font-bold text-slate-500">
          +{project.total_members - 3}
        </div>
      {/if}
    </div>
    <div class="flex items-center gap-1 text-slate-400 text-xs font-medium">
      <i data-lucide="clock" class="w-3.5 h-3.5"></i>
      {formatJalaliDate(project.due_date)}
    </div>
  </div>
</div>
```

**File**: `frontend/src/components/TaskListItem.svelte`

```svelte
<script>
  let { task, onComplete } = $props();
  
  const priorityColors = {
    4: 'bg-red-100 text-red-700',
    3: 'bg-orange-100 text-orange-700',
    2: 'bg-blue-100 text-blue-700',
    1: 'bg-slate-100 text-slate-700',
  };
  
  let completed = $state(false);
  
  async function handleCheck() {
    completed = true;
    await onComplete(task.id);
    // Fade out handled by parent
  }
</script>

<div class="flex items-center justify-between p-4 hover:bg-slate-50 rounded-xl transition-colors border-b border-slate-50 last:border-0"
     class:opacity-50={completed}>
  <div class="flex items-center gap-4">
    <button
      onclick={handleCheck}
      disabled={completed}
      class="text-slate-300 hover:text-indigo-600 transition-colors"
    >
      <div class="w-5 h-5 border-2 border-current rounded-md flex items-center justify-center">
        {#if completed}
          <i data-lucide="check" class="w-3 h-3"></i>
        {/if}
      </div>
    </button>
    <div>
      <p class="text-sm font-semibold text-slate-700" class:line-through={completed}>
        {task.title}
      </p>
      <p class="text-xs text-slate-400">{task.project_name}</p>
    </div>
  </div>
  <div class="flex items-center gap-4">
    <span class="px-2 py-1 rounded text-xs font-medium {priorityColors[task.priority]}">
      {task.priority_label}
    </span>
    <span class="text-xs text-slate-500 w-16 text-right">
      {formatJalaliDate(task.due_date, 'short')}
    </span>
  </div>
</div>
```

**File**: `frontend/src/components/MeetingCard.svelte`

```svelte
<script>
  import Avatar from './Avatar.svelte';
  import { formatJalaliDate } from '../lib/utils';
  
  let { meeting } = $props();
</script>

{#if meeting}
  <div class="bg-gradient-to-br from-indigo-600 to-purple-700 rounded-2xl p-6 text-white shadow-lg shadow-indigo-200">
    <h4 class="font-bold text-lg mb-2">{meeting.title}</h4>
    <p class="text-indigo-100 text-sm mb-6 line-clamp-2">{meeting.description}</p>
    <div class="flex items-center justify-between">
      <div class="flex -space-x-2">
        {#each meeting.attendees.slice(0, 3) as attendee}
          <div class="w-8 h-8 rounded-full bg-indigo-400 border-2 border-indigo-600"></div>
        {/each}
      </div>
      <span class="text-sm font-bold bg-white/20 px-3 py-1 rounded-lg backdrop-blur-sm">
        {formatJalaliDate(meeting.meeting_date, 'time')}
      </span>
    </div>
  </div>
{/if}
```

---

### Step 8: Frontend Dashboard Page

**File**: `frontend/src/components/Dashboard.svelte`

```svelte
<script>
  import { onMount } from 'svelte';
  import { getDashboardData, completeTask } from '../lib/api/dashboard';
  import StatCard from './StatCard.svelte';
  import ProjectCard from './ProjectCard.svelte';
  import TaskListItem from './TaskListItem.svelte';
  import MeetingCard from './MeetingCard.svelte';
  
  let dashboardData = $state(null);
  let loading = $state(true);
  let error = $state(null);
  
  async function loadDashboard() {
    try {
      loading = true;
      dashboardData = await getDashboardData();
      error = null;
    } catch (err) {
      error = err.message;
      console.error('Failed to load dashboard:', err);
    } finally {
      loading = false;
    }
  }
  
  async function handleTaskComplete(taskId) {
    try {
      await completeTask(taskId);
      // Remove completed task after 2 seconds
      setTimeout(() => {
        dashboardData.user_tasks = dashboardData.user_tasks.filter(t => t.id !== taskId);
      }, 2000);
    } catch (err) {
      console.error('Failed to complete task:', err);
      alert('خطا در تکمیل وظیفه');
    }
  }
  
  function navigateToProject(projectId) {
    window.location.href = `/projects/${projectId}`;
  }
  
  // Auto-refresh every 30 seconds
  $effect(() => {
    loadDashboard();
    
    const interval = setInterval(loadDashboard, 30000);
    
    return () => clearInterval(interval);
  });
</script>

<div class="max-w-7xl mx-auto space-y-8 p-8">
  {#if loading && !dashboardData}
    <div class="text-center py-12">
      <p class="text-slate-500">در حال بارگذاری...</p>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-700">خطا: {error}</p>
    </div>
  {:else if dashboardData}
    <!-- Statistics Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatCard
        title="پروژه‌های فعال"
        value={dashboardData.statistics.active_projects.current}
        change={dashboardData.statistics.active_projects.change}
        icon="folder-kanban"
        iconColor="blue"
      />
      <StatCard
        title="وظایف در انتظار"
        value={dashboardData.statistics.pending_tasks.current}
        change={dashboardData.statistics.pending_tasks.change}
        icon="check-square"
        iconColor="orange"
      />
      <StatCard
        title="اعضای تیم"
        value={dashboardData.statistics.team_members.current}
        change={dashboardData.statistics.team_members.change}
        icon="users"
        iconColor="purple"
      />
      <StatCard
        title="ددلاین‌های پیش رو"
        value={dashboardData.statistics.upcoming_deadlines.current}
        change={dashboardData.statistics.upcoming_deadlines.change}
        icon="calendar"
        iconColor="red"
      />
    </div>
    
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Projects Section -->
      <div class="lg:col-span-2 space-y-6">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-slate-800">پروژه‌های اخیر</h3>
          <a href="/projects" class="text-sm text-indigo-600 font-medium hover:text-indigo-700 flex items-center gap-1">
            مشاهده همه
            <i data-lucide="chevron-left" class="w-4 h-4"></i>
          </a>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          {#each dashboardData.recent_projects as project (project.id)}
            <ProjectCard {project} onclick={() => navigateToProject(project.id)} />
          {/each}
        </div>
      </div>
      
      <!-- Sidebar: Tasks & Meeting -->
      <div class="lg:col-span-1 space-y-6">
        <!-- Tasks Section -->
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-slate-800">وظایف شما</h3>
        </div>
        <div class="bg-white rounded-2xl border border-slate-100 shadow-sm p-2">
          {#each dashboardData.user_tasks as task (task.id)}
            <TaskListItem {task} onComplete={handleTaskComplete} />
          {/each}
          <button class="w-full py-3 text-sm text-slate-500 hover:text-indigo-600 font-medium flex items-center justify-center gap-2 transition-colors border-t border-slate-50 mt-2">
            <i data-lucide="plus" class="w-4 h-4"></i>
            افزودن وظیفه جدید
          </button>
        </div>
        
        <!-- Meeting Card -->
        <MeetingCard meeting={dashboardData.next_meeting} />
      </div>
    </div>
  {/if}
</div>

<script context="module">
  // Initialize Lucide icons after mount
  import lucide from 'lucide';
  onMount(() => {
    lucide.createIcons();
  });
</script>
```

**Update**: `frontend/src/App.svelte` to add dashboard route

```svelte
<script>
  import { Router, Route } from 'svelte-routing';
  import Dashboard from './components/Dashboard.svelte';
  // ... other imports
</script>

<Router>
  <Route path="/" component={Dashboard} />
  <Route path="/dashboard" component={Dashboard} />
  <!-- ... other routes -->
</Router>
```

---

### Step 9: Utility Functions

**File**: `frontend/src/lib/utils.js`

```javascript
import jalaaliMoment from 'jalali-moment';

/**
 * Format date in Jalali calendar
 * @param {string|Date} date - Date to format
 * @param {string} format - 'full'|'short'|'time'
 * @returns {string} Formatted date
 */
export function formatJalaliDate(date, format = 'short') {
  if (!date) return '';
  
  const m = jalaaliMoment(date);
  
  const now = jalaaliMoment();
  const tomorrow = jalaaliMoment().add(1, 'day');
  
  if (m.isSame(now, 'day')) {
    return 'امروز';
  }
  if (m.isSame(tomorrow, 'day')) {
    return 'فردا';
  }
  
  switch (format) {
    case 'full':
      return m.format('jYYYY/jMM/jDD');
    case 'short':
      return m.format('jMM/jDD');
    case 'time':
      return m.format('HH:mm');
    default:
      return m.format('jYYYY/jMM/jDD');
  }
}
```

---

## Testing

### Manual Testing Checklist

**Dashboard Load**:
- [ ] Dashboard loads within 2 seconds
- [ ] All 4 statistics display correctly
- [ ] Statistics show change indicators with correct colors
- [ ] Up to 4 recent projects displayed
- [ ] Projects show progress bars matching task completion
- [ ] Up to 5 user tasks displayed
- [ ] Tasks sorted by priority then due date
- [ ] Next meeting card appears (if meeting exists)
- [ ] Meeting card hidden (if no upcoming meetings)

**Task Interactions**:
- [ ] Clicking task checkbox marks task complete
- [ ] Completed task fades out after 2 seconds
- [ ] Task list refreshes showing next task
- [ ] Pending tasks count decreases

**Project Interactions**:
- [ ] Clicking project card navigates to project details
- [ ] Project avatars display initials correctly
- [ ] Avatar colors consistent for same user

**Auto-refresh**:
- [ ] Dashboard auto-refreshes every 30 seconds
- [ ] No flickering during refresh
- [ ] User interactions not disrupted by refresh

**Role-based Display**:
- [ ] Team member sees only assigned projects
- [ ] Project manager sees managed projects
- [ ] Admin sees all projects
- [ ] Statistics filtered by user role

**Edge Cases**:
- [ ] Empty state when no projects exist
- [ ] Empty state when no tasks assigned
- [ ] No meeting card when no meetings scheduled
- [ ] Long project names truncate with ellipsis
- [ ] Overdue tasks highlighted appropriately

### API Testing

Use curl or Postman to test endpoints:

```bash
# Get dashboard data
curl -X GET http://localhost:3000/api/dashboard \
  -H "Cookie: access_token=YOUR_JWT_TOKEN"

# Complete task
curl -X PATCH http://localhost:3000/api/tasks/TASK_UUID \
  -H "Cookie: access_token=YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'

# Get next meeting
curl -X GET http://localhost:3000/api/meetings/next \
  -H "Cookie: access_token=YOUR_JWT_TOKEN"

# Create meeting
curl -X POST http://localhost:3000/api/meetings \
  -H "Cookie: access_token=YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Sync",
    "description": "Weekly standup",
    "meeting_date": "2023-10-20T10:00:00Z",
    "duration_minutes": 30,
    "attendee_ids": ["USER_UUID_1", "USER_UUID_2"]
  }'
```

### Integration Testing

Run backend tests:
```bash
cd backend
go test ./... -v
```

Run frontend tests (if test framework setup):
```bash
cd frontend
npm test
```

---

## Troubleshooting

### Common Issues

**Database Connection Error**:
- Check PostgreSQL is running: `sudo systemctl status postgresql`
- Verify connection string in `backend/config/database.go`
- Check user permissions: `GRANT ALL ON DATABASE your_db TO your_user;`

**Migration Fails**:
- Check migration SQL syntax
- Ensure UUID extension enabled: `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
- Verify no table name conflicts
- Check foreign key references exist

**CORS Errors**:
- Ensure frontend and backend on same origin, or
- Add CORS middleware in backend:
  ```go
  app.Use(cors.New(cors.Config{
      AllowOrigins:     "http://localhost:5173",
      AllowCredentials: true,
  }))
  ```

**JWT Authentication Fails**:
- Check cookie settings (httpOnly, SameSite, Secure)
- Verify token expiration
- Ensure auth middleware runs before dashboard handler

**Dashboard Data Not Loading**:
- Check browser console for errors
- Verify API endpoint returns 200 status
- Check network tab for request/response
- Verify JWT token present in cookie

**Auto-refresh Not Working**:
- Check `$effect` hook syntax (Svelte 5)
- Verify interval cleanup on component unmount
- Console log to ensure refresh triggered

**Persian Dates Not Displaying**:
- Verify `jalali-moment` installed: `npm install jalali-moment`
- Check import statement in utils.js
- Test formatJalaliDate function independently

---

## Performance Optimization

### Backend
- **Indexes**: Ensure all indexes created (run migration)
- **Connection Pooling**: Configure pgxpool max connections (default 4)
- **Query Optimization**: Use EXPLAIN ANALYZE to check query plans
- **Caching**: Consider adding Redis for statistics (future enhancement)

### Frontend
- **Lazy Loading**: Only load dashboard on /dashboard route
- **Debouncing**: Debounce rapid refresh requests
- **Code Splitting**: Use dynamic imports for large components
- **Image Optimization**: Use lazy loading for avatars if images added

---

## Security Considerations

- **Authentication**: All endpoints protected by JWT middleware
- **Authorization**: Service layer filters data by user role
- **SQL Injection**: Use parameterized queries (pgx automatically handles)
- **XSS**: Svelte escapes HTML by default
- **CSRF**: SameSite cookie attribute prevents CSRF
- **Rate Limiting**: Consider adding rate limit middleware (future)

---

## Deployment

### Production Build

**Backend**:
```bash
cd backend
go build -o dashboard-server
./dashboard-server
```

**Frontend**:
```bash
cd frontend
npm run build
# Serve dist/ directory with nginx or CDN
```

### Environment Variables

**Backend** (`backend/.env`):
```
DATABASE_URL=postgres://user:pass@localhost:5432/dbname
JWT_SECRET=your-secret-key
PORT=3000
```

**Frontend** (`frontend/.env.production`):
```
VITE_API_BASE_URL=https://api.yourdomain.com
```

### Docker Deployment

See existing `docker-compose.yml` and add dashboard service if needed.

---

## Next Steps

After dashboard implementation:

1. **Testing**: Write comprehensive unit and integration tests
2. **Documentation**: Update README with dashboard usage
3. **User Feedback**: Collect feedback from project managers
4. **Performance Monitoring**: Add logging and metrics
5. **Feature Enhancements**:
   - Dashboard customization
   - Data visualization (charts)
   - Advanced filtering
   - Export functionality
   - Mobile optimization

---

## Resources

**Documentation**:
- [Svelte 5 Runes](https://svelte.dev/docs/svelte/runes)
- [Fiber v2 Docs](https://docs.gofiber.io/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Lucide Icons](https://lucide.dev/)

**Related Files**:
- [spec.md](spec.md) - Feature specification
- [research.md](research.md) - Technical research
- [data-model.md](data-model.md) - Database schema
- [api-contract.md](contracts/api-contract.md) - API specification

---

**Questions or Issues?**  
Refer to research.md for design decisions and alternatives considered.
