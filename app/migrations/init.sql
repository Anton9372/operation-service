CREATE SCHEMA IF NOT EXISTS public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.categories
(
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID         NOT NULL,
    name    VARCHAR(100) NOT NULL,
    type    VARCHAR(10)  NOT NULL
);

CREATE TABLE public.operations
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID           NOT NULL,
    money_sum   NUMERIC(15, 2) NOT NULL,
    description VARCHAR(255),
    date_time   TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT category_fk FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);