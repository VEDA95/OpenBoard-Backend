CREATE TABLE open_board_workspace (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    date_created TIMESTAMP NOT NULL DEFAULT (now()),
    date_updated TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE open_board_board (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    user_id UUID NOT NULL,
    workspace_id UUID,
    name VARCHAR(255) NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
    date_updated TIMESTAMP WITH TIME ZONE,
    is_public BOOLEAN NOT NULL DEFAULT (TRUE)
);

CREATE TABLE open_board_board_list (
    id UUID UNIQUE PRIMARY KEY,
    board_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7),
    position INTEGER NOT NULL
);

CREATE TABLE open_board_board_label (
    id UUID UNIQUE PRIMARY KEY,
    board_id UUID NOT NULL,
    date_created TIMESTAMP NOT NULL DEFAULT (now()),
    date_updated TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7)
);

CREATE TABLE open_board_board_field (
    id UUID UNIQUE PRIMARY KEY,
    board_id UUID NOT NULL,
    date_created TIMESTAMP NOT NULL DEFAULT (now()),
    date_updated TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    field_properties JSONB NOT NULL
);

CREATE TABLE open_board_workspace_permissions (
    workspace_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (workspace_id, permission_id)
);

CREATE TABLE open_board_board_permissions (
    board_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (board_id, permission_id)
);

ALTER TABLE open_board_board ADD FOREIGN KEY (workspace_id) REFERENCES open_board_workspace (id) ON DELETE CASCADE;
ALTER TABLE open_board_board ADD FOREIGN KEY (user_id) REFERENCES open_board_user (id) ON DELETE CASCADE;
ALTER TABLE open_board_board_list ADD FOREIGN KEY (board_id) REFERENCES open_board_board (id) ON DELETE CASCADE;
ALTER TABLE open_board_board_label ADD FOREIGN KEY (board_id) REFERENCES open_board_board (id) ON DELETE CASCADE;
ALTER TABLE open_board_board_field ADD FOREIGN KEY (board_id) REFERENCES open_board_board (id) ON DELETE CASCADE;
ALTER TABLE open_board_workspace_permissions ADD FOREIGN KEY (workspace_id) REFERENCES open_board_workspace (id) ON DELETE CASCADE;
ALTER TABLE open_board_workspace_permissions ADD FOREIGN KEY (permission_id) REFERENCES open_board_role_permission (id) ON DELETE CASCADE;
ALTER TABLE open_board_board_permissions ADD FOREIGN KEY (board_id) REFERENCES open_board_board (id) ON DELETE CASCADE;
ALTER TABLE open_board_board_permissions ADD FOREIGN KEY (permission_id) REFERENCES open_board_role_permission (id) ON DELETE CASCADE;
