package domain

func CanAlertTransition(from, to AlertStatus) bool {
	switch from {
	case AlertOpen:
		return to == AlertInvestigating
	case AlertInvestigating:
		return to == AlertResolved
	case AlertResolved:
		return false
	}
	return false
}
func CanInspectionTransition(from, to InspectionStatus) bool {
	switch from {
	case InspectionPending:
		return to == InspectionRunning || to == InspectionFailed
	case InspectionRunning:
		return to == InspectionCompleted || to == InspectionFailed
	case InspectionCompleted, InspectionFailed:
		return false
	}
	return false
}
func CanManifestTransition(from, to ManifestStatus) bool {
	switch from {
	case ManifestCreated:
		return to == ManifestInTransit
	case ManifestInTransit:
		return to == ManifestAccepted
	case ManifestAccepted:
		return to == ManifestDisposed
	case ManifestDisposed:
		return false
	}
	return false
}
func CanProjectTransition(from, to ProjectStatus) bool {
	switch from {
	case ProjectPlanned:
		return to == ProjectBuilding
	case ProjectBuilding:
		return to == ProjectAccepted
	case ProjectAccepted:
		return false
	}
	return false
}

func ConcurrentTransitionAllowed(from, to InspectionStatus) bool { return CanInspectionTransition(from, to) }
