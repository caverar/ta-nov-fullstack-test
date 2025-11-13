CREATE TABLE IF NOT EXISTS stock_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker TEXT NOT NULL,
    target_from TEXT NOT NULL,
    target_to TEXT NOT NULL,
    company TEXT NOT NULL,
    action TEXT NOT NULL,
    brokerage TEXT NOT NULL,
    rating_from TEXT NOT NULL,
    rating_to TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);