-- Create database
CREATE DATABASE gochat;

COMMENT ON DATABASE gochat
    IS 'Database for managing users and messages for gochat';

-- Connect to database
\c gochat

-- Initialize tables
CREATE TABLE users (
    user_id BIGSERIAL NOT NULL PRIMARY KEY,
    username VARCHAR (20) NOT NULL,
    remote_addr VARCHAR (64) NOT NULL,
    registered TIMESTAMP
);

CREATE TABLE messages (
    message_id BIGSERIAL NOT NULL  PRIMARY KEY,
    message VARCHAR (255) NOT NULL,
    user_id INT NOT NULL,
    sent TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);
