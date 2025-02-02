package digest

import (
	"cmp"
	"time"

	"github.com/soerenschneider/alert-forge/internal/model"
)

const (
	KeyGone = "gone"
	KeySeen = "seen"
	KeyNew  = "new"
)

type DigestData struct {
	Severities     []SeverityInfo
	FiringAlerts   []model.Alert
	ResolvedAlerts []model.Alert
}

type SeverityInfo struct {
	Severity  string
	Total     int
	NewCount  int
	GoneCount int
	SeenCount int
}

// StatusReport holds the count of each severity in the active alerts and stores alert fingerprints.
type StatusReport struct {
	Id   int64
	Date time.Time

	// example
	// map[warning][fingerprint]{}
	SeverityCount map[string]map[string]struct{} `json:"severity_count"`
}

func (s *StatusReport) AlertCount() int {
	cnt := 0
	for _, fingerprints := range s.SeverityCount {
		cnt += len(fingerprints)
	}
	return cnt
}

// returns a map that contains information about how many alerts a new, gone and already seen compared to the

func GroupAlertsBySeverity(alerts []model.Alert) map[string]struct{} {
	groupedAlerts := make(map[string]struct{})

	for _, alert := range alerts {
		severity := alert.Labels["severity"]
		groupedAlerts[severity] = struct{}{}
	}

	return groupedAlerts
}

func Convert(i *SeverityInfo, severities map[string]map[string]int) {
	info := severities[i.Severity]

	i.GoneCount = cmp.Or(info[KeyGone], 0)
	i.SeenCount = cmp.Or(info[KeySeen], 0)
	i.NewCount = cmp.Or(info[KeyNew], 0)

	i.Total = i.SeenCount + i.NewCount
}
