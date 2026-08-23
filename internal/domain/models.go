package domain

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleInspector Role = "inspector"
	RoleRegional  Role = "regional"
)

type QualityClass string

const (
	QualityI   QualityClass = "I"
	QualityII  QualityClass = "II"
	QualityIII QualityClass = "III"
	QualityIV  QualityClass = "IV"
	QualityV   QualityClass = "V"
)

type AlertStatus string

const (
	AlertOpen          AlertStatus = "open"
	AlertInvestigating AlertStatus = "investigating"
	AlertResolved      AlertStatus = "resolved"
)

type InspectionStatus string

const (
	InspectionPending   InspectionStatus = "pending"
	InspectionRunning   InspectionStatus = "running"
	InspectionCompleted InspectionStatus = "completed"
	InspectionFailed    InspectionStatus = "failed"
)

type ManifestStatus string

const (
	ManifestCreated   ManifestStatus = "created"
	ManifestInTransit ManifestStatus = "in_transit"
	ManifestAccepted  ManifestStatus = "accepted"
	ManifestDisposed  ManifestStatus = "disposed"
)

type ProjectStatus string

const (
	ProjectPlanned  ProjectStatus = "planned"
	ProjectBuilding ProjectStatus = "building"
	ProjectAccepted ProjectStatus = "accepted"
)

type User struct {
	ID        int64
	Username  string
	Role      Role
	Region    string
	Disabled  bool
	CreatedAt time.Time
}
type Station struct {
	ID                               int64
	Code, Name, Region, River, Level string
	Active                           bool
	Version                          int64
	CreatedAt                        time.Time
}
type Sample struct {
	ID, StationID, CreatedBy int64
	SampledAt                time.Time
	Quality                  QualityClass
	Metrics                  map[string]float64
	CreatedAt                time.Time
}
type Alert struct {
	ID, StationID, SampleID int64
	Status                  AlertStatus
	Severity                string
	OpenedAt                time.Time
	ClosedAt                *time.Time
	Version                 int64
}
type Inspection struct {
	ID, AlertID, OwnerID   int64
	Status                 InspectionStatus
	DueAt                  time.Time
	StartedAt, CompletedAt *time.Time
	Notes                  string
	Version                int64
}
type RemediationAction struct {
	ID, InspectionID         int64
	Action, Status, Evidence string
	CreatedAt                time.Time
	CompletedAt              *time.Time
}
type Agreement struct {
	ID                                               int64
	UpstreamRegion, DownstreamRegion, ThresholdClass string
	Active                                           bool
	CreatedAt                                        time.Time
}
type Warning struct {
	ID, AgreementID, StationID int64
	Direction, Payload, Status string
	SentAt                     time.Time
	AcknowledgedAt             *time.Time
}
type Compensation struct {
	ID, AgreementID                   int64
	Period, Direction, Reason, Status string
	AmountCents                       int64
	CreatedAt, SettledAt              *time.Time
}
type Manifest struct {
	ID                                                               int64
	Number, ProducerRegion, CarrierRegion, FacilityRegion, WasteType string
	WeightKG                                                         int64
	Status                                                           ManifestStatus
	CreatedAt, AcceptedAt, DisposedAt                                *time.Time
	Version                                                          int64
}
type Project struct {
	ID             int64
	Name, Region   string
	TargetHectares float64
	Status         ProjectStatus
	BudgetCents    int64
	CreatedAt      time.Time
	Version        int64
}
type Milestone struct {
	ID, ProjectID int64
	Name          string
	TargetDate    time.Time
	Status        string
	Evidence      string
	CompletedAt   *time.Time
}
