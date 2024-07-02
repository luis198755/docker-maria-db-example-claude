CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    Program JSON
);

INSERT INTO products (name, price, Program) VALUES 
    ('Product 1', 30.99, NULL)
    ('Product 2', 29.99, NULL),
    ('Product 3', 39.99, NULL);