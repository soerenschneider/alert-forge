package report

import (
	"bytes"
	"cmp"
	"fmt"
	"html/template"
	"maps"
	"slices"

	"github.com/soerenschneider/alert-forge/internal"
	"github.com/soerenschneider/alert-forge/internal/model"
	"github.com/soerenschneider/alert-forge/internal/model/digest"
	"github.com/soerenschneider/alert-forge/internal/templating"
	"github.com/soerenschneider/alert-forge/pkg"
)

type DigestCreator struct {
	htmlTemplate *template.Template
}

func NewDigestCreator() (*DigestCreator, error) {
	templateData, err := templating.GetTemplate("digest.html")
	if err != nil {
		return nil, err
	}

	ret := &DigestCreator{}
	ret.htmlTemplate, err = template.New("digest").Funcs(template.FuncMap{
		"isZeroTime":  templating.IsZeroTime,
		"mapSeverity": templating.MapSeverityToColor,
	}).Parse(templateData)

	return ret, err
}

func (b *DigestCreator) Render(data digest.DigestData) ([]byte, error) {
	var buf bytes.Buffer
	err := b.htmlTemplate.Execute(&buf, data)
	return buf.Bytes(), err
}

func FormatSubjectLine(alerts []model.Alert) string {
	totalAlerts := len(alerts)
	severityCounts := make(map[string]int) // map to store the count for each severity

	for _, alert := range alerts {
		severity, found := alert.Labels["severity"]
		if found {
			severityCounts[severity]++
		} else {
			severityCounts["unspecified"]++
		}
	}

	var severities []string
	for severity := range severityCounts {
		severities = append(severities, severity)
	}

	slices.SortFunc(severities, func(a, b string) int {
		severityA := int(cmp.Or(internal.DefaultSeverityMapping[a], -1))
		severityB := int(cmp.Or(internal.DefaultSeverityMapping[b], -1))

		return severityB - severityA
	})

	severitySummary := ""
	for _, severity := range severities {
		if severitySummary != "" {
			severitySummary += ", "
		}
		severitySummary += fmt.Sprintf("%d %s", severityCounts[severity], severity)
	}

	return fmt.Sprintf("%d alerts total (%s)", totalAlerts, severitySummary)
}

func CompareStatusReports(previous, current digest.StatusReport) map[string]map[string]int {
	comparison := make(map[string]map[string]int)

	// Initialize result maps for keyNew, keyGone, and keySeen alerts per severity
	for severity := range pkg.Concat(maps.Keys(previous.SeverityCount), maps.Keys(current.SeverityCount)) {
		comparison[severity] = map[string]int{
			digest.KeyNew:  0,
			digest.KeyGone: 0,
			digest.KeySeen: 0,
		}
	}

	// Compare the fingerprints for new, gone, and seen alerts
	for severity, fingerprints := range previous.SeverityCount {
		for fingerprint := range fingerprints {
			_, found := current.SeverityCount[severity][fingerprint]
			if !found {
				comparison[severity][digest.KeyGone] += 1
			}
		}
	}

	for severity, fingerprints := range current.SeverityCount {
		for fingerprint := range fingerprints {
			_, found := previous.SeverityCount[severity][fingerprint]
			if found {
				comparison[severity][digest.KeySeen] += 1
			} else {
				comparison[severity][digest.KeyNew] += 1
			}
		}
	}

	return comparison
}
