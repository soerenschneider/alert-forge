package webhooks

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/soerenschneider/alert-forge/internal/model"
	"github.com/soerenschneider/alert-forge/internal/templating"
)

type StatsRenderer struct {
	htmlTemplate *template.Template
}

func NewStatsRenderer() (*StatsRenderer, error) {
	ret := &StatsRenderer{}
	templateData, err := templating.GetTemplate("stats_duration.html")
	if err != nil {
		return nil, err
	}
	ret.htmlTemplate, err = template.New("statsDuration").Funcs(template.FuncMap{
		"mapSeverity": templating.MapSeverityToColor,
		"isZeroTime":  templating.IsZeroTime,
		"json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
	}).Parse(templateData)
	return ret, err
}

func (b *StatsRenderer) Render(stats model.AlertStats) ([]byte, error) {
	var buf bytes.Buffer
	err := b.htmlTemplate.Execute(&buf, stats)
	return buf.Bytes(), err
}
