-- Drop tables if they already exist
DROP TABLE IF EXISTS goal;
DROP TABLE IF EXISTS book;
DROP TABLE IF EXISTS user;

-- user table
CREATE TABLE user (
    id INT NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- goal table
CREATE TABLE goal (
    user_id INT NOT NULL,
    type ENUM('books', 'pages') NOT NULL,
    frequency ENUM('daily', 'monthly') NOT NULL,
    value INT NOT NULL CHECK (value >= 1),
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

-- book table
CREATE TABLE book (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(50) NOT NULL,
    rating FLOAT CHECK (rating >= 0 AND rating <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);
