ALTER TABLE rooms_restrictions
DROP COLUMN reservations_id,
DROP COLUMN restrictions_id,
ADD reservation_id INTEGER NOT NULL REFERENCES reservations(id),
ADD restriction_id INTEGER NOT NULL REFERENCES restrictions(id);