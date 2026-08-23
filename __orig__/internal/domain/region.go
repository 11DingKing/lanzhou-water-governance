package domain

import "strings"

type RegionRelation string

const (
	Upstream    RegionRelation = "upstream"
	Downstream  RegionRelation = "downstream"
	CrossBorder RegionRelation = "cross-border"
)

func NormalizeRegion(value string) string { return strings.ToUpper(strings.TrimSpace(value)) }
func Relation(from, to string) RegionRelation {
	left, right := NormalizeRegion(from), NormalizeRegion(to)
	if left == right {
		return CrossBorder
	}
	if left < right {
		return Upstream
	}
	return Downstream
}
func CanExchangeData(from, to string, agreement Agreement) bool {
	return agreement.Active && NormalizeRegion(agreement.UpstreamRegion) == NormalizeRegion(from) && NormalizeRegion(agreement.DownstreamRegion) == NormalizeRegion(to)
}
func CompensationDirection(relation RegionRelation) string {
	switch relation {
	case Upstream:
		return "upstream-to-downstream"
	case Downstream:
		return "downstream-to-upstream"
	}
	return "bilateral"
}
func RegionSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := NormalizeRegion(value)
		if normalized != "" {
			result[normalized] = struct{}{}
		}
	}
	return result
}
