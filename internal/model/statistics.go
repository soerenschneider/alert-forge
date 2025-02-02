package model

type AlertStats struct {
	AlertDuration                  []AlertDuration                  `json:"alert_duration"`
	AlertsByDay                    []AlertsByDay                    `json:"alerts_by_day"`
	AverageAlertDurationBySeverity []AverageAlertDurationBySeverity `json:"avg_alert_duration_by_severity"`
	AlertsBySeverity               []AlertsBySeverity               `json:"alerts_by_severity"`
	AlertsByInstance               []AlertsByInstance               `json:"alerts_by_instance"`
}

type AlertDuration struct {
	Name            string `json:"name"`
	AverageDuration int64  `json:"average_duration"`
}

type AlertsByDay struct {
	AlertDate  string `json:"alert_date"`
	AlertCount int64  `json:"alert_count"`
}

type AverageAlertDurationBySeverity struct {
	Severity        string `json:"severity"`
	AverageDuration int64  `json:"average_duration"`
}

type AlertsBySeverity struct {
	Severity   string `json:"severity"`
	AlertCount int64  `json:"alert_count"`
}

type AlertsByInstance struct {
	Instance      string `json:"instance"`
	InstanceCount int64  `json:"instance_count"`
}
