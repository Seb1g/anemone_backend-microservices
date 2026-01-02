CREATE TABLE boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title VARCHAR(255) NOT NULL,
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_boards_user_id ON boards (user_id);

CREATE TABLE columns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    column_title VARCHAR(255) NOT NULL,
    board_id UUID NOT NULL,
    position INTEGER NOT NULL,
    FOREIGN KEY (board_id) REFERENCES boards (id) ON DELETE CASCADE,
    UNIQUE (board_id, position)
);

CREATE INDEX idx_columns_board_id ON columns (board_id);

CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    content TEXT NOT NULL,
    column_id UUID NOT NULL,
    position INTEGER NOT NULL,
    FOREIGN KEY (column_id) REFERENCES columns (id) ON DELETE CASCADE,
    UNIQUE (column_id, position)
);

CREATE INDEX idx_cards_column_id ON cards (column_id);