-- Create orders table
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(10) NOT NULL,  
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(50) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard VARCHAR(50) NOT NULL
);

-- Create delivery table
CREATE TABLE delivery (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

-- Create payment table
CREATE TABLE payment (
    transaction VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

-- Create items table
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(50) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL
);

-- Create index for better performance on order_uid in items table
CREATE INDEX idx_items_order_uid ON items(order_uid);

-- Create inbox table (for inbox kafka messages)
CREATE TABLE IF NOT EXISTS inbox (
    message_id TEXT PRIMARY KEY,
    topic TEXT NOT NULL,
    payload TEXT NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);