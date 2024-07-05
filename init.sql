-- CREATE TABLE IF NOT EXISTS products (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     price DECIMAL(10, 2) NOT NULL,
--     Program JSON
-- );

-- INSERT INTO products (name, price, Program) VALUES 
--     ('Product 1', 30.99, NULL)
--     ('Product 2', 29.99, NULL),
--     ('Product 3', 39.99, NULL);

-- Create the table
CREATE TABLE IF NOT EXISTS centurion_projects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    centurion_id VARCHAR(15) UNIQUE NOT NULL,
    matricula VARCHAR(15) NOT NULL,
    ubicacion TEXT NOT NULL,
    proyecto VARCHAR(50) NOT NULL,
    fecha_entrega DATE NOT NULL,
    programa JSON,
    password VARCHAR(64) NOT NULL
);

-- Insert sample data
INSERT INTO centurion_projects (centurion_id, matricula, ubicacion, proyecto, fecha_entrega, programa, password)
VALUES
    ('Centurión_00001', 'TC-CL-24-0001', '16 oriente y 32 norte', 'Perseo', '2024-06-17', NULL, '84f9c98988b508b0'),
    ('Centurión_00002', 'TC-CL-24-0002', '16 oriente y 32 norte', 'Alquimia', '2024-06-19', NULL, 'a1b2c3d4e5f6g7h8'),
    ('Centurión_00003', 'TC-CL-24-0003', '16 oriente y 32 norte', 'Perseo', '2024-06-21', NULL, 'h8g7f6e5d4c3b2a1');

-- You can add more INSERT statements as needed