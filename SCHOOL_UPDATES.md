# Admin Dashboard School Updates

## Pages to Remove/Repurpose

### 1. Orders Page (`/admin/orders`)
- **Current**: Manages e-commerce orders with payment tracking
- **Action**: Remove entirely (no purchases in school platform)
- **Files**: 
  - `frontend/src/pages/admin/Orders.jsx`
  - `backend/internal/handlers/order_handler.go` (keep for now, might repurpose)
  - `backend/internal/services/order_service.go`

### 2. Shipping Page (`/admin/shipping`)
- **Current**: Manages physical book shipping
- **Action**: Remove entirely (digital-only platform)
- **Files**:
  - `frontend/src/pages/admin/Shipping.jsx`
  - `backend/internal/handlers/shipping_handler.go`
  - `backend/internal/services/shipping_service.go`

## Pages to Update for School Context

### 3. Dashboard Overview (`/admin`)
- **Current**: Shows sales, revenue, orders
- **Update to show**:
  - Total students, teachers, school admins
  - Total books in library
  - Active reading sessions
  - Books added to libraries this week/month
  - Most popular books
  - Reading completion rates
- **Files**: `frontend/src/pages/admin/Dashboard.jsx`

### 4. Users Page (`/admin/users`)
- **Current**: Generic user management
- **Updates needed**:
  - Add school_name, class_level columns to table
  - Filter by role (student, teacher, school_admin, platform_admin)
  - Filter by school (for platform admins managing multiple schools)
  - Show student count per class/grade
  - Bulk actions: assign books to students, create class groups
- **Status**: ✅ Role display already fixed

### 5. Books Page (`/admin/books`)
- **Current**: Book catalog management
- **Updates needed**:
  - Remove price/purchase fields
  - Add "Assign to Students" bulk action
  - Add "Recommended Grade Level" field
  - Add "Subject/Category" more prominently
  - Show "Times Added to Library" instead of "Times Purchased"
- **Files**: `frontend/src/pages/admin/Books.jsx`

### 6. Library Management (`/admin/library`)
- **Current**: Might not exist or is generic
- **Updates needed**:
  - View all books in all student libraries
  - See which students have which books
  - Bulk assign/remove books from student libraries
- **Files**: `frontend/src/pages/admin/Library.jsx`

### 7. Reviews Page (`/admin/reviews`)
- **Current**: Book reviews moderation
- **Updates needed**:
  - Show student name and class level
  - Filter by grade/class
- **Status**: Probably fine as-is, minor updates

### 8. Reading Analytics (`/admin/reading`)
- **Current**: Reading session analytics
- **Updates needed**:
  - Group by class/grade
  - Show average reading time per student
  - Identify struggling readers (low completion rates)
  - Most/least read books by grade level
  - Reading streaks and achievements
- **Files**: `frontend/src/pages/admin/Reading.jsx`

### 9. Reports Page (`/admin/reports`)
- **Current**: Sales and revenue reports
- **Update to show**:
  - Student reading reports (by class, grade, individual)
  - Book popularity reports
  - Reading completion reports
  - Library usage reports
- **Files**: `frontend/src/pages/admin/Reports.jsx`

### 10. Email Templates (`/admin/email-templates`)
- **Status**: ✅ Already updated (removed order confirmation)
- **Additional updates**: Add templates for:
  - Book assignment notifications
  - Reading goal reminders
  - Achievement unlocked emails
  - Weekly reading summary

### 11. Blog Page (`/admin/blog`)
- **Status**: ✅ Already updated for school focus
- **Keep as-is**: For educational content, reading tips, announcements

### 12. About Page (`/admin/about`)
- **Status**: ✅ Already updated for school focus
- **Keep as-is**: Platform information

### 13. Contact Page (`/admin/contact`)
- **Status**: ✅ Already updated for school focus
- **Keep as-is**: Support messages from schools

### 14. FAQ Page (`/admin/faq`)
- **Status**: ✅ Already updated for school focus
- **Keep as-is**: School-related FAQs

### 15. Settings Page (`/admin/settings`)
- **Current**: General platform settings
- **Updates needed**:
  - Remove payment gateway settings
  - Add school-specific settings:
    - Default reading goals
    - Grade levels configuration
    - School year settings
    - Student registration settings (require approval, etc.)
- **Files**: `frontend/src/pages/admin/Settings.jsx`

## New Pages to Add

### 16. Classes/Groups Management (`/admin/classes`)
- **Purpose**: Create and manage class groups
- **Features**:
  - Create classes (e.g., "Grade 5A", "Year 10 Science")
  - Assign teachers to classes
  - Assign students to classes
  - Bulk assign books to entire class
  - View class reading statistics

### 17. Assignments Page (`/admin/assignments`)
- **Purpose**: Manage reading assignments
- **Features**:
  - Create reading assignments (book + deadline)
  - Assign to individual students or entire classes
  - Track completion status
  - Set reading goals (pages per day, finish by date)
  - Send reminders

## Navigation Updates

### Admin Sidebar Menu
**Remove:**
- Orders
- Shipping

**Keep:**
- Overview (Dashboard)
- Users
- Roles
- Audit Log
- Books
- Library Management
- Reviews
- Reading Analytics
- Reports
- Email Templates
- Blog
- About
- Contact
- FAQ
- Settings

**Add:**
- Classes/Groups
- Assignments

## Priority Order

1. **High Priority** (Remove e-commerce):
   - Remove Orders page
   - Remove Shipping page
   - Update Dashboard overview metrics

2. **Medium Priority** (School-specific features):
   - Update Books page (remove prices, add grade levels)
   - Update Library Management
   - Update Reading Analytics grouping
   - Update Reports content

3. **Low Priority** (Nice to have):
   - Add Classes/Groups management
   - Add Assignments page
   - Enhanced filtering and bulk actions

## Database Changes Needed

- Remove: `orders`, `order_items`, `shipping_addresses`, `payments` tables (or keep for historical data)
- Add: `classes` table (optional - can use tags/metadata on users)
- Add: `assignments` table (book_id, assigned_to, assigned_by, due_date, status)
- Update: `books` table - add `recommended_grade_level`, remove `price`
