CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "open_board_user" (
    "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "date_created" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
    "date_updated" TIMESTAMP WITH TIME ZONE,
    "last_login" TIMESTAMP WITH TIME ZONE,
    "username" VARCHAR(255) UNIQUE NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "first_name" VARCHAR(255),
    "last_name" VARCHAR(255),
    "hashed_password" TEXT,
    "enabled" BOOLEAN NOT NULL DEFAULT (true),
    "email_verified" BOOLEAN NOT NULL DEFAULT (false)
);

CREATE TABLE "open_board_user_session" (
   id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
   user_id UUID NOT NULL,
   date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
   date_updated TIMESTAMP WITH TIME ZONE,
   expires_on TIMESTAMP WITH TIME ZONE NOT NULL,
   refresh_expires_on TIMESTAMP WITH TIME ZONE,
   session_type VARCHAR(32) NOT NULL,
   access_token TEXT UNIQUE,
   refresh_token TEXT UNIQUE,
   ip_address VARCHAR(255) NOT NULL,
   user_agent VARCHAR(255) NOT NULL,
   additional_info JSONB
);

CREATE TABLE "open_board_role" (
   id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
   name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE "open_board_role_permission" (
  id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  path VARCHAR(255) NOT NULL
);

CREATE TABLE "open_board_role_permissions" (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE "open_board_user_roles" (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

ALTER TABLE open_board_user_session ADD FOREIGN KEY (user_id) REFERENCES open_board_user (id);
ALTER TABLE open_board_role_permissions ADD FOREIGN KEY (role_id) REFERENCES open_board_role (id);
ALTER TABLE open_board_role_permissions ADD FOREIGN KEY (permission_id) REFERENCES open_board_role_permission (id);
ALTER TABLE open_board_user_roles ADD FOREIGN KEY (user_id) REFERENCES open_board_user (id);
ALTER TABLE open_board_user_roles ADD FOREIGN KEY (role_id) REFERENCES open_board_role (id);