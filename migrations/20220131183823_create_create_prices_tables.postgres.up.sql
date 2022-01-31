CREATE TABLE prices (
  id SERIAL PRIMARY KEY,
  winter_price INT,
  summer_price INT,
  room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE
);