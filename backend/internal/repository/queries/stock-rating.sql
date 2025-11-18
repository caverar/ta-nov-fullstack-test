-- name: AddStockRatings :copyfrom
INSERT INTO stock_rating (
    ticker, company, brokerage, target_from, target_to, action, raw_action, rating_from, raw_rating_from, rating_to, raw_rating_to, at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: ClearStockRating :exec
TRUNCATE TABLE stock_rating;


-- name: GetDetailStockRatingList :many
SELECT ticker, company, target_from, target_to, action, rating_from, rating_to, at
FROM stock_rating
LIMIT $1
OFFSET $2;

-- name: GetStockRatingList :many
SELECT ticker, company
FROM stock_rating
ORDER BY ticker ASC;

