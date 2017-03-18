CREATE TABLE IF NOT EXISTS products
(
    product_id serial NOT NULL,
    name varchar(255) NOT NULL,
    created date NOT NULL,
    CONSTRAINT product_id_pkey PRIMARY KEY (product_id)
);

CREATE INDEX products_product_id_idx ON products (product_id);