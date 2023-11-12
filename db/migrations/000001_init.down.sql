-- Start with tables that do not reference other tables or are only referenced by other tables
DROP TABLE IF EXISTS location CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- Next, drop tables that reference the above tables
DROP TABLE IF EXISTS car CASCADE;
DROP TABLE IF EXISTS driver CASCADE;
DROP TABLE IF EXISTS trip CASCADE;
DROP TABLE IF EXISTS chat_history CASCADE;
DROP TABLE IF EXISTS rating CASCADE;
DROP TABLE IF EXISTS trip_station CASCADE;
DROP TABLE IF EXISTS trip_passenger CASCADE;
DROP TABLE IF EXISTS alert CASCADE;
DROP TABLE IF EXISTS report CASCADE;
DROP TABLE IF EXISTS favorite_driver CASCADE;
-- Drop extensions last
DROP EXTENSION IF EXISTS "postgis";
DROP EXTENSION IF EXISTS "uuid-ossp";