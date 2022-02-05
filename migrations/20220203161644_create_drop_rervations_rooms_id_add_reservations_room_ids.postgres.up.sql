ALTER TABLE reservations
DROP COLUMN rooms_id,
ADD room_id INTEGER NOT NULL REFERENCES rooms(id);