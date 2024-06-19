CREATE EXTENSION IF NOT EXISTS "uuid-ossp"

CREATE TABLE categories
(
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(10)         NOT NULL
)

CREATE TABLE operations
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID           NOT NULL,
    money_sum   NUMERIC(15, 2) NOT NULL,
    description TEXT,
    date_time   TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT category_fk FOREIGN KEY (category_id) REFERENCES categories (id)
);