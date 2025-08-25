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
    "enabled" boolean NOT NULL,
    "email_verified" boolean NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "uni_users_email" UNIQUE ("email"),
    CONSTRAINT "uni_users_username" UNIQUE ("username")
);
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE INDEX "idx_users_username" ON "public"."users" ("username");
-- Create "password_reset_tokens" table
CREATE TABLE "public"."password_reset_tokens" (
    "id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL,
    "expires_on" timestamptz NOT NULL,
    "type" character varying(16) NOT NULL,
    "token" character varying(6) NOT NULL,
    "user_id" uuid NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_password_reset_tokens_user" FOREIGN KEY (
        "user_id"
    ) REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "permissions" table
CREATE TABLE "public"."permissions" (
    "id" uuid NOT NULL,
    "path" text NULL,
    PRIMARY KEY ("id")
);
-- Create "roles" table
CREATE TABLE "public"."roles" (
    "id" uuid NOT NULL,
    "name" character varying(255) NOT NULL,
    PRIMARY KEY ("id")
);
-- Create "role_permissions" table
CREATE TABLE "public"."role_permissions" (
    "role_id" uuid NOT NULL,
    "permission_id" uuid NOT NULL,
    PRIMARY KEY ("role_id", "permission_id"),
    CONSTRAINT "fk_role_permissions_permission" FOREIGN KEY (
        "permission_id"
    ) REFERENCES "public"."permissions" (
        "id"
    ) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT "fk_role_permissions_role" FOREIGN KEY (
        "role_id"
    ) REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
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
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_users_sessions" FOREIGN KEY (
        "user_id"
    ) REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create "user_roles" table
CREATE TABLE "public"."user_roles" (
    "user_id" uuid NOT NULL,
    "role_id" uuid NOT NULL,
    PRIMARY KEY ("user_id", "role_id"),
    CONSTRAINT "fk_user_roles_role" FOREIGN KEY (
        "role_id"
    ) REFERENCES "public"."roles" (
        "id"
    ) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT "fk_user_roles_user" FOREIGN KEY (
        "user_id"
    ) REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
