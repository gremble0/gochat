-- Create database
CREATE DATABASE gochat;

COMMENT ON DATABASE gochat
    IS 'Database for managing users and messages for gochat';

-- Connect to database
\c gochat

-- Initialize tables
CREATE TABLE users (
    username VARCHAR (20) NOT NULL,
    remote_addr VARCHAR (64) NOT NULL,
    registered TIMESTAMP,
    PRIMARY KEY (username, remote_addr)
);

CREATE TABLE messages (
    message_id BIGSERIAL NOT NULL PRIMARY KEY,
    message VARCHAR (255) NOT NULL,
    sender VARCHAR (20) NOT NULL,
    sender_addr VARCHAR (64) NOT NULL,
    sent TIMESTAMP,
    FOREIGN KEY (sender, sender_addr) REFERENCES users (username, remote_addr)
);
