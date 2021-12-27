package dataaccess

import (
	mdl "github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/model"
)

type Repository interface {
	MergePrecedingKMers(kPlus1Mer []byte) error
	AddGenomeMeta(genome *mdl.Genome) error
	AddSequenceKMer(sequence *mdl.SequenceMeta, KMer string)
}
