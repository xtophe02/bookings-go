ALTER TABLE reservations
DROP COLUMN user_id,
ADD COLUMN first_name VARCHAR(30) NOT NULL,
ADD COLUMN last_name VARCHAR(30) NOT NULL,
ADD COLUMN email VARCHAR(50) NOT NULL,
ADD COLUMN phone VARCHAR(20);