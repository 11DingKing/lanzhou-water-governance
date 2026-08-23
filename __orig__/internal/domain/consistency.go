package domain

import "fmt"

type ConsistencyIssue struct {
	Object  string
	Field   string
	Message string
}

func ValidateSample(s Sample, station Station, user User) []ConsistencyIssue {
	issues := make([]ConsistencyIssue, 0)
	if s.StationID != station.ID {
		issues = append(issues, ConsistencyIssue{"sample", "station_id", "sample station mismatch"})
	}
	if user.Role != RoleAdmin && station.Region != user.Region {
		issues = append(issues, ConsistencyIssue{"sample", "region", "operator outside region"})
	}
	if len(s.Metrics) == 0 {
		issues = append(issues, ConsistencyIssue{"sample", "metrics", "empty metrics"})
	}
	if s.SampledAt.IsZero() {
		issues = append(issues, ConsistencyIssue{"sample", "sampled_at", "missing sample time"})
	}
	return issues
}
func ValidateManifest(m Manifest) []ConsistencyIssue {
	issues := make([]ConsistencyIssue, 0)
	if m.Number == "" {
		issues = append(issues, ConsistencyIssue{"manifest", "number", "manifest number required"})
	}
	if m.WeightKG <= 0 {
		issues = append(issues, ConsistencyIssue{"manifest", "weight_kg", "weight must be positive"})
	}
	if m.ProducerRegion == m.FacilityRegion {
		issues = append(issues, ConsistencyIssue{"manifest", "facility_region", "facility must be independently registered"})
	}
	return issues
}
func ValidateProject(p Project) error {
	if p.Name == "" {
		return fmt.Errorf("project name required")
	}
	if p.TargetHectares <= 0 {
		return fmt.Errorf("target hectares must be positive")
	}
	if p.BudgetCents <= 0 {
		return fmt.Errorf("budget must be positive")
	}
	return nil
}
