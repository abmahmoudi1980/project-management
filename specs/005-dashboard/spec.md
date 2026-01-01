# Feature Specification: Project Manager Dashboard

**Feature ID**: 005  
**Feature Name**: dashboard  
**Created**: 2026-01-01  
**Status**: Draft  

---

## Overview

### Purpose

Provide users with a centralized dashboard that displays a comprehensive overview of all active projects, pending tasks, team activities, and upcoming deadlines. The dashboard serves as the primary landing page after login, enabling users to quickly assess project status, identify bottlenecks, and prioritize their work.

### Business Value

- **Increased productivity**: Users can immediately see their most urgent tasks and project status without navigating through multiple pages
- **Better decision making**: Real-time statistics and progress indicators help managers identify projects that need attention
- **Reduced context switching**: All critical information is available in one place, reducing time spent searching for data
- **Improved team coordination**: Visibility into team member assignments and upcoming meetings facilitates better collaboration

### User Personas

- **Project Manager**: Needs to monitor overall project health, track deadlines, and ensure team capacity is optimized
- **Team Member**: Wants to see assigned tasks, upcoming work, and project priorities
- **Admin**: Requires high-level view of all organizational activities and team performance

---

## User Scenarios & Testing

### Primary User Flows

#### Scenario 1: Morning Check-in (Project Manager)
**Actor**: Sarah, Project Manager  
**Goal**: Review overnight updates and plan the day  

**Steps**:
1. Sarah logs into the system
2. Dashboard immediately displays with current statistics
3. She scans the stat cards to see: 12 active projects (+2 new), 48 pending tasks (-5 completed), 24 team members (+4 new), 7 upcoming deadlines
4. She notices the "Website Redesign" project is at 75% progress
5. She clicks on her urgent task "Design Homepage Mockups" marked as High priority
6. She reviews the team meeting card showing a sync at 10:00 AM

**Expected Outcome**: Sarah has a complete picture of project status and her priorities within 30 seconds

#### Scenario 2: Task Management (Team Member)
**Actor**: Mike, Frontend Developer  
**Goal**: Identify and complete daily tasks  

**Steps**:
1. Mike opens the dashboard
2. He views the "Your Tasks" section on the right side
3. He sees 3 tasks: one marked Critical (Fix Navigation Bug), one High (Design Homepage Mockups), and one Medium (Content Review)
4. He notices the Critical task is due tomorrow
5. He checks the checkbox to mark the task as complete
6. The task moves to completed state and the counter updates

**Expected Outcome**: Mike can quickly identify task priorities and due dates, and mark tasks complete

#### Scenario 3: Project Drill-down (Project Manager)
**Actor**: Sarah, Project Manager  
**Goal**: Get detailed information about a specific project  

**Steps**:
1. Sarah views the "Recent Projects" section
2. She sees 4 project cards with status badges, progress bars, team avatars, and deadlines
3. She notices "Mobile App Alpha" has only 30% progress with deadline Nov 12
4. She clicks on the project card to view full details
5. System navigates to the project detail page

**Expected Outcome**: Sarah can quickly assess project health and navigate to details when needed

### Edge Cases

1. **No active projects**: Display empty state with "Create New Project" call-to-action
2. **No tasks assigned**: Show encouraging message "You're all caught up!" with option to view all tasks
3. **Overdue tasks**: Highlight overdue tasks in red with "Overdue" badge
4. **Large task count (50+)**: Show only top 5 tasks ordered by priority and due date, with "View All" link
5. **Long project names**: Truncate project names to 2 lines with ellipsis
6. **Multiple team meetings**: Show only the next upcoming meeting in the sidebar card

### Acceptance Criteria

#### Statistics Display
- [ ] Dashboard displays 4 stat cards: Active Projects, Pending Tasks, Team Members, Upcoming Deadlines
- [ ] Each stat card shows current count and change indicator (+/- with number)
- [ ] Change indicators display in green (positive) or red (negative)
- [ ] Statistics update in real-time when data changes elsewhere in the system

#### Project Cards
- [ ] Recent projects section displays up to 4 most recently updated projects
- [ ] Each project card shows: status badge, project name, client/organization name, progress percentage, progress bar, team member avatars (max 3 visible), deadline date
- [ ] Status badges use color coding: Planning (gray), In Progress (blue), On Track (green), Review (purple)
- [ ] Progress bar visually represents completion percentage
- [ ] Clicking project card navigates to project detail page

#### Task List
- [ ] "Your Tasks" section displays up to 5 highest priority tasks for logged-in user
- [ ] Each task shows: checkbox, task name, project name, priority badge, due date
- [ ] Priority badges use color coding: Critical (red), High (orange), Medium (blue), Low (gray)
- [ ] Checking task checkbox marks it as complete
- [ ] Tasks are sorted by: priority (Critical > High > Medium > Low), then due date (soonest first)
- [ ] "Add New Task" button opens task creation modal

