CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE TABLE users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  bio TEXT,
  accept_notification BOOLEAN NOT NULL DEFAULT TRUE,
  sub TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE driver (
  user_id UUID PRIMARY KEY,
  license TEXT NOT NULL,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  brand TEXT NOT NULL,
  model TEXT NOT NULL,
  color TEXT NOT NULL,
  plate TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE locations (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  name TEXT NOT NULL,
  address TEXT,
  coordinates GEOMETRY(Point, 4326)
);
CREATE TABLE trip (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  driver_id UUID NOT NULL,
  passenger_limit INTEGER NOT NULL DEFAULT 1,
  source_id UUID NOT NULL,
  destination_id UUID NOT NULL,
  status TEXT CHECK (status IN ('not_start', 'in_trip', 'finished')) DEFAULT 'not_start',
  start_time TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (driver_id) REFERENCES driver(user_id),
  FOREIGN KEY (source_id) REFERENCES locations(id),
  FOREIGN KEY (destination_id) REFERENCES locations(id)
);
CREATE TABLE chat_history (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  trip_id UUID NOT NULL,
  sender_id UUID NOT NULL,
  msg_content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (sender_id) REFERENCES users(id)
);
CREATE TABLE rating (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  trip_id UUID NOT NULL,
  commenter_id UUID NOT NULL,
  comment TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (commenter_id) REFERENCES users(id)
);
CREATE TABLE trip_station (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  trip_id UUID NOT NULL,
  name TEXT NOT NULL,
  FOREIGN KEY (trip_id) REFERENCES trip(id)
);
CREATE TABLE trip_passenger (
  trip_id UUID NOT NULL,
  passenger_id UUID NOT NULL,
  station_source_id UUID NOT NULL,
  station_destination_id UUID NOT NULL,
  status TEXT CHECK (status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (trip_id, passenger_id),
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (passenger_id) REFERENCES users(id),
  FOREIGN KEY (station_source_id) REFERENCES locations(id),
  FOREIGN KEY (station_destination_id) REFERENCES locations(id)
);
CREATE TABLE alert (
  trip_id UUID PRIMARY KEY,
  passenger_id UUID NOT NULL,
  comment TEXT,
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (passenger_id) REFERENCES users(id)
);
CREATE TABLE report (
  trip_id UUID NOT NULL,
  complainant UUID NOT NULL,
  defendant UUID NOT NULL,
  comment TEXT,
  PRIMARY KEY (trip_id, complainant, defendant),
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (complainant) REFERENCES users(id),
  FOREIGN KEY (defendant) REFERENCES users(id)
);
CREATE TABLE favorite_driver (
  user_id UUID NOT NULL,
  driver_id UUID NOT NULL,
  note TEXT,
  PRIMARY KEY (user_id, driver_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (driver_id) REFERENCES driver(user_id)
);
-- ------------------ view ------------------
-- trip_view: trip + driverinfo + tripinfo + station_source + station_destination
CREATE VIEW trip_view AS
SELECT trip.*,
  udriver.name AS driver_name,
  driver.brand AS driver_brand,
  driver.model AS driver_model,
  driver.color AS driver_color,
  driver.plate AS driver_plate,
  location_source.name AS source_name,
  location_source.address AS source_address,
  location_source.coordinates AS source_coordinates,
  location_destination.name AS destination_name,
  location_destination.address AS destination_address,
  location_destination.coordinates AS destination_coordinates
FROM trip
  JOIN users AS udriver ON trip.driver_id = udriver.id
  JOIN driver ON trip.driver_id = driver.user_id
  JOIN locations AS location_source ON trip.source_id = location_source.id
  JOIN locations AS location_destination ON trip.destination_id = location_destination.id;
-- trip_passenger_view
CREATE VIEW trip_passenger_view AS
SELECT trip_passenger.*,
  udriver.name AS driver_name,
  driver.brand AS driver_brand,
  driver.model AS driver_model,
  driver.color AS driver_color,
  driver.plate AS driver_plate,
  location_source.name AS source_name,
  location_source.address AS source_address,
  location_source.coordinates AS source_coordinates,
  location_destination.name AS destination_name,
  location_destination.address AS destination_address,
  location_destination.coordinates AS destination_coordinates,
  passenger_location_source.name AS passenger_location_source_name,
  passenger_location_source.address AS passenger_location_source_address,
  passenger_location_source.coordinates AS passenger_location_source_coordinates,
  passenger_location_destination.name AS passenger_location_destination_name,
  passenger_location_destination.address AS passenger_location_destination_address,
  passenger_location_destination.coordinates AS passenger_location_destination_coordinates
FROM trip_passenger
  JOIN trip ON trip_passenger.trip_id = trip.id
  JOIN users AS udriver ON trip.driver_id = udriver.id
  JOIN locations AS location_source ON trip.source_id = location_source.id
  JOIN locations AS location_destination ON trip.destination_id = location_destination.id
  JOIN locations AS passenger_location_source ON trip_passenger.station_source_id = passenger_location_source.id
  JOIN locations AS passenger_location_destination ON trip_passenger.station_destination_id = passenger_location_destination.id;