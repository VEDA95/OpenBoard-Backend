-- Create "auth_settings" table
CREATE TABLE "public"."auth_settings" (
  "id" bigserial NOT NULL,
  "updated_at" timestamptz NULL,
  "allow_public_registration" boolean NULL DEFAULT true,
  "require_email_verification" boolean NULL DEFAULT true,
  "require_admin_approval" boolean NULL DEFAULT false,
  "default_user_role" character varying(50) NULL DEFAULT 'user',
  "registration_welcome_email" boolean NULL DEFAULT true,
  "allow_user_invitations" boolean NULL DEFAULT true,
  "invitation_expiry" bigint NULL DEFAULT 604800,
  "invite_only_mode" boolean NULL DEFAULT false,
  "session_timeout" bigint NULL DEFAULT 3600,
  "session_idle_timeout" bigint NULL DEFAULT 1800,
  "remember_me_duration" bigint NULL DEFAULT 86400,
  "max_login_attempts" bigint NULL DEFAULT 5,
  "lockout_duration" bigint NULL DEFAULT 18000,
  "two_factor_authentication" boolean NULL DEFAULT false,
  "pending_two_factor_auth_expiry" bigint NULL DEFAULT 900,
  "two_factor_required" boolean NULL DEFAULT false,
  "enable_o_auth" boolean NULL DEFAULT false,
  "cors_domain" character varying(255) NOT NULL,
  "web_authn_rp_id" character varying(255) NULL DEFAULT 'localhost',
  "web_authn_rp_display_name" character varying(255) NULL DEFAULT 'Open Board',
  "web_authn_rp_origins" text NULL DEFAULT '',
  PRIMARY KEY ("id"),
  CONSTRAINT "check_single_row_auth" CHECK (id = 1)
);
-- Create "board_permissions" table
CREATE TABLE "public"."board_permissions" (
  "board_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("board_id", "permission_id")
);
-- Create "boards" table
CREATE TABLE "public"."boards" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "is_public" boolean NOT NULL DEFAULT true,
  "user_id" uuid NOT NULL,
  "workspace_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "card_activities" table
CREATE TABLE "public"."card_activities" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "activity" character varying(255) NOT NULL,
  "user_id" uuid NOT NULL,
  "card_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "card_attachments" table
CREATE TABLE "public"."card_attachments" (
  "card_id" uuid NOT NULL,
  "file_upload_id" uuid NOT NULL,
  PRIMARY KEY ("card_id", "file_upload_id")
);
-- Create "card_labels" table
CREATE TABLE "public"."card_labels" (
  "card_id" uuid NOT NULL,
  "label_id" uuid NOT NULL,
  PRIMARY KEY ("card_id", "label_id")
);
-- Create "cards" table
CREATE TABLE "public"."cards" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "color" character varying(7) NULL,
  "position" bigint NOT NULL,
  "reminder_date" timestamptz NULL,
  "due_date" timestamptz NULL,
  "time_spent" bigint NULL,
  "estimated_time_spent" bigint NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "list_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "check_list_items" table
CREATE TABLE "public"."check_list_items" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" text NOT NULL,
  "is_checked" boolean NOT NULL DEFAULT true,
  "position" bigint NOT NULL,
  "card_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "comment" text NOT NULL,
  "user_id" uuid NOT NULL,
  "card_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "email_settings" table
