-- Create the database if it doesn't exist
CREATE DATABASE exampledb;

-- Create a new user with a password
CREATE USER new_user WITH ENCRYPTED PASSWORD '';

ALTER USER new_user WITH SUPERUSER;