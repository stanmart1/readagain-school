-- PostgreSQL Database Schema for ReadAgain
-- This file creates all necessary tables and initial data

-- Create sequences
CREATE SEQUENCE IF NOT EXISTS roles_id_seq;
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS permissions_id_seq;
CREATE SEQUENCE IF NOT EXISTS books_id_seq;
CREATE SEQUENCE IF NOT EXISTS authors_id_seq;
CREATE SEQUENCE IF NOT EXISTS categories_id_seq;
CREATE SEQUENCE IF NOT EXISTS reviews_id_seq;
CREATE SEQUENCE IF NOT EXISTS blogs_id_seq;
CREATE SEQUENCE IF NOT EXISTS about_pages_id_seq;
CREATE SEQUENCE IF NOT EXISTS user_libraries_id_seq;
CREATE SEQUENCE IF NOT EXISTS reading_sessions_id_seq;
CREATE SEQUENCE IF NOT EXISTS faqs_id_seq;
CREATE SEQUENCE IF NOT EXISTS reading_goals_id_seq;
CREATE SEQUENCE IF NOT EXISTS system_settings_id_seq;
CREATE SEQUENCE IF NOT EXISTS audit_logs_id_seq;
CREATE SEQUENCE IF NOT EXISTS auth_logs_id_seq;
CREATE SEQUENCE IF NOT EXISTS token_blacklists_id_seq;
CREATE SEQUENCE IF NOT EXISTS notifications_id_seq;
CREATE SEQUENCE IF NOT EXISTS user_achievements_id_seq;
CREATE SEQUENCE IF NOT EXISTS achievements_id_seq;
CREATE SEQUENCE IF NOT EXISTS activities_id_seq;
CREATE SEQUENCE IF NOT EXISTS contact_messages_id_seq;
CREATE SEQUENCE IF NOT EXISTS groups_id_seq;
CREATE SEQUENCE IF NOT EXISTS group_members_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_rooms_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_messages_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_members_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_reactions_id_seq;
CREATE SEQUENCE IF NOT EXISTS highlights_id_seq;
CREATE SEQUENCE IF NOT EXISTS notes_id_seq;
CREATE SEQUENCE IF NOT EXISTS licenses_id_seq;

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
  id bigint PRIMARY KEY DEFAULT nextval('roles_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  name text NOT NULL UNIQUE,
  description text
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
  id bigint PRIMARY KEY DEFAULT nextval('users_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  email text NOT NULL UNIQUE,
  username text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  first_name text,
  last_name text,
  phone_number text,
  school_name text,
  school_category text,
  class_level text,
  department text,
  role_id bigint REFERENCES roles(id),
  is_active boolean DEFAULT true,
  is_email_verified boolean DEFAULT false,
  verification_token text,
  verification_token_expires timestamp with time zone,
  last_login timestamp with time zone
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
  id bigint PRIMARY KEY DEFAULT nextval('permissions_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  name text NOT NULL UNIQUE,
  description text,
  category text,
  resource text,
  action text
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
  role_id bigint REFERENCES roles(id) ON DELETE CASCADE,
  permission_id bigint REFERENCES permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);

-- Create books table
CREATE TABLE IF NOT EXISTS books (
  id bigint PRIMARY KEY DEFAULT nextval('books_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  author_id bigint,
  title text NOT NULL,
  subtitle text,
  description text,
  isbn text,
  publication_date date,
  page_count integer,
  language text DEFAULT 'en',
  file_path text,
  cover_image text,
  file_size bigint,
  status text DEFAULT 'draft',
  category_id bigint,
  featured boolean DEFAULT false,
  price decimal(10,2) DEFAULT 0.00
);

-- Create authors table
CREATE TABLE IF NOT EXISTS authors (
  id bigint PRIMARY KEY DEFAULT nextval('authors_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  user_id bigint,
  business_name text,
  bio text,
  website text,
  social_links jsonb
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
  id bigint PRIMARY KEY DEFAULT nextval('categories_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  name text NOT NULL UNIQUE,
  description text,
  status text DEFAULT 'active'
);

-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
  id bigint PRIMARY KEY DEFAULT nextval('reviews_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  user_id bigint REFERENCES users(id),
  book_id bigint REFERENCES books(id),
  rating integer CHECK (rating >= 1 AND rating <= 5),
  comment text,
  status text DEFAULT 'published'
);

-- Create blogs table
CREATE TABLE IF NOT EXISTS blogs (
  id bigint PRIMARY KEY DEFAULT nextval('blogs_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  title text NOT NULL,
  slug text UNIQUE,
  content text,
  excerpt text,
  featured_image text,
  author_id bigint REFERENCES users(id),
  status text DEFAULT 'draft',
  published_at timestamp with time zone
);

-- Create about_pages table
CREATE TABLE IF NOT EXISTS about_pages (
  id bigint PRIMARY KEY DEFAULT nextval('about_pages_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  title text,
  content text,
  mission text,
  vision text,
  values jsonb
);

-- Create user_libraries table
CREATE TABLE IF NOT EXISTS user_libraries (
  id bigint PRIMARY KEY DEFAULT nextval('user_libraries_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  user_id bigint REFERENCES users(id),
  book_id bigint REFERENCES books(id),
  status text DEFAULT 'not_started',
  progress decimal(5,2) DEFAULT 0.00,
  current_page integer DEFAULT 1,
  started_at timestamp with time zone,
  completed_at timestamp with time zone,
  last_read_at timestamp with time zone
);

-- Create other necessary tables
CREATE TABLE IF NOT EXISTS reading_sessions (
  id bigint PRIMARY KEY DEFAULT nextval('reading_sessions_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  user_id bigint REFERENCES users(id),
  book_id bigint REFERENCES books(id),
  start_time timestamp with time zone,
  end_time timestamp with time zone,
  pages_read integer DEFAULT 0,
  session_duration integer DEFAULT 0
);

CREATE TABLE IF NOT EXISTS faqs (
  id bigint PRIMARY KEY DEFAULT nextval('faqs_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  question text NOT NULL,
  answer text NOT NULL,
  category text,
  order_index integer DEFAULT 0,
  is_published boolean DEFAULT true
);

CREATE TABLE IF NOT EXISTS system_settings (
  id bigint PRIMARY KEY DEFAULT nextval('system_settings_id_seq'),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp with time zone,
  key text NOT NULL UNIQUE,
  value text,
  description text,
  category text
);

-- Insert default data
INSERT INTO roles (name, description) VALUES 
  ('platform_admin', 'Platform administrator with full system access'),
  ('school_admin', 'School administrator who manages school users and library'),
  ('teacher', 'Teacher who assigns books and tracks student progress'),
  ('student', 'Student who reads books')
ON CONFLICT (name) DO NOTHING;

-- Insert admin user
INSERT INTO users (email, username, password_hash, first_name, last_name, phone_number, school_name, school_category, department, role_id, is_active, is_email_verified) 
SELECT 'admin@readagain.com', 'admin', '$2a$10$fOoVIiYwORELVQctuZtMiuvKH7D2KPUVWjf7v3ZyKUotCmrygaUxa', 'Super', 'Admin', '0896655577', 'University Of Lagos', 'Tertiary', 'Computer Science', r.id, true, true
FROM roles r WHERE r.name = 'platform_admin'
ON CONFLICT (email) DO NOTHING;

-- Insert default categories
INSERT INTO categories (name, description) VALUES 
  ('Fiction', 'Fiction books and novels'),
  ('Non-Fiction', 'Non-fiction books'),
  ('Science', 'Science and technology'),
  ('History', 'Historical books'),
  ('Biography', 'Biographies and memoirs'),
  ('Education', 'Educational materials')
ON CONFLICT (name) DO NOTHING;


-- Table: role_permissions
DROP TABLE IF EXISTS "role_permissions";
CREATE TABLE "role_permissions" (
  "role_id" bigint NOT NULL,
  "permission_id" bigint NOT NULL
);

-- Data for table: role_permissions
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '1');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '2');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '3');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '4');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '5');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '6');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '7');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '8');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '9');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '10');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '11');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '12');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '13');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '14');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '15');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '16');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '17');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '18');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '19');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '20');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '21');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '22');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '23');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '24');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '25');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '26');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '27');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '28');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '29');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '30');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '31');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '32');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '33');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '34');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '35');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '36');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '37');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '38');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '39');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '40');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '41');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '42');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '43');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '44');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '45');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '46');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '47');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '48');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '49');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '50');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '51');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '52');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '53');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '54');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '55');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '56');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '57');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '58');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('8', '59');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '3');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '4');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '15');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '16');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '30');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '20');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '21');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '32');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '39');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '48');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '50');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('7', '53');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('6', '15');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('6', '30');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('6', '20');
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES ('5', '15');


