-- name: AddStockRatings :copyfrom
INSERT INTO stock_rating (
    ticker, company, target_from, target_to, action, rating_from, rating_to, at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: ClearStockRating :exec
TRUNCATE TABLE stock_rating;