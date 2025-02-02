-- name: GetReport :one
SELECT id, date, severity_count FROM reports WHERE id = sqlc.arg(id);

-- name: GetLatestReport :one
SELECT id, date, severity_count
FROM reports
ORDER BY date DESC
LIMIT 1;

-- name: SaveReport :exec
INSERT INTO reports (date, severity_count) VALUES (sqlc.arg(date), sqlc.arg(severity_count));

-- name: GetAlertsToday :many
SELECT data
FROM alerts
WHERE ended IS NULL OR (
    ended >= strftime('%s', 'now', 'start of day') AND
    ended < strftime('%s', 'now', 'start of day', '+1 day')
    )
ORDER BY severity DESC, started DESC, ended DESC;

-- name: GetAlertsYesterday :many
SELECT data
FROM alerts
WHERE ended >= strftime('%s', 'now', 'start of day', '-1 day') AND
      ended < strftime('%s', 'now', 'start of day')
ORDER BY severity DESC, started DESC, ended DESC;

-- name: GetAlertsTodayResolved :many
SELECT data
FROM alerts
WHERE ended >= strftime('%s', 'now', 'start of day') AND
      ended < strftime('%s', 'now', 'start of day', '+1 day')
ORDER BY severity DESC, started DESC, ended DESC;

-- name: GetAlert :one
SELECT data
FROM alerts
WHERE fingerprint = sqlc.arg(fingerprint);

-- name: GetAlerts :many
SELECT data
FROM alerts
ORDER BY severity DESC, started DESC, ended DESC;

-- name: GetActiveAlerts :many
SELECT data
FROM alerts
WHERE ended IS NULL;

-- name: GetAlertsByInstance :many
SELECT data
FROM alerts
WHERE instance = sqlc.arg(instance)
ORDER BY severity DESC, started DESC, ended DESC;

-- name: GetAlertsBySeverity :many
SELECT data
FROM alerts
WHERE severity = sqlc.arg(severity)
ORDER BY started DESC, ended DESC;

-- name: GetAlertsBetween :many
SELECT data
FROM alerts
WHERE started BETWEEN sqlc.arg(start) AND sqlc.arg(end)
ORDER BY severity DESC, started DESC, ended DESC;

-- name: SaveAlert :exec
INSERT INTO alerts (
    fingerprint,
    name,
    started,
    ended,
    severity,
    instance,
    data
)
VALUES (
    json_extract(sqlc.arg(data), '$.fingerprint'),
    sqlc.arg(name),
    sqlc.arg(started),
    sqlc.arg(ended),
    sqlc.arg(severity),
    sqlc.arg(instance),
    sqlc.arg(data)
)
ON CONFLICT(fingerprint) DO UPDATE SET
    ended=excluded.ended,
    data=excluded.data;
