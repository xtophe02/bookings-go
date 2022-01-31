CREATE TABLE rooms_restrictions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  room_id INTEGER NOT NULL REFERENCES rooms(id),
  reservations_id INTEGER NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
  restrictions_id INTEGER NOT NULL REFERENCES restrictions(id) ON DELETE CASCADE
);