package model

type ModelStatus int64

const (
	Initing      ModelStatus = iota
	AddingGenome ModelStatus = iota
)

type Model struct {
	K            int64
	GenomesCount int64
	Status       ModelStatus
}

type SequenceMeta struct {
	Name string `json:"name"`
}

type Sequence struct {
	SequenceMeta
	Data []byte `json:"data"`
}

type Genome struct {
	Name      string     `json:"genome"`
	Sequences []Sequence `json:"sequences"`
}

func EmptyModel() *Model {
	return &Model{}
}

func (model *Model) SetK(k int64) {
	model.K = k
}

func (model *Model) UpdateGenomesCount() {
	model.GenomesCount++
}
