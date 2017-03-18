CREATE TYPE order_status AS ENUM ('created', 'processed', 'canceled', 'finished');

CREATE TABLE IF NOT EXISTS orders
(
    order_id serial NOT NULL,
    shop_id integer NOT NULL,
    customer_id integer NOT NULL,
    status order_status NOT NULL,
    created date NOT NULL,
    CONSTRAINT order_id_pkey PRIMARY KEY (order_id)
);

CREATE INDEX orders_shop_id_idx ON orders (shop_id);
CREATE INDEX customer_id_idx ON orders (customer_id);
CREATE INDEX status_idx ON orders (status);
CREATE INDEX created_idx ON orders (created);
