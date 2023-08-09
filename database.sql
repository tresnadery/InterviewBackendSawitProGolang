/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	phone_number VARCHAR (13) UNIQUE NOT NULL,
	full_name VARCHAR ( 60 ) NOT NULL,
	password text,
	password_salt VARCHAR (15),
	successfully_login int DEFAULT 0,
	last_login timestamptz,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
	created_by uuid,
	modified_at timestamptz,
	modified_by uuid,
	deleted_at timestamptz
);
