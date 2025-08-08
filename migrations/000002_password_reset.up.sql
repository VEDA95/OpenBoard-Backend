CREATE TABLE "open_board_password_reset_token" (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
    expires_on TIMESTAMP WITH TIME ZONE NOT NULL,
    type VARCHAR(16) NOT NULL,
    token VARCHAR(6) NOT NULL,
    user_id UUID NOT NULL
);

ALTER TABLE open_board_password_reset_token ADD FOREIGN KEY (user_id) REFERENCES open_board_user (id) ON DELETE CASCADE;