package controller

import (
	"log"

	repo "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/data-access"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/chunkreader"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/workerspool"
)

// команда AddGenome делится на маленькие таски
type AddGenomeTask struct {
	ID          int
	genomeChunk *[]byte                  // что записать
	repo        repo.Repository          // куда записать
	reader      *chunkreader.ChunkReader // потом куда вернуть чанк
}

var (
	currentID = 0
)

func NewAddGenomeTask(genomeChunk *[]byte, repo repo.Repository, reader *chunkreader.ChunkReader) workerspool.Task {
	currentID++
	return &AddGenomeTask{
		ID:          currentID,
		genomeChunk: genomeChunk,
		repo:        repo,
		reader:      reader,
	}
}

func (task *AddGenomeTask) GetID() int {
	return task.ID
}

func (task *AddGenomeTask) Process() error {
	log.Printf("[task %d] merging sequence [but not really]: %s\n", task.GetID(), string(*task.genomeChunk))
	return nil
}

func (task *AddGenomeTask) Cleanup() error {
	task.reader.FreeChunk(task.genomeChunk)
	return nil
}
