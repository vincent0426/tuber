CREATE VIEW driver_view AS
SELECT users.*,
  driver.license,
  driver.verified,
  driver.brand,
  driver.model,
  driver.color,
  driver.plate,
  driver.created_at AS driver_created_at
FROM driver
  JOIN users ON users.id = driver.user_id;
CREATE VIEW trip_view AS
SELECT trip.id AS id,
  trip.passenger_limit AS passenger_limit,
  trip.status AS status,
  trip.start_time AS start_time,
  trip.created_at AS created_at,
  trip.updated_at AS updated_at,
  users.id AS driver_id,
  users.name AS driver_name,
  users.image_url AS driver_image_url,
  driver.brand AS driver_brand,
  driver.model AS driver_model,
  driver.color AS driver_color,
  driver.plate AS driver_plate,
  location_source.id AS source_id,
  location_source.name AS source_name,
  location_source.place_id AS source_place_id,
  location_source.lat_lon AS source_lat_lon,
  location_destination.id AS destination_id,
  location_destination.name AS destination_name,
  location_destination.place_id AS destination_place_id,
  location_destination.lat_lon AS destination_lat_lon
FROM trip
  JOIN users ON users.id = trip.driver_id
  JOIN driver ON trip.driver_id = driver.user_id
  JOIN locations AS location_source ON trip.source_id = location_source.id
  JOIN locations AS location_destination ON trip.destination_id = location_destination.id;
-- trip_passenger_view
CREATE VIEW trip_passenger_view AS
SELECT trip_passenger.trip_id AS trip_id,
  single_trip_passenger.id AS passenger_id,
  single_trip_passenger.name AS passenger_name,
  single_trip_passenger.image_url AS passenger_image_url,
  users.id AS driver_id,
  users.name AS driver_name,
  users.image_url AS driver_image_url,
  driver.brand AS driver_brand,
  driver.model AS driver_model,
  driver.color AS driver_color,
  driver.plate AS driver_plate,
  location_source.name AS source_name,
  location_source.place_id AS source_place_id,
  location_source.lat_lon AS source_lat_lon,
  location_destination.name AS destination_name,
  location_destination.place_id AS destination_place_id,
  location_destination.lat_lon AS destination_lat_lon,
  passenger_location_source.name AS passenger_location_source_name,
  passenger_location_source.place_id AS passenger_location_source_place_id,
  passenger_location_source.lat_lon AS passenger_location_source_lat_lon,
  passenger_location_destination.name AS passenger_location_destination_name,
  passenger_location_destination.place_id AS passenger_location_destination_place_id,
  passenger_location_destination.lat_lon AS passenger_location_destination_lat_lon
FROM trip_passenger
  JOIN trip ON trip_passenger.trip_id = trip.id
  JOIN users AS single_trip_passenger ON single_trip_passenger.id = trip_passenger.passenger_id
  JOIN users ON users.id = trip.driver_id
  JOIN driver ON trip.driver_id = driver.user_id
  JOIN locations AS location_source ON trip.source_id = location_source.id
  JOIN locations AS location_destination ON trip.destination_id = location_destination.id
  JOIN locations AS passenger_location_source ON trip_passenger.source_id = passenger_location_source.id
  JOIN locations AS passenger_location_destination ON trip_passenger.destination_id = passenger_location_destination.id;
CREATE VIEW favorite_driver_view AS
SELECT favorite_driver.*,
  users.name AS driver_name,
  users.image_url AS driver_image_url,
  driver.brand AS driver_brand,
  driver.model AS driver_model,
  driver.color AS driver_color,
  driver.plate AS driver_plate,
  driver.created_at AS driver_created_at
FROM favorite_driver
  JOIN users ON users.id = favorite_driver.driver_id
  JOIN driver ON favorite_driver.driver_id = driver.user_id;