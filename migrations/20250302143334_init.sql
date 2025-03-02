-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    password            VARCHAR(255) NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tests(
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    description         TEXT,
    data                TEXT,
    author_id           uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS test_results(
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    test_id             uuid NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
    user_id             uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    result              TEXT NOT NULL,
    score               INTEGER NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS test_results;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
