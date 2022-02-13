ALTER TABLE room_restrictions
DROP CONSTRAINT rooms_restrictions_reservation_id_fkey,
ADD CONSTRAINT room_restrictions_reservation_id_fkey
FOREIGN KEY (reservation_id)
REFERENCES reservations(id)
ON DELETE CASCADE;