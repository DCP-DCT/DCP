package DCP

import "github.com/google/uuid"

type ControlEntity struct {
	Id       uuid.UUID
	BranchId uuid.UUID
}

type CalculationProcessControlEntity struct {
	NodesContributedToUpdates map[ControlEntity]int
}

func (dpce *CalculationProcessControlEntity) RegisterContribution(id uuid.UUID, branchId uuid.UUID, added int) {
	ce := ControlEntity{
		Id:       id,
		BranchId: branchId,
	}

	dpce.NodesContributedToUpdates[ce] = added
}

type Diagnosis struct {
	NumberOfBroadcasts             int
	NumberOfUpdates                int
	NumberOfRejectedDueToThreshold int
	NumberOfDuplicates             int
	NumberOfPkMatches              int
	NumberOfInternalUpdates        int
	NumberOfPacketsDropped         int
	Control                        *CalculationProcessControlEntity
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
