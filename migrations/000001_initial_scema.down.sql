-- Drop tables in reverse order to respect foreign key constraints
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS credential_access;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;