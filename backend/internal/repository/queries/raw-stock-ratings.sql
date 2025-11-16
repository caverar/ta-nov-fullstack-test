-- name: AddRawStockRatings :copyfrom
INSERT INTO raw_stock_ratings (
    ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: DeleteRawStockRatings :exec
TRUNCATE TABLE raw_stock_ratings;