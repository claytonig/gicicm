CREATE USER goicm with password 'pass';
ALTER ROLE goicm with superuser;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id        SERIAL PRIMARY KEY,
    name       varchar(40) NOT NULL,
    email         varchar(40) UNIQUE NOT NULL,
    password    varchar(200) NOT NULL
);
GRANT ALL PRIVILEGES ON TABLE users TO goicm;
GRANT ALL ON SEQUENCE users_id_seq TO goicm;

INSERT into users(name, email, password) VALUES ('user to be deleted','delete@me.com','$2a$10easdasd$21hx81mFFbdlAn4Q9iEw5eYg86MPugTrd5HSxbw0s.PtlUB4XQlLu');
INSERT into users(name, email, password) VALUES ('superadmin','clayton@test.com','$2a$10$21hx81mFFbdlAn4Q9iEw5eYg86MPugTrd5HSxbw0s.PtlUB4XQlLu');
INSERT into users(name, email, password) VALUES ('test user 2','testtwo@mail.com','123123');
INSERT into users(name, email, password) VALUES ('test user 1','test@mail.com','123123');