#### Navigation & Layout
- [ ] Sidebar contains navigation links: Dashboard, Projects, My Tasks, Team, Settings
- [ ] Current page (Dashboard) is highlighted in sidebar
- [ ] Header contains search bar and notification bell with badge
- [ ] "New Project" button in header opens project creation modal
- [ ] Dashboard is responsive on mobile, tablet, and desktop
- [ ] User profile card in sidebar shows avatar, name, and role

#### Team Meeting Card
- [ ] Displays next upcoming team meeting with title, description, time, and attendee avatars
- [ ] Only shown if user has upcoming meetings within next 7 days
- [ ] Hidden if no meetings scheduled

---

## Functional Requirements

### Data Display Requirements

**FR-1: Statistics Overview**
- System shall display 4 key statistics in card format at top of dashboard
- Each statistic shall show current value and change from previous period
- Statistics shall include:
  - Active Projects: count of projects with status "In Progress" or "On Track"
  - Pending Tasks: count of tasks assigned to any team member with status not "Completed"
  - Team Members: count of active users in the system
  - Upcoming Deadlines: count of tasks and milestones due within next 7 days
- Change indicators shall calculate difference from same period 7 days ago

**FR-2: Recent Projects Grid**
- System shall display up to 4 most recently updated projects
- Projects shall be ordered by last update timestamp (most recent first)
- Each project card shall display:
  - Project name (max 50 characters, truncated with ellipsis)
  - Client/organization name
  - Status badge (color-coded)
  - Progress percentage
  - Visual progress bar
  - Up to 3 team member avatars (show "+N" if more members exist)
  - Deadline date in localized format

**FR-3: Task List**
- System shall display up to 5 tasks assigned to logged-in user
- Tasks shall be filtered to show only incomplete tasks
- Tasks shall be sorted by priority (Critical, High, Medium, Low) then due date
- Each task shall display:
  - Interactive checkbox
  - Task name (max 60 characters, truncated with ellipsis)
  - Parent project name
  - Priority badge (color-coded)
  - Due date (formatted as "Today", "Tomorrow", or date)
- Section shall show "Add New Task" button at bottom

**FR-4: Team Meeting Card**
- System shall display next upcoming meeting for logged-in user
- Meeting card shall show only if meeting exists within next 7 days
- Card shall display:
  - Meeting title
  - Meeting description (max 100 characters)
  - Meeting time
  - Up to 3 attendee avatars

### Interaction Requirements

**FR-5: Navigation**
- Clicking project card shall navigate to project detail page
- Clicking task name shall navigate to task detail page
- Clicking "View All" link in projects section shall navigate to projects list page
- Clicking "New Project" button shall open project creation modal
- Clicking "Add New Task" button shall open task creation modal
- Search bar shall filter across projects, tasks, and team members

**FR-6: Task Completion**
- Checking task checkbox shall mark task as complete
- Task completion shall update Pending Tasks statistic
- Completed task shall remain visible for 2 seconds with strikethrough, then fade out
- Task list shall refresh to show next highest priority task if available

**FR-7: Real-time Updates**
- Dashboard data shall refresh automatically every 30 seconds
- Statistics shall update without full page reload
- New tasks or projects created elsewhere shall appear on dashboard after next refresh

### Access Control Requirements

**FR-8: Role-based Display**
- Project Managers and Admins shall see all projects in their organization
- Team Members shall see only projects they are assigned to
- Task list shall show only tasks assigned to logged-in user
- Team meeting card shall show only meetings user is invited to

---

## Success Criteria

### Performance Metrics
- [ ] Dashboard loads all content within 2 seconds on standard broadband connection
- [ ] Dashboard renders correctly on viewports from 320px (mobile) to 2560px (desktop)
- [ ] Auto-refresh completes within 500ms without disrupting user interaction
- [ ] Page supports up to 100 concurrent users without degradation

### User Experience Metrics
- [ ] Users can identify their highest priority task within 5 seconds of viewing dashboard
- [ ] Users can assess overall project health without clicking any links
- [ ] 90% of users successfully complete their first task from dashboard within 1 minute
- [ ] Dashboard provides visual feedback for all user actions within 200ms

### Business Metrics
- [ ] Dashboard reduces time spent searching for task information by 50%
- [ ] Users access project details 30% faster compared to navigating from menu
- [ ] 80% of users start their session from dashboard page
- [ ] Task completion rate increases by 20% due to better visibility

---

## Key Entities

### Dashboard Statistics
- **Active Projects Count**: Integer count of active projects
- **Pending Tasks Count**: Integer count of incomplete tasks
- **Team Members Count**: Integer count of active users
- **Upcoming Deadlines Count**: Integer count of items due within 7 days
- **Change Indicators**: Positive or negative integer showing change from 7 days ago

### Project Card
- **Project ID**: Unique identifier
- **Project Name**: String (max 50 characters displayed)
- **Client Name**: String
- **Status**: Enum (Planning, In Progress, On Track, Review, Completed)
- **Progress**: Integer percentage (0-100)
- **Team Members**: Array of user objects (show max 3)
- **Deadline**: Date
- **Last Updated**: Timestamp

