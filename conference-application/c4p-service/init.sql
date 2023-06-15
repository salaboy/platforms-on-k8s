CREATE TABLE proposals(
    id VARCHAR PRIMARY KEY NOT NULL,
    title VARCHAR NOT NULL,
    description TEXT NOT NULL,
    author VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    approved boolean ,
    status VARCHAR NOT NULL
);