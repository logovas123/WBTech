CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(10),
    internal_signature VARCHAR(100),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created VARCHAR(50),
    oof_shard VARCHAR(10)
);

CREATE TABLE delivery (
    order_uid VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(100),
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE payment (
    transaction_id VARCHAR(50) PRIMARY KEY,
    order_uid VARCHAR(50) NOT NULL,
    request_id VARCHAR(50),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount NUMERIC(10, 2),
    payment_dt BIGINT,
    bank VARCHAR(50),
    delivery_cost NUMERIC(10, 2),
    goods_total NUMERIC(10, 2),
    custom_fee NUMERIC(10, 2),
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) NOT NULL,
    chrt_id BIGINT,
    track_number VARCHAR(50),
    price NUMERIC(10, 2),
    rid VARCHAR(50),
    name VARCHAR(255),
    sale NUMERIC(10, 2),
    size VARCHAR(50),
    total_price NUMERIC(10, 2),
    nm_id BIGINT,
    brand VARCHAR(100),
    status INT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

