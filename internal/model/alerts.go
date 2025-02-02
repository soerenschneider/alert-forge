package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type Status struct {
	InhibitedBy []string `json:"inhibitedBy"` // Array of inhibited by
	SilencedBy  []string `json:"silencedBy"`  // Array of silenced by
	State       string   `json:"state"`       // The state (active, etc.)
}

type Alert struct {
	Annotations map[string]string `json:"annotations"` // Annotations field
	EndsAt      time.Time         `json:"endsAt"`      // End time
	Fingerprint string            `json:"fingerprint"` // Fingerprint
	Receivers   []struct {
		Name string `json:"name"` // Receiver name
	} `json:"receivers"` // Receivers array
	StartsAt     time.Time         `json:"startsAt"`     // Start time
	Status       Status            `json:"status"`       // Status is now a structured field
	UpdatedAt    string            `json:"updatedAt"`    // Updated time
	GeneratorURL string            `json:"generatorURL"` // URL for the alert generator
	Labels       map[string]string `json:"labels"`       // Labels field
}

func (a *Alert) IsActive() bool {
	return a.EndsAt.IsZero() || time.Now().After(a.EndsAt)
}

func (a *Alert) Hash() string {
	return Hash(a.StartsAt.String(), a.Labels)
}

func Hash(startsAt string, labels map[string]string) string {
	input := fmt.Sprintf("%s%v", startsAt, labels)
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
