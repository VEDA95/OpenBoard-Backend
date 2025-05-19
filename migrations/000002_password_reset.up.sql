CREATE TABLE "open_board_password_reset_token" (
    id TEXT PRIMARY KEY,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
    expires_on TIMESTAMP WITH TIME ZONE NOT NULL,
    type VARCHAR(16) NOT NULL,
    user_id UUID NOT NULL
);

ALTER TABLE open_board_password_reset_token ADD FOREIGN KEY (user_id) REFERENCES open_board_user (id) ON DELETE CASCADE;