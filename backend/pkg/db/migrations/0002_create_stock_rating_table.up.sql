CREATE TYPE STOCK_RATING_TYPE AS ENUM ( 'buy', 'hold', 'sell', 'pending');
CREATE TYPE STOCK_ACTION_TYPE AS ENUM ( 'up', 'down', 'reiterated');
CREATE TABLE IF NOT EXISTS stock_rating (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker TEXT NOT NULL,
    company TEXT NOT NULL,
    target_from NUMERIC(10,2) NOT NULL,
    target_to NUMERIC(10,2) NOT NULL,
    action STOCK_ACTION_TYPE NOT NULL,
    rating_from STOCK_RATING_TYPE NOT NULL,
    rating_to STOCK_RATING_TYPE NOT NULL,
    at TIMESTAMPTZ NOT NULL
);