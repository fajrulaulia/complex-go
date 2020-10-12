USE modulordb;

DROP TABLE IF EXISTS users;
CREATE TABLE users(
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    username VARCHAR(50),
    phonenumber CHAR(15),
    email VARCHAR(50),
    password TEXT,
    status CHAR(10),
    created_at DATETIME,
    updated_at DATETIME
);
