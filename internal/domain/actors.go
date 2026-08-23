package domain

type Actor struct {
	User   User
	Scopes []string
}

func (a Actor) Can(scope string) bool {
	if a.User.Role == RoleAdmin {
		return true
	}
	for _, candidate := range a.Scopes {
		if candidate == scope {
			return true
		}
	}
	return false
}
func (a Actor) CanRegion(region string) bool {
	return a.User.Role == RoleAdmin || NormalizeRegion(a.User.Region) == NormalizeRegion(region)
}
func (a Actor) AddScope(scope string) Actor {
	copy := a
	for _, existing := range copy.Scopes {
		if existing == scope {
			return copy
		}
	}
	copy.Scopes = append(copy.Scopes, scope)
	return copy
}
func (a Actor) RemoveScope(scope string) Actor {
	copy := a
	result := make([]string, 0, len(copy.Scopes))
	for _, existing := range copy.Scopes {
		if existing != scope {
			result = append(result, existing)
		}
	}
	copy.Scopes = result
	return copy
}
func (a Actor) ScopeSet() map[string]struct{} {
	result := make(map[string]struct{}, len(a.Scopes))
	for _, scope := range a.Scopes {
		result[scope] = struct{}{}
	}
	return result
}