### Task Item
- **Task ID**: Unique identifier
- **Task Name**: String (max 60 characters displayed)
- **Project Name**: String
- **Priority**: Enum (Critical, High, Medium, Low)
- **Due Date**: Date
- **Completed**: Boolean
- **Assigned User**: User object

### Meeting Card
- **Meeting Title**: String
- **Meeting Description**: String (max 100 characters displayed)
- **Meeting Time**: Timestamp
- **Attendees**: Array of user objects (show max 3)

---

## Assumptions

1. **Data Availability**: The system already has existing projects, tasks, and users to display (not building these entities)
2. **Authentication**: User authentication and session management is already implemented
3. **Persian Language**: System uses Persian language and Jalali calendar (based on existing codebase using jalali-moment)
4. **Real-time Infrastructure**: System has capability to push updates or poll for changes
5. **Responsive Design**: System already uses Tailwind CSS for consistent styling
6. **Icon Library**: Using Lucide icons as shown in reference design
7. **Data Refresh**: 30-second auto-refresh interval is acceptable to users and server load
8. **Task Visibility**: Users should only see their own assigned tasks (not team-wide tasks)
9. **Mobile Support**: Dashboard should work on mobile devices but may show simplified layout
10. **Browser Support**: Modern browsers (Chrome, Firefox, Safari, Edge) with ES6+ JavaScript support

---

## Constraints

### Technical Constraints
- Must integrate with existing Go backend API
- Must use Svelte 5 with runes syntax for frontend components
- Must use existing PostgreSQL database schema (projects, tasks, users tables)
- Must follow existing authentication/authorization patterns (JWT tokens)
- Must maintain consistency with existing Tailwind CSS design system

### Design Constraints
- Dashboard must match visual style of reference HTML design (dashboard.html)
- Must use existing color scheme: indigo primary, slate neutrals
- Card-based layout with rounded corners and subtle shadows
- Sidebar navigation must remain accessible on all pages

### Business Constraints
- Dashboard must load quickly enough for daily use (target <2 seconds)
- Must not introduce additional infrastructure costs (use existing polling/refresh patterns)
- Must work with existing team permissions and project visibility rules

---

## Dependencies

### Internal Dependencies
- Existing backend API endpoints for projects, tasks, users
- Existing authentication system (JWT tokens in httpOnly cookies)
- Existing PostgreSQL database with project_management schema
- Existing Svelte 5 frontend application structure
- Existing stores for state management (projects store, tasks store, auth store)

### External Dependencies
- Lucide icons library (for UI icons)
- Tailwind CSS (for styling)
- jalali-moment (for Persian date formatting)

---

## Out of Scope

The following items are explicitly NOT part of this feature:

1. **Customizable Dashboard**: Users cannot customize which widgets appear or their layout
2. **Dashboard Widgets**: No drag-and-drop widget system
3. **Advanced Filtering**: No filtering options on dashboard cards (use full list pages instead)
4. **Charts and Graphs**: No data visualization beyond progress bars
5. **Activity Feed**: No timeline of recent activities or changes
6. **Notifications Panel**: Notification bell shows count but clicking doesn't show details
7. **Calendar View**: No calendar widget for deadlines and meetings
8. **Time Tracking Widget**: No time logging interface on dashboard
9. **Project Creation**: "New Project" button exists but form is out of scope for this feature
10. **Task Creation**: "Add New Task" button exists but form is out of scope for this feature
11. **Inline Editing**: No ability to edit projects or tasks directly from dashboard cards
12. **Keyboard Shortcuts**: No keyboard navigation for dashboard
13. **Dark Mode**: Dashboard uses light theme only
14. **Export**: No ability to export dashboard data
15. **Print View**: No optimized print layout

---

## Related Features

- **Feature 001**: Enhanced entities with Redmine fields (provides project and task data models)
- **Feature 002**: Svelte 5 upgrade (establishes frontend framework and patterns)
- **Feature 003**: User authentication (provides login and session management)
- **Feature 004**: Advanced task search (provides task filtering capabilities, not used on dashboard)

---

## Notes

### Design Reference
The visual design follows the attached `dashboard.html` static mockup, which demonstrates:
- Clean card-based layout with subtle shadows
- Indigo/purple gradient accent cards
- Responsive grid system
- Consistent spacing and typography
- Professional, modern aesthetic

### Persian Language Considerations
- All dates must be displayed in Jalali calendar format using jalali-moment
- Text direction should be RTL-aware but current design uses LTR layout
- Number formatting should follow Persian conventions (e.g., ۱۲۳ vs 123) [NEEDS CLARIFICATION: Number format preference]

### Performance Considerations
- Dashboard should lazy-load avatars and images
- Consider pagination if task/project lists exceed display limits
- Auto-refresh should be debounced to prevent excessive API calls
- Statistics should be calculated efficiently (consider caching)

### Future Enhancements
Potential additions for future iterations (not in current scope):
- Dashboard customization options
- Data visualization (charts, graphs)
- Activity timeline
- Quick actions (inline task creation, status updates)
- Dashboard themes (dark mode)
