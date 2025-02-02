package webhooks

import (
	"bytes"
	"html/template"

	"github.com/soerenschneider/alert-forge/internal/model"
	"github.com/soerenschneider/alert-forge/internal/report"
	"github.com/soerenschneider/alert-forge/internal/templating"
)

type AlertRenderer struct {
	htmlTemplate *template.Template
}

func NewAlertRenderer() (*AlertRenderer, error) {
	ret := &AlertRenderer{}
	templateData, err := templating.GetTemplate("alerts.html")
	if err != nil {
		return nil, err
	}
	ret.htmlTemplate, err = template.New("alertReport").Funcs(template.FuncMap{
		"mapSeverity": templating.MapSeverityToColor,
		"isZeroTime":  templating.IsZeroTime,
	}).Parse(templateData)
	return ret, err
}

type templateData struct {
	ActiveAlerts       []model.Alert
	ResolvedAlerts     []model.Alert
	SeveritiesActive   string
	SeveritiesResolved string
}

func (b *AlertRenderer) Render(alerts []model.Alert) ([]byte, error) {
	data := templateData{
		ActiveAlerts:   make([]model.Alert, 0, len(alerts)/2),
		ResolvedAlerts: make([]model.Alert, 0, len(alerts)/2),
	}

	for _, alert := range alerts {
		if alert.EndsAt.IsZero() {
			data.ActiveAlerts = append(data.ActiveAlerts, alert)
		} else {
			data.ResolvedAlerts = append(data.ResolvedAlerts, alert)
		}
	}

	data.SeveritiesActive = report.FormatSubjectLine(data.ActiveAlerts)
	data.SeveritiesResolved = report.FormatSubjectLine(data.ResolvedAlerts)

	var buf bytes.Buffer
	err := b.htmlTemplate.Execute(&buf, data)
	return buf.Bytes(), err
}
