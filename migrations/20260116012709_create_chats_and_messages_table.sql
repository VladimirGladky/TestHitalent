-- +goose Up
CREATE TABLE chats (
                       id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE messages (
                          id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                          chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
                          text TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
