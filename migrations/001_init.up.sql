-- Create enum types for fields with limited possible values
CREATE TYPE order_locale AS ENUM ('en', 'ru', 'de', 'fr', 'es', 'zh');
CREATE TYPE payment_provider AS ENUM ('wbpay', 'paypal', 'stripe', 'other');

-- Create delivery table
CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);

-- Create payment table
CREATE TABLE payment (
    transaction VARCHAR(50) PRIMARY KEY,
    request_id VARCHAR(50),
    currency VARCHAR(3) NOT NULL,
    provider payment_provider NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL DEFAULT 0
);

-- Create order table
CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(10) NOT NULL,
    delivery_id INTEGER REFERENCES delivery(id) NOT NULL,
    payment_transaction VARCHAR(50) REFERENCES payment(transaction) NOT NULL,
    locale order_locale NOT NULL,
    internal_signature VARCHAR(100),
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

-- Create item table
CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) NOT NULL,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INTEGER NOT NULL  
);

-- Create indexes for better performance
CREATE INDEX idx_item_order_uid ON item(order_uid);
CREATE INDEX idx_orders_track_number ON orders(track_number);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_date_created ON orders(date_created);