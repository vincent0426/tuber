CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE TABLE users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  bio TEXT,
  accept_notification BOOLEAN NOT NULL DEFAULT TRUE,
  sub TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE car (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  driver_id UUID NOT NULL,
  license TEXT NOT NULL,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (driver_id) REFERENCES users(id)
);
CREATE TABLE driver (
  user_id UUID PRIMARY KEY,
  car_id UUID,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (car_id) REFERENCES car(id)
);
CREATE TABLE location (
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
  FOREIGN KEY (source_id) REFERENCES location(id),
  FOREIGN KEY (destination_id) REFERENCES location(id)
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
  FOREIGN KEY (station_source_id) REFERENCES location(id),
  FOREIGN KEY (station_destination_id) REFERENCES location(id)
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
-- CREATE VIEW user_tokens_view AS
-- SELECT users.*
-- FROM users
--   JOIN tokens ON users.id = tokens.user_id
-- WHERE tokens.expiry > NOW();