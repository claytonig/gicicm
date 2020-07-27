CREATE USER goicm with password 'pass';
ALTER ROLE goicm with superuser;
CREATE TABLE users (
    id        SERIAL PRIMARY KEY,
    name       varchar(40) NOT NULL,
    email         varchar(40) UNIQUE NOT NULL,
    password    varchar(200) NOT NULL
);
GRANT ALL PRIVILEGES ON TABLE users TO goicm;
GRANT ALL ON SEQUENCE users_id_seq TO goicm;

