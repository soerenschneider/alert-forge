package report

import (
	"reflect"
	"testing"

	"github.com/soerenschneider/alert-forge/internal/model/digest"
)

func TestCompareStatusReports(t *testing.T) {
	type args struct {
		previous digest.StatusReport
		current  digest.StatusReport
	}
	tests := []struct {
		name string
		args args
		want map[string]map[string]int
	}{
		{
			name: "one new, one seen, one gone",
			args: args{
				previous: digest.StatusReport{
					SeverityCount: map[string]map[string]struct{}{
						"warning": {
							"fingerprint-1": {},
							"fingerprint-2": {},
						},
					},
				},
				current: digest.StatusReport{
					SeverityCount: map[string]map[string]struct{}{
						"warning": {
							"fingerprint-2": {},
							"fingerprint-3": {},
						},
					},
				},
			},
			want: map[string]map[string]int{
				"warning": {
					digest.KeyGone: 1,
					digest.KeyNew:  1,
					digest.KeySeen: 1,
				},
			},
		},
		{
			name: "both seen",
			args: args{
				previous: digest.StatusReport{
					SeverityCount: map[string]map[string]struct{}{
						"warning": {
							"fingerprint-1": {},
							"fingerprint-2": {},
						},
					},
				},
				current: digest.StatusReport{
					SeverityCount: map[string]map[string]struct{}{
						"warning": {
							"fingerprint-1": {},
							"fingerprint-2": {},
						},
					},
				},
			},
			want: map[string]map[string]int{
				"warning": {
					digest.KeyGone: 0,
					digest.KeyNew:  0,
					digest.KeySeen: 2,
				},
			},
		},
		{
			name: "one gone, one seen",
			args: args{
				previous: digest.StatusReport{
					SeverityCount: map[string]map[string]struct{}{
						"warning": {
							"fingerprint-1": {},
							"fingerprint-2": {},
						},
					},
				},
				current: digest.StatusReport{
					SeverityCount: map[string]map[string]struct{}{
						"warning": {
							"fingerprint-1": {},
						},
					},
				},
			},
			want: map[string]map[string]int{
				"warning": {
					digest.KeyGone: 1,
					digest.KeyNew:  0,
					digest.KeySeen: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareStatusReports(tt.args.previous, tt.args.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareStatusReports() = %v, want %v", got, tt.want)
			}
		})
	}
}
