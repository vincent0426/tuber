-- Seeding users Table
INSERT INTO users (name, email, bio)
VALUES (
    'John Doe',
    'johndoe@example.com',
    'Loves long drives'
  ),
  (
    'Jane Smith',
    'janesmith@example.com',
    'Enthusiast of eco-friendly travel'
  ),
  (
    'Alice Johnson',
    'alicejohnson@example.com',
    'Frequent traveler for business'
  ),
  (
    'Bob Brown',
    'bobbrown@example.com',
    'Enjoy exploring new places'
  );
-- Seeding car Table (assuming John Doe is a driver)
INSERT INTO car (driver_id, license)
VALUES (
    (
      SELECT id
      FROM users
      WHERE name = 'John Doe'
    ),
    'ABC123'
  ),
  (
    (
      SELECT id
      FROM users
      WHERE name = 'Jane Smith'
    ),
    'XYZ789'
  );
-- Seeding driver Table (assuming John Doe is a driver)
INSERT INTO driver (user_id, car_id)
VALUES (
    (
      SELECT id
      FROM users
      WHERE name = 'John Doe'
    ),
    (
      SELECT id
      FROM car
      WHERE driver_id = (
          SELECT id
          FROM users
          WHERE name = 'John Doe'
        )
    )
  ),
  (
    (
      SELECT id
      FROM users
      WHERE name = 'Jane Smith'
    ),
    (
      SELECT id
      FROM car
      WHERE driver_id = (
          SELECT id
          FROM users
          WHERE name = 'Jane Smith'
        )
    )
  );
-- Seeding location Table
INSERT INTO location (name, address, coordinates)
VALUES (
    'Central Park',
    '5th Ave, New York, NY',
    ST_GeomFromText('POINT(-73.965355 40.782865)', 4326)
  ),
  (
    'Golden Gate Bridge',
    'San Francisco, CA',
    ST_GeomFromText('POINT(-122.478255 37.819929)', 4326)
  );
-- Seeding trip Table (assuming John Doe started a trip)
INSERT INTO trip (
    driver_id,
    passenger_limit,
    source_id,
    destination_id,
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
      WHERE name = 'Central Park'
    ),
    (
      SELECT id
      FROM location
      WHERE name = 'Golden Gate Bridge'
    ),
    '2023-11-15 08:00:00'
  );
-- Seeding chat_history Table (sample chat between John and Jane)
INSERT INTO chat_history (trip_id, sender_id, msg_content)
VALUES (
    (
      SELECT id
      FROM trip
      LIMIT 1
    ), (
      SELECT id
      FROM users
      WHERE name = 'Alice Johnson'
    ),
    'Is there space for one more?'
  );
-- Seeding rating Table (Jane rates John's trip)
INSERT INTO rating (trip_id, commenter_id, comment)
VALUES (
    (
      SELECT id
      FROM trip
      LIMIT 1
    ), (
      SELECT id
      FROM users
      WHERE name = 'Alice Johnson'
    ),
    'Great trip, comfortable car!'
  );
-- Seeding trip_Station Table (One station for John's trip)
INSERT INTO trip_station (trip_id, name)
VALUES (
    (
      SELECT id
      FROM trip
      LIMIT 1
    ), 'Midway Point'
  );
-- Seeding trip_passenger Table (Jane joins John's trip)
INSERT INTO trip_passenger (
    trip_id,
    passenger_id,
    station_source_id,
    station_destination_id
  )
VALUES (
    (
      SELECT id
      FROM trip
      LIMIT 1
    ), (
      SELECT id
      FROM users
      WHERE name = 'Alice Johnson'
    ),
    (
      SELECT id
      FROM location
      WHERE name = 'Central Park'
    ),
    (
      SELECT id
      FROM location
      WHERE name = 'Golden Gate Bridge'
    )
  );
-- Seeding alert Table (Jane alerts about something in John's trip)
INSERT INTO alert (trip_id, passenger_id, comment)
VALUES (
    (
      SELECT id
      FROM trip
      LIMIT 1
    ), (
      SELECT id
      FROM users
      WHERE name = 'Alice Johnson'
    ),
    'Left my bag in the car'
  );
-- Seeding report Table (Jane reports John during a trip)
INSERT INTO report (trip_id, complainant, defendant, comment)
VALUES (
    (
      SELECT id
      FROM trip
      LIMIT 1
    ), (
      SELECT id
      FROM users
      WHERE name = 'Alice Johnson'
    ),
    (
      SELECT id
      FROM users
      WHERE name = 'John Doe'
    ),
    'Driver was very courteous'
  );
-- Seeding Favorite_driver Table (Jane adds John to her favorite drivers)
INSERT INTO favorite_driver (user_id, driver_id, note)
VALUES (
    (
      SELECT id
      FROM users
      WHERE name = 'Alice Johnson'
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
    'Always on time!'
  );