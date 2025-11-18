CREATE TYPE STOCK_RATING_TYPE AS ENUM ( 'buy', 'hold', 'sell', 'pending');
CREATE TYPE STOCK_ACTION_TYPE AS ENUM ( 'up', 'down', 'reiterated');
CREATE TABLE IF NOT EXISTS stock_rating (
    ticker TEXT PRIMARY KEY NOT NULL,
    company TEXT NOT NULL,
    brokerage TEXT NOT NULL,
    target_from NUMERIC(10,2) NOT NULL,
    target_to NUMERIC(10,2) NOT NULL,
    action STOCK_ACTION_TYPE NOT NULL,
    raw_action TEXT NOT NULL,
    rating_from STOCK_RATING_TYPE NOT NULL,
    raw_rating_from TEXT NOT NULL,
    rating_to STOCK_RATING_TYPE NOT NULL,
    raw_rating_to TEXT NOT NULL,
    at TIMESTAMPTZ NOT NULL
);