-- +goose Up
-- +goose StatementBegin
CREATE TABLE "office_user" (
    office_id UUID,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (office_id, user_id),
    FOREIGN KEY (office_id) REFERENCES office (id),
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "office_user";
-- +goose StatementEnd