CREATE TABLE "public"."email_settings" (
  "id" bigserial NOT NULL,
  "updated_at" timestamptz NULL,
  "email_provider" character varying(20) NULL DEFAULT 'smtp',
  "email_enabled" boolean NULL DEFAULT true,
  "email_from_address" character varying(255) NULL,
  "email_from_name" character varying(100) NULL,
  "smtp_host" character varying(255) NULL,
  "smtp_port" bigint NULL DEFAULT 587,
  "smtp_username" character varying(255) NULL,
  "smtp_password" character varying(255) NULL,
  "smtp_encryption" character varying(10) NULL DEFAULT 'tls',
  "smtp_auth_method" character varying(20) NULL DEFAULT 'plain',
  "smtp_verify_ssl" boolean NULL DEFAULT true,
  "send_grid_api_key" character varying(255) NULL,
  "mailgun_api_key" character varying(255) NULL,
  "mailgun_domain" character varying(255) NULL,
  "ses_access_key_id" character varying(255) NULL,
  "ses_secret_access_key" character varying(255) NULL,
  "ses_region" character varying(50) NULL DEFAULT 'us-east-1',
  "postmark_server_token" character varying(255) NULL,
  "postmark_account_token" character varying(255) NULL,
  "email_footer_text" text NULL,
  "enable_email_notifications" boolean NULL DEFAULT true,
  "send_password_reset_email" boolean NULL DEFAULT true,
  "notify_on_card_assigned" boolean NULL DEFAULT true,
  "notify_on_card_comment" boolean NULL DEFAULT true,
  "notify_on_card_due" boolean NULL DEFAULT true,
  "notify_on_board_invite" boolean NULL DEFAULT true,
  "notify_on_mention" boolean NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "check_single_row_email" CHECK (id = 1)
);
-- Create "email_verification_tokens" table
CREATE TABLE "public"."email_verification_tokens" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "expires_on" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "external_auth_providers" table
CREATE TABLE "public"."external_auth_providers" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "client_id" character varying(255) NOT NULL,
  "client_secret" text NOT NULL,
  "required_email_domain" character varying(255) NULL,
  "auth_url" character varying(255) NOT NULL,
  "login_url" character varying(255) NOT NULL,
  "user_info_url" character varying(255) NOT NULL,
  "logout_url" character varying(255) NULL,
  "use_pkce" boolean NOT NULL DEFAULT false,
  "default_login_method" boolean NOT NULL DEFAULT false,
  "self_registration_enabled" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create "external_provider_permissions" table
CREATE TABLE "public"."external_provider_permissions" (
  "external_auth_provider_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("external_auth_provider_id", "permission_id")
);
-- Create "file_uploads" table
CREATE TABLE "public"."file_uploads" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "extension" character varying(8) NOT NULL,
  "type" character varying(64) NOT NULL,
  "path" text NOT NULL,
  "size" bigint NOT NULL,
  "additional_details" jsonb NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "general_settings" table
CREATE TABLE "public"."general_settings" (
  "id" bigserial NOT NULL,
  "updated_at" timestamptz NULL,
  "app_name" character varying(100) NOT NULL DEFAULT 'Kanban Board',
  "app_url" character varying(255) NOT NULL DEFAULT 'http://localhost:3000',
  "app_logo_id" uuid NULL,
  "app_favicon_id" uuid NULL,
  "app_description" text NULL,
  "show_announcement_banner" boolean NULL DEFAULT false,
  "announcement_message" text NULL,
  "announcement_type" character varying(20) NULL DEFAULT 'info',
  "default_language" character varying(10) NULL DEFAULT 'en',
  "default_timezone" character varying(50) NULL DEFAULT 'UTC',
  "default_items_per_page" bigint NULL DEFAULT 25,
  "max_file_size" bigint NULL DEFAULT 1024,
  PRIMARY KEY ("id"),
  CONSTRAINT "check_single_row_general" CHECK (id = 1)
);
-- Create "labels" table
CREATE TABLE "public"."labels" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "color" character varying(7) NOT NULL,
  "board_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "lists" table
CREATE TABLE "public"."lists" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "color" character varying(7) NULL,
  "position" bigint NOT NULL,
  "board_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "multi_auth_methods" table
CREATE TABLE "public"."multi_auth_methods" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "type" character varying(255) NOT NULL,
  "credentials" jsonb NULL,
  "user_id" uuid NULL,
  PRIMARY KEY ("id")
);
-- Create "password_reset_tokens" table
CREATE TABLE "public"."password_reset_tokens" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "expires_on" timestamptz NOT NULL,
  "type" character varying(16) NOT NULL,
  "token" character varying(6) NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "permissions" table
CREATE TABLE "public"."permissions" (
  "id" uuid NOT NULL,
  "path" text NULL,
  PRIMARY KEY ("id")
);
-- Create "role_permissions" table
CREATE TABLE "public"."role_permissions" (
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("role_id", "permission_id")
);
-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "sessions" table
CREATE TABLE "public"."sessions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "expires_on" timestamptz NOT NULL,
  "refresh_expires_on" timestamptz NULL,
  "type" character varying(32) NOT NULL,
  "remember_me" boolean NOT NULL DEFAULT false,
  "access_token" text NOT NULL,
  "refresh_token" text NULL,
  "ip_address" character varying(255) NOT NULL,
  "user_agent" character varying(255) NOT NULL,
  "additional_info" jsonb NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "user_roles" table
