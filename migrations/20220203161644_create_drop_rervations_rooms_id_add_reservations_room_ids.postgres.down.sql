ALTER TABLE reservations
DROP COLUMN room_id,
ADD rooms_id INTEGER NOT NULL REFERENCES rooms(id);