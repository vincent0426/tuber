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
SELECT trip.*,
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
  location_destination.lat_lon AS destination_lat_lon
FROM trip
  JOIN users ON users.id = trip.driver_id
  JOIN driver ON trip.driver_id = driver.user_id
  JOIN locations AS location_source ON trip.source_id = location_source.id
  JOIN locations AS location_destination ON trip.destination_id = location_destination.id;
-- trip_passenger_view
CREATE VIEW trip_passenger_view AS
SELECT trip_passenger.*,
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
  JOIN users ON users.id = trip.driver_id
  JOIN driver ON trip.driver_id = driver.user_id
  JOIN locations AS location_source ON trip.source_id = location_source.id
  JOIN locations AS location_destination ON trip.destination_id = location_destination.id
  JOIN locations AS passenger_location_source ON trip_passenger.source_id = passenger_location_source.id
  JOIN locations AS passenger_location_destination ON trip_passenger.destination_id = passenger_location_destination.id;