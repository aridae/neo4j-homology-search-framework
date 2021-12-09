package controller

import (
	repo "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/data-access"
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/utils/workerspool"
)

// команда InitEmptyGraph делится на маленькие таски
type InitGraphTask struct {
	ID        int
	kplus1mer []byte          // что записать
	repo      repo.Repository // куда записать
}

var (
	currID = 0
)

func NewInitGraphTask(kplus1mer []byte, repo repo.Repository) workerspool.Task {
	currID++
	return &InitGraphTask{
		ID:        currID,
		kplus1mer: kplus1mer,
		repo:      repo,
	}
}

func (task *InitGraphTask) GetID() int {
	return task.ID
}

func (task *InitGraphTask) Process() error {
	return task.repo.MergePrecedingKMers(task.kplus1mer)
}

func (task *InitGraphTask) Cleanup() error {
	return nil
}
