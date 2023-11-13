CREATE DATABASE gochat
    WITH OWNER DEFAULT;

CREATE TABLE users (
    user_id serial PRIMARY KEY,
    username VARCHAR (20) NOT NULL,
    remote_addr VARCHAR (64) NOT NULL,
    registered TIMESTAMP NOT NULL
);

-- Maybe combination of user and timestamp for primary key?
CREATE TABLE messages (
    message_id serial PRIMARY KEY,
    message VARCHAR (255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);
