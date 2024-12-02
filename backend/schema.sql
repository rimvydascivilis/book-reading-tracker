-- Drop tables if they already exist
DROP TABLE IF EXISTS progress;
Drop TABLE IF EXISTS reading;
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

-- reading table
CREATE TABLE reading (
    id INT NOT NULL AUTO_INCREMENT,
    book_id INT NOT NULL,
    user_id INT NOT NULL,
    total_pages INT NOT NULL CHECK (total_pages > 0),
    link VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);

-- progress table
CREATE TABLE progress (
    id INT NOT NULL AUTO_INCREMENT,
    reading_id INT NOT NULL,
    user_id INT NOT NULL,
    pages INT NOT NULL CHECK (pages > 0),
    reading_date DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (reading_id) REFERENCES reading(id) ON DELETE CASCADE
);

-- list table
CREATE TABLE list (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

-- list_item table
CREATE TABLE list_item (
    id INT NOT NULL AUTO_INCREMENT,
    list_id INT NOT NULL,
    book_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (list_id) REFERENCES list(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);