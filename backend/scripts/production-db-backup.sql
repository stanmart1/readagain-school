-- PostgreSQL Backup
-- Database: readagain

-- Create sequences first
CREATE SEQUENCE IF NOT EXISTS roles_id_seq;
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS permissions_id_seq;
CREATE SEQUENCE IF NOT EXISTS auth_logs_id_seq;
CREATE SEQUENCE IF NOT EXISTS token_blacklists_id_seq;
CREATE SEQUENCE IF NOT EXISTS authors_id_seq;
CREATE SEQUENCE IF NOT EXISTS books_id_seq;
CREATE SEQUENCE IF NOT EXISTS categories_id_seq;
CREATE SEQUENCE IF NOT EXISTS user_libraries_id_seq;
CREATE SEQUENCE IF NOT EXISTS reading_sessions_id_seq;
CREATE SEQUENCE IF NOT EXISTS blogs_id_seq;
CREATE SEQUENCE IF NOT EXISTS faqs_id_seq;
CREATE SEQUENCE IF NOT EXISTS reading_goals_id_seq;
CREATE SEQUENCE IF NOT EXISTS system_settings_id_seq;
CREATE SEQUENCE IF NOT EXISTS audit_logs_id_seq;
CREATE SEQUENCE IF NOT EXISTS notifications_id_seq;
CREATE SEQUENCE IF NOT EXISTS reviews_id_seq;
CREATE SEQUENCE IF NOT EXISTS about_pages_id_seq;
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


-- Table: roles
DROP TABLE IF EXISTS "roles";
CREATE TABLE "roles" (
  "id" bigint NOT NULL DEFAULT nextval('roles_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "name" text NOT NULL,
  "description" text
);

-- Data for table: roles
INSERT INTO "roles" ("id", "created_at", "updated_at", "deleted_at", "name", "description") VALUES ('8', 2025-12-22 11:39:02+01, 2025-12-22 13:09:02+01, NULL, 'platform_admin', 'Platform administrator with full system access');
INSERT INTO "roles" ("id", "created_at", "updated_at", "deleted_at", "name", "description") VALUES ('7', 2025-12-22 11:39:02+01, 2025-12-22 13:09:02+01, NULL, 'school_admin', 'School administrator who manages school users and library');
INSERT INTO "roles" ("id", "created_at", "updated_at", "deleted_at", "name", "description") VALUES ('6', 2025-12-22 11:39:02+01, 2025-12-22 13:09:02+01, NULL, 'teacher', 'Teacher who assigns books and tracks student progress');
INSERT INTO "roles" ("id", "created_at", "updated_at", "deleted_at", "name", "description") VALUES ('5', 2025-12-22 11:39:02+01, 2025-12-22 13:09:02+01, NULL, 'student', 'Student who reads books');


-- Table: users
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
  "id" bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "email" text NOT NULL,
  "username" text NOT NULL,
  "password_hash" text NOT NULL,
  "first_name" text,
  "last_name" text,
  "phone_number" text,
  "school_name" text,
  "school_category" text,
  "class_level" text,
  "department" text,
  "role_id" bigint,
  "is_active" boolean DEFAULT false,
  "is_email_verified" boolean DEFAULT false,
  "verification_token" text,
  "verification_token_expires" timestamp with time zone,
  "last_login" timestamp with time zone
);