CREATE TABLE "public"."user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "role_id")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "last_login" timestamptz NULL,
  "username" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "first_name" character varying(255) NULL,
  "last_name" character varying(255) NULL,
  "hashed_password" text NOT NULL,
  "thumbnail_id" uuid NULL,
  "enabled" boolean NOT NULL DEFAULT true,
  "email_verified" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email"),
  CONSTRAINT "uni_users_username" UNIQUE ("username")
);
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE INDEX "idx_users_username" ON "public"."users" ("username");
-- Create "workspace_permissions" table
CREATE TABLE "public"."workspace_permissions" (
  "workspace_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("workspace_id", "permission_id")
);
-- Create "workspaces" table
CREATE TABLE "public"."workspaces" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "is_public" boolean NOT NULL DEFAULT true,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Modify "board_permissions" table
ALTER TABLE "public"."board_permissions" ADD CONSTRAINT "fk_board_permissions_board" FOREIGN KEY ("board_id") REFERENCES "public"."boards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_board_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "boards" table
ALTER TABLE "public"."boards" ADD CONSTRAINT "fk_boards_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_workspaces_boards" FOREIGN KEY ("workspace_id") REFERENCES "public"."workspaces" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "card_activities" table
ALTER TABLE "public"."card_activities" ADD CONSTRAINT "fk_card_activities_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_cards_activities" FOREIGN KEY ("card_id") REFERENCES "public"."cards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "card_attachments" table
ALTER TABLE "public"."card_attachments" ADD CONSTRAINT "fk_card_attachments_card" FOREIGN KEY ("card_id") REFERENCES "public"."cards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_card_attachments_file_upload" FOREIGN KEY ("file_upload_id") REFERENCES "public"."file_uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "card_labels" table
ALTER TABLE "public"."card_labels" ADD CONSTRAINT "fk_card_labels_card" FOREIGN KEY ("card_id") REFERENCES "public"."cards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_card_labels_label" FOREIGN KEY ("label_id") REFERENCES "public"."labels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "cards" table
ALTER TABLE "public"."cards" ADD CONSTRAINT "fk_lists_cards" FOREIGN KEY ("list_id") REFERENCES "public"."lists" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "check_list_items" table
ALTER TABLE "public"."check_list_items" ADD CONSTRAINT "fk_check_list_items_card" FOREIGN KEY ("card_id") REFERENCES "public"."cards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "comments" table
ALTER TABLE "public"."comments" ADD CONSTRAINT "fk_cards_comments" FOREIGN KEY ("card_id") REFERENCES "public"."cards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_comments_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "email_verification_tokens" table
ALTER TABLE "public"."email_verification_tokens" ADD CONSTRAINT "fk_email_verification_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "external_provider_permissions" table
ALTER TABLE "public"."external_provider_permissions" ADD CONSTRAINT "fk_external_provider_permissions_external_auth_provider" FOREIGN KEY ("external_auth_provider_id") REFERENCES "public"."external_auth_providers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_external_provider_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "file_uploads" table
ALTER TABLE "public"."file_uploads" ADD CONSTRAINT "fk_users_files" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "general_settings" table
ALTER TABLE "public"."general_settings" ADD CONSTRAINT "fk_general_settings_app_favicon" FOREIGN KEY ("app_favicon_id") REFERENCES "public"."file_uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_general_settings_app_logo" FOREIGN KEY ("app_logo_id") REFERENCES "public"."file_uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "labels" table
ALTER TABLE "public"."labels" ADD CONSTRAINT "fk_labels_board" FOREIGN KEY ("board_id") REFERENCES "public"."boards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "lists" table
ALTER TABLE "public"."lists" ADD CONSTRAINT "fk_boards_lists" FOREIGN KEY ("board_id") REFERENCES "public"."boards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "multi_auth_methods" table
ALTER TABLE "public"."multi_auth_methods" ADD CONSTRAINT "fk_multi_auth_methods_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "password_reset_tokens" table
ALTER TABLE "public"."password_reset_tokens" ADD CONSTRAINT "fk_password_reset_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "role_permissions" table
ALTER TABLE "public"."role_permissions" ADD CONSTRAINT "fk_role_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_role_permissions_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "sessions" table
ALTER TABLE "public"."sessions" ADD CONSTRAINT "fk_users_sessions" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Modify "user_roles" table
ALTER TABLE "public"."user_roles" ADD CONSTRAINT "fk_user_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_user_roles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "users" table
ALTER TABLE "public"."users" ADD CONSTRAINT "fk_users_thumbnail" FOREIGN KEY ("thumbnail_id") REFERENCES "public"."file_uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "workspace_permissions" table
ALTER TABLE "public"."workspace_permissions" ADD CONSTRAINT "fk_workspace_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_workspace_permissions_workspace" FOREIGN KEY ("workspace_id") REFERENCES "public"."workspaces" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "workspaces" table
ALTER TABLE "public"."workspaces" ADD CONSTRAINT "fk_workspaces_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
