package DCP

type Diagnosis struct {
	NumberOfBroadcasts             int
	NumberOfUpdates                int
	NumberOfRejectedDueToThreshold int
	NumberOfDuplicates             int
}

func NewDiagnosis() *Diagnosis {
	return &Diagnosis{
		NumberOfBroadcasts:             0,
		NumberOfUpdates:                0,
		NumberOfRejectedDueToThreshold: 0,
		NumberOfDuplicates:             0,
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
