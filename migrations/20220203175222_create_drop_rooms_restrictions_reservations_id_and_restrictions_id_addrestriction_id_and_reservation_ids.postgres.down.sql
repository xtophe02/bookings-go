ALTER TABLE rooms_restrictions
DROP COLUMN reservation_id,
DROP COLUMN restriction_id,
ADD reservations_id INTEGER NOT NULL REFERENCES reservations(id),
ADD restrictions_id INTEGER NOT NULL REFERENCES restrictions(id);