-- Data for table: users
INSERT INTO "users" ("id", "created_at", "updated_at", "deleted_at", "email", "username", "password_hash", "first_name", "last_name", "phone_number", "school_name", "school_category", "class_level", "department", "role_id", "is_active", "is_email_verified", "verification_token", "verification_token_expires", "last_login") VALUES ('2', 2025-12-11 11:29:18+01, 2025-12-25 17:21:29+01, NULL, 'admin@readagain.com', 'admin', '$2a$10$fOoVIiYwORELVQctuZtMiuvKH7D2KPUVWjf7v3ZyKUotCmrygaUxa', 'Super', 'Admin', '0896655577', 'University Of Lagos', 'Tertiary', '', 'Computer Science', '8', true, true, '', NULL, 2025-12-25 17:21:29+01);


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
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'analytics.view', 'View analytics dashboard', 'analytics');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'analytics.export', 'Export analytics data', 'analytics');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'users.view', 'View users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'users.create', 'Create users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'users.edit', 'Edit users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'users.delete', 'Delete users', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'users.manage', 'Full user management', 'users');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'roles.view', 'View roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'roles.create', 'Create roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'roles.edit', 'Edit roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'roles.delete', 'Delete roles', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('12', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'roles.manage', 'Full role management', 'roles');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'audit_logs.view', 'View audit logs', 'audit');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('14', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'audit_logs.export', 'Export audit logs', 'audit');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'books.view', 'View books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('16', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'books.create', 'Create books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('17', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'books.edit', 'Edit books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('18', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'books.delete', 'Delete books', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('19', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'books.manage', 'Full book management', 'books');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reviews.view', 'View reviews', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('21', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reviews.moderate', 'Moderate reviews', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reviews.delete', 'Delete reviews', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('23', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reviews.manage', 'Full review management', 'reviews');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('24', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'orders.view', 'View orders', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('25', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'orders.edit', 'Edit orders', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('26', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'orders.delete', 'Delete orders', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('27', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'orders.manage', 'Full order management', 'orders');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('28', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'shipping.view', 'View shipping', 'shipping');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('29', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'shipping.manage', 'Manage shipping', 'shipping');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('30', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reading.view_analytics', 'View reading analytics', 'reading');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('31', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reading.manage', 'Manage reading data', 'reading');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('32', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reports.view', 'View reports', 'reports');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('33', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reports.generate', 'Generate reports', 'reports');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('34', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'reports.export', 'Export reports', 'reports');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('35', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'email_templates.view', 'View email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('36', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'email_templates.create', 'Create email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('37', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'email_templates.edit', 'Edit email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('38', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'email_templates.delete', 'Delete email templates', 'email');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('39', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'blog.view', 'View blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('40', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'blog.create', 'Create blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('41', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'blog.edit', 'Edit blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('42', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'blog.delete', 'Delete blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('43', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'blog.publish', 'Publish blog posts', 'blog');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('44', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'works.view', 'View works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('45', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'works.create', 'Create works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('46', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'works.edit', 'Edit works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('47', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'works.delete', 'Delete works', 'works');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('48', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'about.view', 'View about page', 'about');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('49', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'about.edit', 'Edit about page', 'about');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('50', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'contact.view', 'View contact messages', 'contact');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('51', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'contact.reply', 'Reply to contact messages', 'contact');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('52', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'contact.delete', 'Delete contact messages', 'contact');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('53', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'faq.view', 'View FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('54', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'faq.create', 'Create FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('55', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'faq.edit', 'Edit FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('56', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'faq.delete', 'Delete FAQs', 'faq');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('57', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'settings.view', 'View settings', 'settings');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('58', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'settings.edit', 'Edit settings', 'settings');
INSERT INTO "permissions" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "category") VALUES ('59', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'settings.manage', 'Full settings management', 'settings');


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
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('12', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('14', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('16', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('17', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('18', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('19', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('21', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('23', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('24', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('25', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'curl/8.7.1', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('26', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);
INSERT INTO "auth_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "ip_address", "user_agent", "success") VALUES ('27', 2025-12-25 17:21:29+01, 2025-12-25 17:21:29+01, NULL, '2', 'login', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36', true);


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
INSERT INTO "authors" ("id", "created_at", "updated_at", "deleted_at", "user_id", "business_name", "bio", "photo", "website", "email", "status") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'Classic Literature Press', 'Publisher of timeless classic literature for young readers', '', '', 'author@gmail.com', 'active');


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
INSERT INTO "books" ("id", "created_at", "updated_at", "deleted_at", "author_id", "title", "subtitle", "description", "short_description", "cover_image", "file_path", "file_size", "category_id", "isbn", "is_featured", "is_bestseller", "is_new_release", "status", "is_active", "pages", "language", "publisher", "publication_date", "library_count", "view_count", "seo_title", "seo_description", "seo_keywords") VALUES ('11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '1', 'Moby Dick', '', '', '', 'covers/0683010e-42a7-4af3-95c3-4e01b00e19f8.jpg', 'books/e068cce3-0e5c-4c71-bd04-cd0666b87723.epub', '816008', '2', '', true, false, false, 'published', true, '0', 'English', '', CURRENT_TIMESTAMP, '1', '15', '', '', '');


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
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Fiction', 'Fictional stories and novels', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Adventure', 'Adventure and exploration stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Fantasy', 'Fantasy and magical stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Mystery', 'Mystery and detective stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Science Fiction', 'Science fiction stories', 'active');
INSERT INTO "categories" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "status") VALUES ('6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Classic', 'Classic literature', 'active');


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
INSERT INTO "user_libraries" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "progress", "current_page", "last_read_at", "completed_at", "is_favorite", "rating") VALUES ('2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', '35.36493604213695', '0', NULL, NULL, false, '0');


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
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', CURRENT_TIMESTAMP, NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', CURRENT_TIMESTAMP, NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', CURRENT_TIMESTAMP, NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', CURRENT_TIMESTAMP, NULL, '35', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', CURRENT_TIMESTAMP, NULL, '0', '0', '0', '0');
INSERT INTO "reading_sessions" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "start_time", "end_time", "duration", "pages_read", "start_page", "end_page") VALUES ('6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', CURRENT_TIMESTAMP, NULL, '0', '0', '0', '0');


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
INSERT INTO "reading_goals" ("id", "created_at", "updated_at", "deleted_at", "user_id", "type", "target", "current", "start_date", "end_date", "is_completed", "book_id") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '2', 'books', '12', '0', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false, NULL);


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
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'update_review_status', 'review', '50', '', 'approved', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'toggle_review_featured', 'review', '50', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '55', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'update_review_status', 'review', '56', '', 'approved', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'toggle_review_featured', 'review', '56', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '44', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '54', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '42', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '53', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '52', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '51', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('12', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'toggle_review_featured', 'review', '50', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '50', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('14', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '49', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '48', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('16', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '47', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('17', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '46', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('18', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '45', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('19', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '43', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '41', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('21', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '40', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '39', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('23', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '38', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');
INSERT INTO "audit_logs" ("id", "created_at", "updated_at", "deleted_at", "user_id", "action", "entity_type", "entity_id", "old_value", "new_value", "ip_address", "user_agent") VALUES ('24', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', 'delete_review', 'review', '37', '', '', '127.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36');


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
INSERT INTO "reviews" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "rating", "comment", "status", "is_featured") VALUES ('56', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', '2', 'At least image upload is working', 'approved', true);


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
INSERT INTO "about_pages" ("id", "created_at", "updated_at", "deleted_at", "title", "content", "mission", "vision", "team_section", "values") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'About ReadAgain', 'Welcome to ReadAgain - Your digital reading platform.', '', '', '', '');


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
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'First Purchase', 'Purchase your first book', '📚', 'books_purchased', '1', '10', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Book Collector', 'Purchase 10 books', '📖', 'books_purchased', '10', '50', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Library Master', 'Purchase 50 books', '🏛️', 'books_purchased', '50', '200', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'First Finish', 'Complete your first book', '🎉', 'books_completed', '1', '20', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Bookworm', 'Complete 5 books', '🐛', 'books_completed', '5', '100', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Avid Reader', 'Complete 25 books', '📕', 'books_completed', '25', '300', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Reading Legend', 'Complete 100 books', '👑', 'books_completed', '100', '1000', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('8', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Getting Started', 'Read for 60 minutes', '⏱️', 'reading_minutes', '60', '15', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Dedicated Reader', 'Read for 10 hours', '⏰', 'reading_minutes', '600', '100', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Time Master', 'Read for 100 hours', '🕐', 'reading_minutes', '6000', '500', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('11', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'First Session', 'Complete your first reading session', '🎯', 'reading_sessions', '1', '5', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('12', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Consistent Reader', 'Complete 50 reading sessions', '📅', 'reading_sessions', '50', '150', NULL);
INSERT INTO "achievements" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "icon", "type", "target", "points", "category") VALUES ('13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Reading Habit', 'Complete 200 reading sessions', '🔥', 'reading_sessions', '200', '400', NULL);


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
INSERT INTO "groups" ("id", "created_at", "updated_at", "deleted_at", "name", "description", "created_by", "member_count") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'Science Club', '', '2', '1');


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
INSERT INTO "group_members" ("id", "created_at", "updated_at", "deleted_at", "group_id", "user_id") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '1', '2');


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
INSERT INTO "chat_rooms" ("id", "created_at", "updated_at", "deleted_at", "type", "name", "description", "group_id", "book_id", "created_by", "is_active", "last_message", "last_message_at") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'group', 'Group Chat', 'Group discussion', '1', NULL, '2', true, 'hi', CURRENT_TIMESTAMP);


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
INSERT INTO "chat_messages" ("id", "created_at", "updated_at", "deleted_at", "room_id", "user_id", "message", "message_type", "file_url", "file_name", "reply_to_id", "is_edited", "edited_at", "is_deleted") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '1', '2', 'hi', 'text', '', '', NULL, false, NULL, false);


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
INSERT INTO "chat_members" ("id", "created_at", "updated_at", "deleted_at", "room_id", "user_id", "role", "last_read_at", "unread_count", "is_muted", "joined_at") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '1', '2', 'admin', CURRENT_TIMESTAMP, '0', false, CURRENT_TIMESTAMP);


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
INSERT INTO "highlights" ("id", "created_at", "updated_at", "deleted_at", "user_id", "book_id", "text", "color", "start_offset", "end_offset", "context", "cfi_range") VALUES ('1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, '2', '11', 'Then God spake unto the fish; and from the shuddering cold and blackness of the sea, the whale came breeching up towards the warm and pleasant sun, and all the delights of air and earth; and ‘vomited out Jonah upon the dry land;’ when the word of the Lord came a second time; and Jonah, bruised and beaten—his ears, like two sea-shells, still multitudinously murmuring of the ocean—Jonah did the Almighty’s bidding. And what was that, shipmates? To preach the Truth to the face of Falsehood! That was it!', '#2196f3', '0', '504', 'Then God spake unto the fish; and from the shuddering cold and blackness of the sea, the whale came breeching up towards the warm and pleasant sun, and all the delights of air and earth; and ‘vomited out Jonah upon the dry land;’ when the word of the Lord came a second time; and Jonah, bruised and beaten—his ears, like two sea-shells, still multitudinously murmuring of the ocean—Jonah did the Almighty’s bidding. And what was that, shipmates? To preach the Truth to the face of Falsehood! That was it!', 'epubcfi(/6/6!/4/44,/5:967,/5:1519)');


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
INSERT INTO "notes" ("id", "created_at", "updated_at", "deleted_at", "content", "page_number", "position", "color", "user_id", "book_id", "page") VALUES (2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Test', NULL, NULL, 'yellow', '2', '11', '0');


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
INSERT INTO "licenses" ("id", "api_key", "license_key", "device_fingerprint", "is_active", "activated_at", "last_verified_at", "created_at", "updated_at", "is_valid", "grace_period_started_at", "grace_period_days", "is_locked") VALUES ('1', 'pk_live_ey62swS2VTFWQnARkDbAcB1pg7kgousxqXAeP4JBaNo', 'YtCRzVM9LtisY19x_li2rHGydpsThSNdYH61ePuP27w', '4a091dc77282e609f2e15f368dcd3377c7bb9336c50bd087620e1d44335b47d2', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false, CURRENT_TIMESTAMP, '3', false);

