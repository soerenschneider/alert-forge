// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: reports.sql

package generated

import (
	"context"
)

const analyzeAlertsPerDay = `-- name: AnalyzeAlertsPerDay :many
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
ORDER BY days.day
`

type AnalyzeAlertsPerDayRow struct {
	AlertDate  string `json:"alert_date"`
	AlertCount int64  `json:"alert_count"`
}

func (q *Queries) AnalyzeAlertsPerDay(ctx context.Context) ([]AnalyzeAlertsPerDayRow, error) {
	rows, err := q.db.QueryContext(ctx, analyzeAlertsPerDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnalyzeAlertsPerDayRow
	for rows.Next() {
		var i AnalyzeAlertsPerDayRow
		if err := rows.Scan(&i.AlertDate, &i.AlertCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const analyzeAverageAlertDuration = `-- name: AnalyzeAverageAlertDuration :many
SELECT
    name,
    CAST(AVG(ended - started) AS INT) AS average_duration
FROM alerts
WHERE ended IS NOT NULL
GROUP BY name
HAVING COUNT(*) > 3
ORDER BY average_duration DESC
`

type AnalyzeAverageAlertDurationRow struct {
	Name            string `json:"name"`
	AverageDuration int64  `json:"average_duration"`
}

func (q *Queries) AnalyzeAverageAlertDuration(ctx context.Context) ([]AnalyzeAverageAlertDurationRow, error) {
	rows, err := q.db.QueryContext(ctx, analyzeAverageAlertDuration)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnalyzeAverageAlertDurationRow
	for rows.Next() {
		var i AnalyzeAverageAlertDurationRow
		if err := rows.Scan(&i.Name, &i.AverageDuration); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const analyzeAverageAlertDurationBySeverity = `-- name: AnalyzeAverageAlertDurationBySeverity :many
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
ORDER BY severity DESC
`

type AnalyzeAverageAlertDurationBySeverityRow struct {
	Severity        string `json:"severity"`
	AverageDuration int64  `json:"average_duration"`
}

func (q *Queries) AnalyzeAverageAlertDurationBySeverity(ctx context.Context) ([]AnalyzeAverageAlertDurationBySeverityRow, error) {
	rows, err := q.db.QueryContext(ctx, analyzeAverageAlertDurationBySeverity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnalyzeAverageAlertDurationBySeverityRow
	for rows.Next() {
		var i AnalyzeAverageAlertDurationBySeverityRow
		if err := rows.Scan(&i.Severity, &i.AverageDuration); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const analyzeInstances = `-- name: AnalyzeInstances :many
SELECT
    instance,
    COUNT(*) AS instance_count
FROM alerts
GROUP BY instance
ORDER BY instance_count DESC LIMIT 10
`

type AnalyzeInstancesRow struct {
	Instance      string `json:"instance"`
	InstanceCount int64  `json:"instance_count"`
}

func (q *Queries) AnalyzeInstances(ctx context.Context) ([]AnalyzeInstancesRow, error) {
	rows, err := q.db.QueryContext(ctx, analyzeInstances)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnalyzeInstancesRow
	for rows.Next() {
		var i AnalyzeInstancesRow
		if err := rows.Scan(&i.Instance, &i.InstanceCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const analyzeLongestDuration = `-- name: AnalyzeLongestDuration :many
SELECT
    name,
    CAST((ended - started) AS INT) AS duration
FROM alerts
WHERE ended IS NOT NULL
ORDER BY duration DESC LIMIT 10
`

type AnalyzeLongestDurationRow struct {
	Name     string `json:"name"`
	Duration int64  `json:"duration"`
}

func (q *Queries) AnalyzeLongestDuration(ctx context.Context) ([]AnalyzeLongestDurationRow, error) {
	rows, err := q.db.QueryContext(ctx, analyzeLongestDuration)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnalyzeLongestDurationRow
	for rows.Next() {
		var i AnalyzeLongestDurationRow
		if err := rows.Scan(&i.Name, &i.Duration); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const analyzeSeverities = `-- name: AnalyzeSeverities :many
SELECT
    CAST(json_extract(data,'$.labels.severity') AS VARCHAR) AS severity,
    COUNT(*) AS alert_count
FROM alerts
GROUP BY severity
ORDER BY severity
`

type AnalyzeSeveritiesRow struct {
	Severity   string `json:"severity"`
	AlertCount int64  `json:"alert_count"`
}

func (q *Queries) AnalyzeSeverities(ctx context.Context) ([]AnalyzeSeveritiesRow, error) {
	rows, err := q.db.QueryContext(ctx, analyzeSeverities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnalyzeSeveritiesRow
	for rows.Next() {
		var i AnalyzeSeveritiesRow
		if err := rows.Scan(&i.Severity, &i.AlertCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
