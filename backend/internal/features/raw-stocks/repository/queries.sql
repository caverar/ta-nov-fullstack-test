-- name: AddStickEvents :exec
INSERT INTO stock_events (
    ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);