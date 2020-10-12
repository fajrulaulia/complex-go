USE modulordb;

DROP TABLE IF EXISTS blogs;
CREATE TABLE blogs(
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    id_user INT,
    title VARCHAR(50),
    body MEDIUMTEXT,
    created_at DATETIME,
    updated_at DATETIME
);