-- Table: permissions
DROP TABLE IF EXISTS "permissions";
CREATE TABLE "permissions" (
  "id" bigint NOT NULL DEFAULT nextval('permissions_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "name" text NOT NULL,
  "description" text,
  "category" text
);

-- Data for table: permissions
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('1', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'analytics.view', 'View analytics dashboard', 'analytics');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('2', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'analytics.export', 'Export analytics data', 'analytics');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('3', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'users.view', 'View users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('4', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'users.create', 'Create users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('5', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'users.edit', 'Edit users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('6', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'users.delete', 'Delete users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('7', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'users.manage', 'Full user management', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('8', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'roles.view', 'View roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('9', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'roles.create', 'Create roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('10', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'roles.edit', 'Edit roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('11', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'roles.delete', 'Delete roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('12', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'roles.manage', 'Full role management', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('13', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'audit_logs.view', 'View audit logs', 'audit');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('14', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'audit_logs.export', 'Export audit logs', 'audit');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('15', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'books.view', 'View books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('16', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'books.create', 'Create books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('17', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'books.edit', 'Edit books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('18', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'books.delete', 'Delete books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('19', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'books.manage', 'Full book management', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('20', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reviews.view', 'View reviews', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('21', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reviews.moderate', 'Moderate reviews', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('22', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reviews.delete', 'Delete reviews', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('23', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reviews.manage', 'Full review management', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('24', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'orders.view', 'View orders', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('25', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'orders.edit', 'Edit orders', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('26', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'orders.delete', 'Delete orders', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('27', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'orders.manage', 'Full order management', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('28', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'shipping.view', 'View shipping', 'shipping');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('29', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'shipping.manage', 'Manage shipping', 'shipping');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('30', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reading.view_analytics', 'View reading analytics', 'reading');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('31', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reading.manage', 'Manage reading data', 'reading');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('32', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reports.view', 'View reports', 'reports');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('33', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reports.generate', 'Generate reports', 'reports');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('34', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'reports.export', 'Export reports', 'reports');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('35', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'email_templates.view', 'View email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('36', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'email_templates.create', 'Create email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('37', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'email_templates.edit', 'Edit email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('38', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'email_templates.delete', 'Delete email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('39', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'blog.view', 'View blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('40', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'blog.create', 'Create blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('41', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'blog.edit', 'Edit blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('42', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'blog.delete', 'Delete blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('43', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'blog.publish', 'Publish blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('44', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'works.view', 'View works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('45', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'works.create', 'Create works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('46', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'works.edit', 'Edit works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('47', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'works.delete', 'Delete works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('48', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'about.view', 'View about page', 'about');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('49', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'about.edit', 'Edit about page', 'about');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('50', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'contact.view', 'View contact messages', 'contact');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('51', Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:57 GMT+0100 (West Africa Standard Time), NULL, 'contact.reply', 'Reply to contact messages', 'contact');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('52', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'contact.delete', 'Delete contact messages', 'contact');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('53', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'faq.view', 'View FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('54', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'faq.create', 'Create FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('55', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'faq.edit', 'Edit FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('56', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'faq.delete', 'Delete FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('57', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'settings.view', 'View settings', 'settings');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('58', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'settings.edit', 'Edit settings', 'settings');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('59', Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:35:58 GMT+0100 (West Africa Standard Time), NULL, 'settings.manage', 'Full settings management', 'settings');


-- Table: auth_logs
DROP TABLE IF EXISTS "auth_logs";
CREATE TABLE "auth_logs" (
  "id" bigint NOT NULL DEFAULT nextval('auth_logs_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "action" text NOT NULL,
  "ip_address" text,
  "user_agent" text,
  "success" boolean
);

-- Data for table: auth_logs
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('3', Thu Dec 11 2025 11:36:31 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 11:36:31 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('4', Thu Dec 11 2025 11:37:36 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 11:37:36 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('5', Thu Dec 11 2025 12:36:42 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 12:36:42 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('6', Thu Dec 11 2025 13:27:44 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 13:27:44 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('7', Thu Dec 11 2025 13:34:25 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 13:34:25 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('8', Thu Dec 11 2025 13:36:12 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 13:36:12 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('9', Thu Dec 11 2025 13:36:25 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 13:36:25 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('10', Thu Dec 11 2025 16:40:47 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 16:40:47 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('11', Thu Dec 11 2025 17:44:14 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 17:44:14 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('12', Thu Dec 11 2025 17:44:27 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 17:44:27 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('13', Thu Dec 11 2025 17:56:50 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 17:56:50 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('14', Thu Dec 11 2025 18:53:02 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 18:53:02 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('15', Thu Dec 11 2025 22:29:02 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 22:29:02 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('16', Thu Dec 11 2025 22:39:33 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 22:39:33 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('17', Thu Dec 11 2025 22:41:44 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 22:41:44 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('18', Thu Dec 11 2025 22:42:02 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 22:42:02 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('19', Thu Dec 11 2025 23:09:21 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 23:09:21 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('20', Thu Dec 11 2025 23:42:19 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 23:42:19 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('21', Mon Dec 22 2025 10:46:47 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 10:46:47 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('22', Mon Dec 22 2025 13:03:35 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 13:03:35 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('23', Mon Dec 22 2025 13:09:40 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 13:09:40 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('24', Mon Dec 22 2025 13:16:08 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 13:16:08 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('25', Mon Dec 22 2025 13:19:56 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 13:19:56 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('26', Mon Dec 22 2025 22:25:24 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 22:25:24 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('27', Thu Dec 25 2025 17:21:29 GMT+0100 (West Africa Standard Time), Thu Dec 25 2025 17:21:29 GMT+0100 (West Africa Standard Time), NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);


-- Table: token_blacklists
DROP TABLE IF EXISTS "token_blacklists";
CREATE TABLE "token_blacklists" (
  "id" bigint NOT NULL DEFAULT nextval('token_blacklists_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "token" text NOT NULL,
  "expires_at" timestamp with time zone NOT NULL
);


-- Table: authors
DROP TABLE IF EXISTS "authors";
CREATE TABLE "authors" (
  "id" bigint NOT NULL DEFAULT nextval('authors_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "business_name" text NOT NULL,
  "bio" text,
  "photo" text,
  "website" text,
  "email" text,
  "status" text DEFAULT 'active'::text
);

-- Data for table: authors
INSERT INTO "authors" ("id", "created_at", "updated_at", "deleted_at", "user_id", "business_name", "bio", "photo", "website", "email", "status") VALUES ('1', Mon Dec 22 2025 15:35:15 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 10:29:40 GMT+0100 (West Africa Standard Time), NULL, '2', 'Classic Literature Press', 'Publisher of timeless classic literature for young readers', '', '', 'author@gmail.com', 'active');


-- Table: books
DROP TABLE IF EXISTS "books";
CREATE TABLE "books" (
  "id" bigint NOT NULL DEFAULT nextval('books_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "author_id" bigint NOT NULL,
  "title" text NOT NULL,
  "subtitle" text,
  "description" text,
  "short_description" text,
  "cover_image" text,
  "file_path" text,
  "file_size" bigint,
  "category_id" bigint,
  "isbn" text,
  "is_featured" boolean DEFAULT false,
  "is_bestseller" boolean DEFAULT false,
  "is_new_release" boolean DEFAULT false,
  "status" text DEFAULT 'published'::text,
  "is_active" boolean DEFAULT true,
  "pages" bigint,
  "language" text DEFAULT 'English'::text,
  "publisher" text,
  "publication_date" timestamp with time zone,
  "library_count" bigint DEFAULT 0,
  "view_count" bigint DEFAULT 0,
  "seo_title" text,
  "seo_description" text,
  "seo_keywords" text
);

-- Data for table: books
INSERT INTO "books" ("id", "created_at", "updated_at", "deleted_at", "author_id", "title", "subtitle", "description", "short_description", "cover_image", "file_path", "file_size", "category_id", "isbn", "is_featured", "is_bestseller", "is_new_release", "status", "is_active", "pages", "language", "publisher", "publication_date", "library_count", "view_count", "seo_title", "seo_description", "seo_keywords") VALUES ('11', Mon Dec 22 2025 20:52:44 GMT+0100 (West Africa Standard Time), Thu Dec 25 2025 17:29:15 GMT+0100 (West Africa Standard Time), NULL, '1', 'Moby Dick', '', '', '', 'covers/0683010e-42a7-4af3-95c3-4e01b00e19f8.jpg', 'books/e068cce3-0e5c-4c71-bd04-cd0666b87723.epub', '816008', '2', '', true, false, false, 'published', true, '0', 'English', '', Mon Dec 22 2025 20:52:44 GMT+0100 (West Africa Standard Time), '1', '15', '', '', '');


-- Table: categories
DROP TABLE IF EXISTS "categories";
CREATE TABLE "categories" (
  "id" bigint NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "name" text NOT NULL,
  "description" text,
  "status" text DEFAULT 'active'::text
);

-- Data for table: categories
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('1', Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), NULL, 'Fiction', 'Fictional stories and novels', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('2', Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), NULL, 'Adventure', 'Adventure and exploration stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('3', Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), NULL, 'Fantasy', 'Fantasy and magical stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('4', Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), NULL, 'Mystery', 'Mystery and detective stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('5', Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 15:36:19 GMT+0100 (West Africa Standard Time), NULL, 'Science Fiction', 'Science fiction stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('6', Mon Dec 22 2025 15:36:20 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 15:36:20 GMT+0100 (West Africa Standard Time), NULL, 'Classic', 'Classic literature', 'active');


-- Table: user_libraries
DROP TABLE IF EXISTS "user_libraries";
CREATE TABLE "user_libraries" (
  "id" bigint NOT NULL DEFAULT nextval('user_libraries_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "book_id" bigint NOT NULL,
  "progress" numeric DEFAULT 0,
  "current_page" bigint DEFAULT 0,
  "last_read_at" timestamp with time zone,
  "completed_at" timestamp with time zone,
  "is_favorite" boolean DEFAULT false,
  "rating" bigint
);

-- Data for table: user_libraries
INSERT INTO "user_libraries" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "progress", "current_page", "last_read_at", "completed_at", "is_favorite", "rating") VALUES ('2', Mon Dec 22 2025 20:54:25 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 13:07:24 GMT+0100 (West Africa Standard Time), NULL, '2', '11', '35.36493604213695', '0', NULL, NULL, false, '0');


-- Table: reading_sessions
DROP TABLE IF EXISTS "reading_sessions";
CREATE TABLE "reading_sessions" (
  "id" bigint NOT NULL DEFAULT nextval('reading_sessions_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "book_id" bigint NOT NULL,
  "start_time" timestamp with time zone NOT NULL,
  "end_time" timestamp with time zone,
  "duration" bigint,
  "pages_read" bigint,
  "start_page" bigint,
  "end_page" bigint
);

-- Data for table: reading_sessions
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('1', Tue Dec 23 2025 13:06:03 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 13:06:03 GMT+0100 (West Africa Standard Time), NULL, '2', '11', Mon Jan 01 0001 00:13:35 GMT+0013 (West Africa Standard Time), NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('2', Tue Dec 23 2025 13:06:03 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 13:06:03 GMT+0100 (West Africa Standard Time), NULL, '2', '11', Mon Jan 01 0001 00:13:35 GMT+0013 (West Africa Standard Time), NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('3', Tue Dec 23 2025 13:06:51 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 13:06:51 GMT+0100 (West Africa Standard Time), NULL, '2', '11', Mon Jan 01 0001 00:13:35 GMT+0013 (West Africa Standard Time), NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('4', Tue Dec 23 2025 13:06:52 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 13:07:27 GMT+0100 (West Africa Standard Time), NULL, '2', '11', Mon Jan 01 0001 00:13:35 GMT+0013 (West Africa Standard Time), NULL, '35', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('5', Thu Dec 25 2025 17:21:39 GMT+0100 (West Africa Standard Time), Thu Dec 25 2025 17:21:39 GMT+0100 (West Africa Standard Time), NULL, '2', '11', Mon Jan 01 0001 00:13:35 GMT+0013 (West Africa Standard Time), NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('6', Thu Dec 25 2025 17:21:39 GMT+0100 (West Africa Standard Time), Thu Dec 25 2025 17:21:39 GMT+0100 (West Africa Standard Time), NULL, '2', '11', Mon Jan 01 0001 00:13:35 GMT+0013 (West Africa Standard Time), NULL, '0', '0', '0', '0');


-- Table: blogs
DROP TABLE IF EXISTS "blogs";
CREATE TABLE "blogs" (
  "id" bigint NOT NULL DEFAULT nextval('blogs_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "title" text NOT NULL,
  "slug" text NOT NULL,
  "excerpt" text,
  "content" text,
  "featured_image" text,
  "author_id" bigint,
  "category_id" bigint,
  "status" text DEFAULT 'draft'::text,
  "is_featured" boolean DEFAULT false,
  "views" bigint DEFAULT 0,
  "read_time" bigint,
  "tags" text,
  "published_at" timestamp with time zone,
  "seo_title" text,
  "seo_description" text,
  "seo_keywords" text
);


-- Table: faqs
DROP TABLE IF EXISTS "faqs";
CREATE TABLE "faqs" (
  "id" bigint NOT NULL DEFAULT nextval('faqs_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "question" text NOT NULL,
  "answer" text NOT NULL,
  "category" text,
  "order" bigint DEFAULT 0,
  "is_active" boolean DEFAULT true
);


-- Table: reading_goals
DROP TABLE IF EXISTS "reading_goals";
CREATE TABLE "reading_goals" (
  "id" bigint NOT NULL DEFAULT nextval('reading_goals_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "type" text NOT NULL,
  "target" bigint NOT NULL,
  "current" bigint DEFAULT 0,
  "start_date" timestamp with time zone NOT NULL,
  "end_date" timestamp with time zone NOT NULL,
  "is_completed" boolean DEFAULT false,
  "book_id" integer
);

-- Data for table: reading_goals
INSERT INTO "reading_goals" ("id", "created_at", "updated_at", "deleted_at", "user_id", "type", "target", "current", "start_date", "end_date", "is_completed", "book_id") VALUES ('1', Mon Dec 22 2025 23:45:23 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 23:45:23 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 23:46:06 GMT+0100 (West Africa Standard Time), '2', 'books', '12', '0', Mon Dec 22 2025 01:00:00 GMT+0100 (West Africa Standard Time), Fri Dec 26 2025 01:00:00 GMT+0100 (West Africa Standard Time), false, NULL);


-- Table: system_settings
DROP TABLE IF EXISTS "system_settings";
CREATE TABLE "system_settings" (
  "id" bigint NOT NULL DEFAULT nextval('system_settings_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "key" text NOT NULL,
  "value" text,
  "data_type" text,
  "category" text,
  "description" text,
  "is_public" boolean DEFAULT false
);


-- Table: audit_logs
DROP TABLE IF EXISTS "audit_logs";
CREATE TABLE "audit_logs" (
  "id" bigint NOT NULL DEFAULT nextval('audit_logs_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint,
  "action" text NOT NULL,
  "entity_type" text,
  "entity_id" bigint,
  "old_value" text,
  "new_value" text,
  "ip_address" text,
  "user_agent" text
);

-- Data for table: audit_logs
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('1', Mon Dec 22 2025 17:22:34 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 17:22:34 GMT+0100 (West Africa Standard Time), NULL, '2', 'update_review_status', 'review', '50', '', 'approved', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('2', Mon Dec 22 2025 17:22:40 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 17:22:40 GMT+0100 (West Africa Standard Time), NULL, '2', 'toggle_review_featured', 'review', '50', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('3', Mon Dec 22 2025 18:19:54 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 18:19:54 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '55', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('4', Mon Dec 22 2025 20:57:11 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 20:57:11 GMT+0100 (West Africa Standard Time), NULL, '2', 'update_review_status', 'review', '56', '', 'approved', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('5', Mon Dec 22 2025 20:57:16 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 20:57:16 GMT+0100 (West Africa Standard Time), NULL, '2', 'toggle_review_featured', 'review', '56', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('6', Mon Dec 22 2025 21:24:32 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:24:32 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '44', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('7', Mon Dec 22 2025 21:24:35 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:24:35 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '54', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('8', Mon Dec 22 2025 21:24:44 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:24:44 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '42', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('9', Mon Dec 22 2025 21:24:48 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:24:48 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '53', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('10', Mon Dec 22 2025 21:24:59 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:24:59 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '52', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('11', Mon Dec 22 2025 21:25:10 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:10 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '51', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('12', Mon Dec 22 2025 21:25:14 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:14 GMT+0100 (West Africa Standard Time), NULL, '2', 'toggle_review_featured', 'review', '50', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('13', Mon Dec 22 2025 21:25:17 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:17 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '50', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('14', Mon Dec 22 2025 21:25:21 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:21 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '49', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('15', Mon Dec 22 2025 21:25:24 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:24 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '48', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('16', Mon Dec 22 2025 21:25:27 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:27 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '47', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('17', Mon Dec 22 2025 21:25:30 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:30 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '46', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('18', Mon Dec 22 2025 21:25:33 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:33 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '45', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('19', Mon Dec 22 2025 21:25:35 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:35 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '43', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('20', Mon Dec 22 2025 21:25:38 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:38 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '41', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('21', Mon Dec 22 2025 21:25:40 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:40 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '40', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('22', Mon Dec 22 2025 21:25:43 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:43 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '39', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('23', Mon Dec 22 2025 21:25:45 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:45 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '38', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('24', Mon Dec 22 2025 21:25:49 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 21:25:49 GMT+0100 (West Africa Standard Time), NULL, '2', 'delete_review', 'review', '37', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');


-- Table: notifications
DROP TABLE IF EXISTS "notifications";
CREATE TABLE "notifications" (
  "id" bigint NOT NULL DEFAULT nextval('notifications_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "type" text NOT NULL,
  "title" text NOT NULL,
  "message" text,
  "is_read" boolean DEFAULT false,
  "action_url" text
);


-- Table: reviews
DROP TABLE IF EXISTS "reviews";
CREATE TABLE "reviews" (
  "id" bigint NOT NULL DEFAULT nextval('reviews_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "book_id" bigint NOT NULL,
  "rating" bigint NOT NULL,
  "comment" text,
  "status" text DEFAULT 'pending'::text,
  "is_featured" boolean DEFAULT false
);

-- Data for table: reviews
INSERT INTO "reviews" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "rating", "comment", "status", "is_featured") VALUES ('56', Mon Dec 22 2025 20:57:02 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 20:57:16 GMT+0100 (West Africa Standard Time), NULL, '2', '11', '2', 'At least image upload is working', 'approved', true);


-- Table: about_pages
DROP TABLE IF EXISTS "about_pages";
CREATE TABLE "about_pages" (
  "id" bigint NOT NULL DEFAULT nextval('about_pages_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "title" text,
  "content" text,
  "mission" text,
  "vision" text,
  "team_section" text,
  "values" text
);

-- Data for table: about_pages
INSERT INTO "about_pages" ("id", "created_at", "updated_at", "deleted_at", "title", "content", "mission", "vision", "team_section", "values") VALUES ('1', Thu Dec 11 2025 11:10:41 GMT+0100 (West Africa Standard Time), Thu Dec 11 2025 11:10:41 GMT+0100 (West Africa Standard Time), NULL, 'About ReadAgain', 'Welcome to ReadAgain - Your digital reading platform.', '', '', '', '');


-- Table: user_achievements
DROP TABLE IF EXISTS "user_achievements";
CREATE TABLE "user_achievements" (
  "id" bigint NOT NULL DEFAULT nextval('user_achievements_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "achievement_id" bigint NOT NULL,
  "progress" bigint DEFAULT 0,
  "is_unlocked" boolean DEFAULT false,
  "unlocked_at" timestamp without time zone
);


-- Table: achievements
DROP TABLE IF EXISTS "achievements";
CREATE TABLE "achievements" (
  "id" bigint NOT NULL DEFAULT nextval('achievements_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "name" text NOT NULL,
  "description" text,
  "icon" text,
  "type" text,
  "target" bigint,
  "points" bigint,
  "category" character varying(255)
);

-- Data for table: achievements
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('1', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'First Purchase', 'Purchase your first book', '📚', 'books_purchased', '1', '10', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('2', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Book Collector', 'Purchase 10 books', '📖', 'books_purchased', '10', '50', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('3', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Library Master', 'Purchase 50 books', '🏛️', 'books_purchased', '50', '200', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('4', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'First Finish', 'Complete your first book', '🎉', 'books_completed', '1', '20', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('5', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Bookworm', 'Complete 5 books', '🐛', 'books_completed', '5', '100', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('6', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Avid Reader', 'Complete 25 books', '📕', 'books_completed', '25', '300', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('7', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Reading Legend', 'Complete 100 books', '👑', 'books_completed', '100', '1000', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('8', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Getting Started', 'Read for 60 minutes', '⏱️', 'reading_minutes', '60', '15', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('9', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Dedicated Reader', 'Read for 10 hours', '⏰', 'reading_minutes', '600', '100', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('10', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Time Master', 'Read for 100 hours', '🕐', 'reading_minutes', '6000', '500', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('11', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'First Session', 'Complete your first reading session', '🎯', 'reading_sessions', '1', '5', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('12', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Consistent Reader', 'Complete 50 reading sessions', '📅', 'reading_sessions', '50', '150', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('13', Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), Wed Dec 10 2025 12:03:52 GMT+0100 (West Africa Standard Time), NULL, 'Reading Habit', 'Complete 200 reading sessions', '🔥', 'reading_sessions', '200', '400', NULL);


-- Table: activities
DROP TABLE IF EXISTS "activities";
CREATE TABLE "activities" (
  "id" bigint NOT NULL DEFAULT nextval('activities_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "type" text NOT NULL,
  "title" text,
  "description" text,
  "entity_type" text,
  "entity_id" bigint,
  "metadata" text
);


-- Table: contact_messages
DROP TABLE IF EXISTS "contact_messages";
CREATE TABLE "contact_messages" (
  "id" bigint NOT NULL DEFAULT nextval('contact_messages_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "subject" text,
  "message" text NOT NULL,
  "status" text DEFAULT 'new'::text,
  "reply" text
);


-- Table: groups
DROP TABLE IF EXISTS "groups";
CREATE TABLE "groups" (
  "id" bigint NOT NULL DEFAULT nextval('groups_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "name" text NOT NULL,
  "description" text,
  "created_by" bigint NOT NULL,
  "member_count" bigint DEFAULT 0
);

-- Data for table: groups
INSERT INTO "groups" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "created_by", "member_count") VALUES ('1', Mon Dec 22 2025 22:28:32 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 22:28:42 GMT+0100 (West Africa Standard Time), NULL, 'Science Club', '', '2', '1');


-- Table: group_members
DROP TABLE IF EXISTS "group_members";
CREATE TABLE "group_members" (
  "id" bigint NOT NULL DEFAULT nextval('group_members_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "group_id" bigint NOT NULL,
  "user_id" bigint NOT NULL
);

-- Data for table: group_members
INSERT INTO "group_members" ("id", "created_at", "updated_at", "deleted_at", "group_id", "user_id") VALUES ('1', Mon Dec 22 2025 22:28:42 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 22:28:42 GMT+0100 (West Africa Standard Time), NULL, '1', '2');


-- Table: chat_rooms
DROP TABLE IF EXISTS "chat_rooms";
CREATE TABLE "chat_rooms" (
  "id" bigint NOT NULL DEFAULT nextval('chat_rooms_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "type" text NOT NULL,
  "name" text NOT NULL,
  "description" text,
  "group_id" bigint,
  "book_id" bigint,
  "created_by" bigint NOT NULL,
  "is_active" boolean DEFAULT true,
  "last_message" text,
  "last_message_at" timestamp with time zone
);

-- Data for table: chat_rooms
INSERT INTO "chat_rooms" ("id", "created_at", "updated_at", "deleted_at", "type", "name", "description", "group_id", "book_id", "created_by", "is_active", "last_message", "last_message_at") VALUES ('1', Mon Dec 22 2025 23:30:28 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 23:30:47 GMT+0100 (West Africa Standard Time), NULL, 'group', 'Group Chat', 'Group discussion', '1', NULL, '2', true, 'hi', Mon Dec 22 2025 23:30:47 GMT+0100 (West Africa Standard Time));


-- Table: chat_messages
DROP TABLE IF EXISTS "chat_messages";
CREATE TABLE "chat_messages" (
  "id" bigint NOT NULL DEFAULT nextval('chat_messages_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "room_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "message" text NOT NULL,
  "message_type" text DEFAULT 'text'::text,
  "file_url" text,
  "file_name" text,
  "reply_to_id" bigint,
  "is_edited" boolean DEFAULT false,
  "edited_at" timestamp with time zone,
  "is_deleted" boolean DEFAULT false
);

-- Data for table: chat_messages
INSERT INTO "chat_messages" ("id", "created_at", "updated_at", "deleted_at", "room_id", "user_id", "message", "message_type", "file_url", "file_name", "reply_to_id", "is_edited", "edited_at", "is_deleted") VALUES ('1', Mon Dec 22 2025 23:30:47 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 23:30:47 GMT+0100 (West Africa Standard Time), NULL, '1', '2', 'hi', 'text', '', '', NULL, false, NULL, false);


-- Table: chat_members
DROP TABLE IF EXISTS "chat_members";
CREATE TABLE "chat_members" (
  "id" bigint NOT NULL DEFAULT nextval('chat_members_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "room_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "role" text DEFAULT 'member'::text,
  "last_read_at" timestamp with time zone,
  "unread_count" bigint DEFAULT 0,
  "is_muted" boolean DEFAULT false,
  "joined_at" timestamp with time zone NOT NULL
);

-- Data for table: chat_members
INSERT INTO "chat_members" ("id", "created_at", "updated_at", "deleted_at", "room_id", "user_id", "role", "last_read_at", "unread_count", "is_muted", "joined_at") VALUES ('1', Mon Dec 22 2025 23:30:28 GMT+0100 (West Africa Standard Time), Mon Dec 22 2025 23:44:50 GMT+0100 (West Africa Standard Time), NULL, '1', '2', 'admin', Mon Dec 22 2025 23:44:50 GMT+0100 (West Africa Standard Time), '0', false, Mon Dec 22 2025 23:30:28 GMT+0100 (West Africa Standard Time));


-- Table: chat_reactions
DROP TABLE IF EXISTS "chat_reactions";
CREATE TABLE "chat_reactions" (
  "id" bigint NOT NULL DEFAULT nextval('chat_reactions_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "message_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "emoji" text NOT NULL
);


-- Table: highlights
DROP TABLE IF EXISTS "highlights";
CREATE TABLE "highlights" (
  "id" bigint NOT NULL DEFAULT nextval('highlights_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "user_id" bigint NOT NULL,
  "book_id" bigint NOT NULL,
  "text" text NOT NULL,
  "color" text DEFAULT 'yellow'::text,
  "start_offset" bigint,
  "end_offset" bigint,
  "context" text,
  "cfi_range" text
);

-- Data for table: highlights
INSERT INTO "highlights" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "text", "color", "start_offset", "end_offset", "context", "cfi_range") VALUES ('1', Tue Dec 23 2025 10:50:54 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 10:50:54 GMT+0100 (West Africa Standard Time), NULL, '2', '11', 'Then God spake unto the fish; and from the shuddering cold and blackness of the sea, the whale came breeching up towards the warm and pleasant sun, and all the delights of air and earth; and ‘vomited out Jonah upon the dry land;’ when the word of the Lord came a second time; and Jonah, bruised and beaten—his ears, like two sea-shells, still multitudinously murmuring of the ocean—Jonah did the Almighty’s bidding. And what was that, shipmates? To preach the Truth to the face of Falsehood! That was it!', '#2196f3', '0', '504', 'Then God spake unto the fish; and from the shuddering cold and blackness of the sea, the whale came breeching up towards the warm and pleasant sun, and all the delights of air and earth; and ‘vomited out Jonah upon the dry land;’ when the word of the Lord came a second time; and Jonah, bruised and beaten—his ears, like two sea-shells, still multitudinously murmuring of the ocean—Jonah did the Almighty’s bidding. And what was that, shipmates? To preach the Truth to the face of Falsehood! That was it!', 'epubcfi(/6/6!/4/44,/5:967,/5:1519)');


-- Table: notes
DROP TABLE IF EXISTS "notes";
CREATE TABLE "notes" (
  "id" integer NOT NULL DEFAULT nextval('notes_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "content" text NOT NULL,
  "page_number" integer,
  "position" text,
  "color" character varying(50) DEFAULT 'yellow'::character varying,
  "user_id" bigint NOT NULL,
  "book_id" bigint NOT NULL,
  "page" bigint NOT NULL
);

-- Data for table: notes
INSERT INTO "notes" ("id", "created_at", "updated_at", "deleted_at", "content", "page_number", "position", "color", "user_id", "book_id", "page") VALUES (2, Tue Dec 23 2025 11:02:21 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 11:02:21 GMT+0100 (West Africa Standard Time), Tue Dec 23 2025 11:02:30 GMT+0100 (West Africa Standard Time), 'Test', NULL, NULL, 'yellow', '2', '11', '0');


-- Table: licenses
DROP TABLE IF EXISTS "licenses";
CREATE TABLE "licenses" (
  "id" bigint NOT NULL DEFAULT nextval('licenses_id_seq'::regclass),
  "api_key" text NOT NULL,
  "license_key" text NOT NULL,
  "device_fingerprint" text NOT NULL,
  "is_active" boolean DEFAULT false,
  "activated_at" timestamp with time zone,
  "last_verified_at" timestamp with time zone,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "is_valid" boolean DEFAULT true,
  "grace_period_started_at" timestamp with time zone,
  "grace_period_days" bigint DEFAULT 3,
  "is_locked" boolean DEFAULT false
);

-- Data for table: licenses
INSERT INTO "licenses" ("id", "api_key", "license_key", "device_fingerprint", "is_active", "activated_at", "last_verified_at", "created_at", "updated_at", "is_valid", "grace_period_started_at", "grace_period_days", "is_locked") VALUES ('1', 'pk_live_ey62swS2VTFWQnARkDbAcB1pg7kgousxqXAeP4JBaNo', 'YtCRzVM9LtisY19x_li2rHGydpsThSNdYH61ePuP27w', '4a091dc77282e609f2e15f368dcd3377c7bb9336c50bd087620e1d44335b47d2', true, Wed Dec 24 2025 13:48:53 GMT+0100 (West Africa Standard Time), Thu Dec 25 2025 17:20:33 GMT+0100 (West Africa Standard Time), Wed Dec 24 2025 13:48:53 GMT+0100 (West Africa Standard Time), Thu Dec 25 2025 17:20:33 GMT+0100 (West Africa Standard Time), false, Thu Dec 25 2025 17:20:33 GMT+0100 (West Africa Standard Time), '3', false);

