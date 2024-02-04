CREATE TABLE Orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entryD VARCHAR(50),
    locale VARCHAR(50),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(50),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(50)
);

CREATE TABLE Delivery (
    delivery_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255),
    nameD VARCHAR(255),
    phone VARCHAR(50),
    zip VARCHAR(50),
    city VARCHAR(100),
    addressD VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(100),
    FOREIGN KEY (order_uid) REFERENCES Orders(order_uid)
);

CREATE TABLE Payment (
    payment_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255),
    transactionD VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(255),
    providerD VARCHAR(255),
    amount INT,
    payment_dt INT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    FOREIGN KEY(order_uid) REFERENCES Orders(order_uid)
);

CREATE TABLE Items (
    item_id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255),
    chrt_id INT,
    track_number VARCHAR(255),
    price FLOAT,
    rid VARCHAR(255),
    nameD VARCHAR(255),
    sale INT,
    sizeD VARCHAR(255),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    statusD INT,
    FOREIGN KEY(order_uid) REFERENCES Orders(order_uid)
);