-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "enabled" SET DEFAULT true, ALTER COLUMN "email_verified" SET DEFAULT false, ADD COLUMN "thumbnail_id" uuid NULL;
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
  "user_id" uuid NULL,
  "card_id" uuid NULL,
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
-- Create "file_uploads" table
CREATE TABLE "public"."file_uploads" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "extension" character varying(8) NOT NULL,
  "type" character varying(64) NOT NULL,
  "path" text NULL,
  "size" bigint NOT NULL,
  "additional_details" jsonb NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
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
-- Create "worksapces" table
CREATE TABLE "public"."worksapces" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "is_public" boolean NOT NULL DEFAULT true,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "workspace_permissions" table
CREATE TABLE "public"."workspace_permissions" (
  "worksapce_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("worksapce_id", "permission_id")
);
-- Modify "users" table
ALTER TABLE "public"."users" ADD CONSTRAINT "fk_users_thumbnail" FOREIGN KEY ("thumbnail_id") REFERENCES "public"."file_uploads" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "board_permissions" table
ALTER TABLE "public"."board_permissions" ADD CONSTRAINT "fk_board_permissions_board" FOREIGN KEY ("board_id") REFERENCES "public"."boards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_board_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "boards" table
ALTER TABLE "public"."boards" ADD CONSTRAINT "fk_boards_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_worksapces_boards" FOREIGN KEY ("workspace_id") REFERENCES "public"."worksapces" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
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
-- Modify "file_uploads" table
ALTER TABLE "public"."file_uploads" ADD CONSTRAINT "fk_users_files" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "labels" table
ALTER TABLE "public"."labels" ADD CONSTRAINT "fk_labels_board" FOREIGN KEY ("board_id") REFERENCES "public"."boards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "lists" table
ALTER TABLE "public"."lists" ADD CONSTRAINT "fk_boards_lists" FOREIGN KEY ("board_id") REFERENCES "public"."boards" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "worksapces" table
ALTER TABLE "public"."worksapces" ADD CONSTRAINT "fk_worksapces_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "workspace_permissions" table
ALTER TABLE "public"."workspace_permissions" ADD CONSTRAINT "fk_workspace_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_workspace_permissions_worksapce" FOREIGN KEY ("worksapce_id") REFERENCES "public"."worksapces" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
