CREATE TABLE IF NOT EXISTS stock_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker TEXT NOT NULL,
    company TEXT NOT NULL,
    target_from NUMERIC NOT NULL,
    target_to NUMERIC NOT NULL,
    action TEXT NOT NULL,
    rating_from TEXT NOT NULL,
    rating_to TEXT NOT NULL,
    at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);