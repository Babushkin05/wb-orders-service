-- Drop tables in reverse order of creation due to foreign key constraints
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS delivery;

-- Drop custom types
DROP TYPE IF EXISTS payment_provider;
DROP TYPE IF EXISTS order_locale;