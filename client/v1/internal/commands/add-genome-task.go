package commands

import (
	repo "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/data-access"
	mdl "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/model"
	"github.com/aridae/neo4j-homology-search-framework/client/v1/internal/utils/workerspool"
)

// команда AddGenome делится на маленькие таски
type AddGenomeTask struct {
	ID    int
	seq   *mdl.Sequence
	kmer  string
	count int64
	repo  repo.Repository // куда записать
}

var (
	currentID = 0
)

func NewAddGenomeTask(repo repo.Repository, seq *mdl.Sequence, kmer string, count int64) workerspool.Task {
	currentID++
	return &AddGenomeTask{
		ID:    currentID,
		seq:   seq,
		kmer:  kmer,
		count: count,
		repo:  repo,
	}
}

func (task *AddGenomeTask) GetID() int {
	return task.ID
}

func (task *AddGenomeTask) Process() error {
	return task.repo.AddSequenceKMer(task.seq, task.kmer, task.count)
}

func (task *AddGenomeTask) Cleanup() error {
	return nil
}
