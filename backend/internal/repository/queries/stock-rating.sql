-- name: AddStockRatings :copyfrom
INSERT INTO stock_rating (
    ticker, company, brokerage, target_from, target_to, action, raw_action, rating_from, raw_rating_from, rating_to, raw_rating_to, at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: ClearStockRating :exec
TRUNCATE TABLE stock_rating;


-- List
-- name: GetStockRatings :many
WITH scored_stock_ratings AS (
    SELECT
        ticker,
        company,
        brokerage,
        target_from,
        target_to,
        action,
        raw_action,
        rating_from,
        rating_to,
        at,
        (target_to - target_from)::Numeric(10,2) AS target_delta,
        (TRUNC((10 * COALESCE( (target_to - target_from) / target_from, 0))
        + (2 * (CASE rating_to
            WHEN 'buy' THEN 1
            WHEN 'hold' THEN 0
            WHEN 'pending' THEN 0
            WHEN 'sell' THEN -1
        END))
        + (1 * (CASE action
            WHEN 'up' THEN 1
            WHEN 'down' THEN -1
            WHEN 'reiterated' THEN 0
        END)), 3)*1000)::DECIMAL "score"
    FROM stock_rating
    WHERE
        (sqlc.arg('ticker_like')::text IS NULL OR ticker ILIKE '%' || sqlc.arg('ticker_like')::text || '%')
        AND (sqlc.arg('company_like')::text IS NULL OR company ILIKE '%' || sqlc.arg('company_like')::text || '%')
)
SELECT
    ticker,
    company,
    brokerage,
    target_from::text,
    target_to::text,
    action,
    raw_action,
    rating_from,
    rating_to,
    at,
    target_delta::text,
    score::INTEGER
FROM scored_stock_ratings
ORDER BY
    -- Numeric ordering
    CASE WHEN sqlc.arg('sort_order')::text = 'desc' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'target_from' THEN target_from
            WHEN 'target_to' THEN target_to
            WHEN 'target_delta' THEN target_delta
            WHEN 'score' THEN score
            ELSE NULl
        END
    END DESC,
    CASE WHEN sqlc.arg('sort_order')::text = 'asc' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'target_from' THEN target_from
            WHEN 'target_to' THEN target_to
            WHEN 'target_delta' THEN target_delta
            WHEN 'score' THEN score
            ELSE NULl
        END
    END ASC,
    -- String Ordering
    CASE WHEN sqlc.arg('sort_order')::text = 'desc' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'ticker' THEN ticker::text
            WHEN 'company' THEN company::text
            WHEN 'brokerage' THEN brokerage::text
            WHEN 'action' THEN action::text
            WHEN 'rating_from' THEN rating_from::text
            WHEN 'rating_to' THEN rating_to::text
            ELSE NULl
        END
    END DESC,
    CASE WHEN sqlc.arg('sort_order')::text = 'asc' THEN
        CASE sqlc.arg('sort_by')::text
            WHEN 'ticker' THEN ticker::text
            WHEN 'company' THEN company::text
            WHEN 'brokerage' THEN brokerage::text
            WHEN 'action' THEN action::text
            WHEN 'rating_from' THEN rating_from::text
            WHEN 'rating_to' THEN rating_to::text
            ELSE NULl
        END
    END,
    ticker ASC

LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- Recommendations Dashboard

-- name: GetOverallMarketStockRatings :many
SELECT
    rating_to AS rating,
    COUNT(*) AS count
FROM stock_rating
GROUP BY action
ORDER BY count DESC;

-- name: GetOverallAnalystActions :many
SELECT
    action AS action,
    COUNT(*) AS count
FROM stock_rating
GROUP BY action
ORDER BY count DESC;

-- name GetStockRating
SELECT * FROM stock_rating WHERE ticker = $1;

