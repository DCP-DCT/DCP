package DCP

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type ControlEntity struct {
	Id       uuid.UUID `json:"id"`
	BranchId uuid.UUID `json:"branch_id"`
}

func (ce ControlEntity) MarshalText() ([]byte, error) {
	type x ControlEntity
	return json.Marshal(x(ce))
}

func (ce *ControlEntity) UnmarshalText(text []byte) error {
	type x ControlEntity
	return json.Unmarshal(text, (*x)(ce))
}

type CalculationProcessControlEntity struct {
	NodesContributedToUpdates map[ControlEntity]int `json:"nodes_contributed_to_updates"`
}

func (dpce *CalculationProcessControlEntity) RegisterContribution(id uuid.UUID, branchId uuid.UUID, added int) {
	ce := ControlEntity{
		Id:       id,
		BranchId: branchId,
	}

	dpce.NodesContributedToUpdates[ce] = added
}

type Diagnosis struct {
	NumberOfBroadcasts             int                              `json:"number_of_broadcasts"`
	NumberOfUpdates                int                              `json:"number_of_updates"`
	NumberOfRejectedDueToThreshold int                              `json:"number_of_rejected_due_to_threshold"`
	NumberOfDuplicates             int                              `json:"number_of_duplicates"`
	NumberOfPkMatches              int                              `json:"number_of_pk_matches"`
	NumberOfInternalUpdates        int                              `json:"number_of_internal_updates"`
	NumberOfPacketsDropped         int                              `json:"number_of_packets_dropped"`
	Control                        *CalculationProcessControlEntity `json:"control"`
	Timers                         *Timer                           `json:"timers"`
}

func NewDiagnosis() *Diagnosis {
	return &Diagnosis{
		NumberOfBroadcasts:             0,
		NumberOfUpdates:                0,
		NumberOfRejectedDueToThreshold: 0,
		NumberOfDuplicates:             0,
		NumberOfPkMatches:              0,
		NumberOfInternalUpdates:        0,
		NumberOfPacketsDropped:         0,
		Control: &CalculationProcessControlEntity{
			NodesContributedToUpdates: make(map[ControlEntity]int),
		},
		Timers: &Timer{
			Timers: make(map[string]time.Duration),
		},
	}
}

func (d *Diagnosis) Init() {
	d.NumberOfBroadcasts = 0
	d.NumberOfUpdates = 0
	d.NumberOfRejectedDueToThreshold = 0
	d.NumberOfDuplicates = 0
}

func (d *Diagnosis) IncrementNumberOfBroadcasts() {
	d.NumberOfBroadcasts++
}

func (d *Diagnosis) IncrementNumberOfUpdates() {
	d.NumberOfUpdates++
}

func (d *Diagnosis) IncrementNumberOgRejectedDueToThreshold() {
	d.NumberOfRejectedDueToThreshold++
}

func (d *Diagnosis) IncrementNumberOfDuplicates() {
	d.NumberOfDuplicates++
}

func (d *Diagnosis) IncrementNumberOfPkMatches() {
	d.NumberOfPkMatches++
}

func (d *Diagnosis) IncrementNumberOfInternalUpdates() {
	d.NumberOfInternalUpdates++
}

func (d *Diagnosis) IncrementNumberOfPacketsDropped() {
	d.NumberOfPacketsDropped++
}
