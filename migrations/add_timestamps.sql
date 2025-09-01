CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_timestamp
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_guide_timestamp
BEFORE UPDATE ON guides
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_booking_timestamp
BEFORE UPDATE ON bookings
FOR EACH ROW EXECUTE FUNCTION update_timestamp();
