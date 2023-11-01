-- Seeding users Table
INSERT INTO users (name, email, bio, language)
VALUES (
    'John Doe',
    'john.doe@email.com',
    'Just a regular John.',
    'en'
  ),
  (
    'Jane Smith',
    'jane.smith@email.com',
    'Love to travel!',
    'en'
  );
-- Seeding car Table (assuming John Doe is a driver)
INSERT INTO car (driver_id, license)
SELECT id,
  'XYZ-1234'
FROM users
WHERE name = 'John Doe';
-- Seeding driver Table (assuming John Doe is a driver)
INSERT INTO driver (user_id, car_id)
SELECT users.id,
  car.id
FROM users,
  car
WHERE users.name = 'John Doe'
  AND car.driver_id = users.id;
-- Seeding location Table
INSERT INTO location (name, address, coordinates)
VALUES (
    'Central Station',
    '89 E 42nd St, New York, NY 10017',
    ST_GeomFromText('POINT(-73.977622 40.752726)', 4326)
  ),
  (
    'Central Park',
    '14 E 60th St, New York, NY 10022',
    ST_GeomFromText('POINT(-73.9712 40.7648)', 4326)
  );
-- Sample coordinates for New York's Central Park
-- Seeding trip Table (assuming John Doe started a trip)
INSERT INTO trip (
    driver_id,
    passenger_limit,
    source_id,
    destination_id,
    status,
    start_time
  )
VALUES (
    (
      SELECT user_id
      FROM driver
      WHERE user_id = (
          SELECT id
          FROM users
          WHERE name = 'John Doe'
        )
    ),
    3,
    (
      SELECT id
      FROM location
      WHERE name = 'Central Station'
    ),
    (
      SELECT id
      FROM location
      WHERE name = 'Central Park'
    ),
    'not_start',
    '2020-01-01 00:00:00'
  );
-- Seeding chat_history Table (sample chat between John and Jane)
INSERT INTO chat_history (trip_id, sender_id, msg_content)
SELECT trip.id,
  users.id,
  'Hey John, are you near?'
FROM trip,
  users
WHERE users.name = 'Jane Smith';
-- Seeding rating Table (Jane rates John's trip)
INSERT INTO rating (trip_id, commenter_id, comment)
SELECT trip.id,
  users.id,
  'Great trip!'
FROM trip,
  users
WHERE users.name = 'Jane Smith';
-- Seeding trip_Station Table (One station for John's trip)
INSERT INTO trip_Station (trip_id, name)
SELECT trip.id,
  'Central Station'
FROM trip
WHERE driver_id = (
    SELECT user_id
    FROM driver
    WHERE user_id = (
        SELECT id
        FROM users
        WHERE name = 'John Doe'
      )
  );
INSERT INTO trip_Station (trip_id, name)
SELECT trip.id,
  'Central Park'
FROM trip
WHERE driver_id = (
    SELECT user_id
    FROM driver
    WHERE user_id = (
        SELECT id
        FROM users
        WHERE name = 'John Doe'
      )
  );
-- Seeding trip_passenger Table (Jane joins John's trip)
INSERT INTO trip_passenger (
    trip_id,
    passenger_id,
    station_source_id,
    station_destination_id,
    status
  )
SELECT trip.id,
  users.id,
  (
    SELECT id
    FROM location
    WHERE name = 'Central Station'
  ),
  (
    SELECT id
    FROM location
    WHERE name = 'Central Park'
  ),
  'accepted'
FROM trip,
  users
WHERE users.name = 'Jane Smith';
-- Seeding alert Table (Jane alerts about something in John's trip)
INSERT INTO alert (trip_id, passenger_id, comment)
SELECT trip.id,
  users.id,
  'driver is driving a bit fast.'
FROM trip,
  users
WHERE users.name = 'Jane Smith';
-- Seeding report Table (Jane reports John during a trip)
INSERT INTO report (trip_id, complainant, defendant, comment)
VALUES (
    (
      SELECT id
      FROM trip
      WHERE driver_id = (
          SELECT user_id
          FROM driver
          WHERE user_id = (
              SELECT id
              FROM users
              WHERE name = 'John Doe'
            )
        )
    ),
    (
      SELECT id
      FROM users
      WHERE name = 'Jane Smith'
    ),
    (
      SELECT id
      FROM users
      WHERE name = 'John Doe'
    ),
    'Rash Driving'
  );
-- Seeding Favorite_driver Table (Jane adds John to her favorite drivers)
INSERT INTO favorite_driver (user_id, driver_id, note)
VALUES (
    (
      SELECT id
      FROM users
      WHERE name = 'Jane Smith'
    ),
    (
      SELECT user_id
      FROM driver
      WHERE user_id = (
          SELECT id
          FROM users
          WHERE name = 'John Doe'
        )
    ),
    'Always on time.'
  );