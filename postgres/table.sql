-- Create the schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS message_schema;

-- Create the 'messages' table within 'message_schema'
CREATE TABLE IF NOT EXISTS message_schema.messages (
    id SERIAL PRIMARY KEY,
    message_body TEXT NOT NULL
);