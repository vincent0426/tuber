CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE TABLE users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  image_url TEXT,
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
  place_id TEXT NOT NULL,
  lat_lon GEOGRAPHY(POINT, 4326) NOT NULL
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
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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
  rate INTEGER NOT NULL,
  trip_id UUID NOT NULL,
  commenter_id UUID NOT NULL,
  comment TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (commenter_id) REFERENCES users(id)
);
CREATE TABLE trip_location (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  trip_id UUID NOT NULL,
  location_id UUID NOT NULL,
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (location_id) REFERENCES locations(id)
);
CREATE TABLE trip_passenger (
  trip_id UUID NOT NULL,
  passenger_id UUID NOT NULL,
  source_id UUID NOT NULL,
  destination_id UUID NOT NULL,
  status TEXT DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  roles TEXT DEFAULT 'passenger',
  PRIMARY KEY (trip_id, passenger_id),
  FOREIGN KEY (trip_id) REFERENCES trip(id),
  FOREIGN KEY (passenger_id) REFERENCES users(id),
  FOREIGN KEY (source_id) REFERENCES locations(id),
  FOREIGN KEY (destination_id) REFERENCES locations(id)
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
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  driver_id UUID NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (driver_id) REFERENCES driver(user_id)
);