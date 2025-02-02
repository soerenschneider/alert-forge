-- name: AnalyzeLongestDuration :many
SELECT
    name,
    CAST((ended - started) AS INT) AS duration
FROM alerts
WHERE ended IS NOT NULL
ORDER BY duration DESC LIMIT 10;

-- name: AnalyzeInstances :many
SELECT
    instance,
    COUNT(*) AS instance_count
FROM alerts
GROUP BY instance
ORDER BY instance_count DESC LIMIT 10;

-- name: AnalyzeSeverities :many
SELECT
    CAST(json_extract(data,'$.labels.severity') AS VARCHAR) AS severity,
    COUNT(*) AS alert_count
FROM alerts
GROUP BY severity
ORDER BY severity;

-- name: AnalyzeAlertsPerDay :many
WITH days AS (
    SELECT DATE('now', '-6 days') AS day
UNION ALL
SELECT DATE(day, '+1 day') FROM days WHERE day < DATE('now')
    )
SELECT
    CAST(days.day AS varchar) AS alert_date,
    COUNT(alerts.id) AS alert_count
FROM days
LEFT JOIN alerts ON DATE(alerts.started, 'unixepoch') = days.day
GROUP BY days.day
ORDER BY days.day;

-- name: AnalyzeAverageAlertDuration :many
SELECT
    name,
    CAST(AVG(ended - started) AS INT) AS average_duration
FROM alerts
WHERE ended IS NOT NULL
GROUP BY name
HAVING COUNT(*) > 3
ORDER BY average_duration DESC;

-- name: AnalyzeAverageAlertDurationBySeverity :many
SELECT
    CAST(json_extract(data,'$.labels.severity') AS varchar) AS severity,
    CAST(AVG(ended - started) AS INT) AS average_duration
FROM alerts
WHERE ended IS NOT NULL
GROUP BY severity

UNION ALL

SELECT
    "$COMBINED" AS severity,
    CAST(AVG(ended - started) AS INT) AS average_duration
FROM alerts
WHERE ended IS NOT NULL
ORDER BY severity DESC;

