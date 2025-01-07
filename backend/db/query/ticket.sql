-- name: GetRowTickets :many
SELECT * FROM tickets
WHERE event_id = $1
AND zone_id = $2
AND row = $3;
