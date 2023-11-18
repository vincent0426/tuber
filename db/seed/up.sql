-- Insert data into 'users' table
INSERT INTO users (name, email, image_url, bio)
VALUES (
    'John Doe',
    'john.doe@example.com',
    'https://example.com/john.doe.jpg',
    'Loves traveling.'
  ),
  (
    'Jane Smith',
    'jane.smith@example.com',
    'https://example.com/jane.smith.jpg',
    'Enjoys long drives.'
  );
-- Insert data into 'driver' table
INSERT INTO driver (
    user_id,
    brand,
    model,
    color,
    plate,
    license,
    verified
  )
SELECT id,
  'Toyota',
  'Camry',
  'White',
  'ABC123',
  'DL123456',
  TRUE
FROM users
WHERE email = 'john.doe@example.com';
-- Insert data into 'locations' table
INSERT INTO locations (name, place_id, lat_lon)
VALUES (
    'Location 1',
    'ChIJpTvG15DL1IkRd8S0KlBVNTI',
    ST_GeomFromText('POINT(-71.060316 48.432044)', 4326)
  ),
  (
    'Location 2',
    'D9iJyWEHuEmuEmsRm9hTkapTCrk',
    ST_GeomFromText('POINT(-69.445469 43.769196)', 4326)
  );
-- Insert a trip (assuming John Doe is the driver)
INSERT INTO trip (
    driver_id,
    passenger_limit,
    source_id,
    destination_id,
    start_time
  )
SELECT (
    SELECT user_id
    FROM driver
    WHERE license = 'DL123456'
  ),
  3,
  (
    SELECT id
    FROM locations
    WHERE name = 'Location 1'
  ),
  (
    SELECT id
    FROM locations
    WHERE name = 'Location 2'
  ),
  '2023-01-01 08:00:00';
-- Insert data into 'chat_history'
INSERT INTO chat_history (trip_id, sender_id, msg_content)
SELECT (
    SELECT id
    FROM trip
    LIMIT 1
  ), (
    SELECT id
    FROM users
    WHERE email = 'john.doe@example.com'
  ),
  'Hello, I am your driver.';
-- Insert data into 'rating'
INSERT INTO rating (trip_id, commenter_id, comment)
SELECT (
    SELECT id
    FROM trip
    LIMIT 1
  ), (
    SELECT id
    FROM users
    WHERE email = 'jane.smith@example.com'
  ),
  'Great trip, very comfortable.';
-- Insert data into 'trip_station' table
INSERT INTO trip_station (trip_id, name)
SELECT (
    SELECT id
    FROM trip
    LIMIT 1
  ), 'Intermediate Stop';
-- Insert data into 'trip_passenger' table
INSERT INTO trip_passenger (
    trip_id,
    passenger_id,
    source_id,
    destination_id
  )
SELECT (
    SELECT id
    FROM trip
    LIMIT 1
  ), (
    SELECT id
    FROM users
    WHERE email = 'jane.smith@example.com'
  ),
  (
    SELECT id
    FROM locations
    WHERE name = 'Location 1'
  ),
  (
    SELECT id
    FROM locations
    WHERE name = 'Location 2'
  );
-- Insert data into 'alert' table
INSERT INTO alert (trip_id, passenger_id, comment)
SELECT (
    SELECT id
    FROM trip
    LIMIT 1
  ), (
    SELECT id
    FROM users
    WHERE email = 'jane.smith@example.com'
  ),
  'Left a bag in the car.';
-- Insert data into 'report' table
INSERT INTO report (trip_id, complainant, defendant, comment)
SELECT (
    SELECT id
    FROM trip
    LIMIT 1
  ), (
    SELECT id
    FROM users
    WHERE email = 'jane.smith@example.com'
  ),
  (
    SELECT id
    FROM users
    WHERE email = 'john.doe@example.com'
  ),
  'Driver was late.';
-- Insert data into 'favorite_driver' table
INSERT INTO favorite_driver (user_id, driver_id)
SELECT (
    SELECT id
    FROM users
    WHERE email = 'jane.smith@example.com'
  ),
  (
    SELECT user_id
    FROM driver
    WHERE license = 'DL123456'
  );