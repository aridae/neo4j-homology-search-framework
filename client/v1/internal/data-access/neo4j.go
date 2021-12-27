package dataaccess

import (
	"github.com/aridae/neo4j-homology-search-framework/client/v1/internal/dbdriver"
	mdl "github.com/aridae/neo4j-homology-search-framework/client/v1/internal/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jRepository struct {
	neo4jClient *dbdriver.Neo4jClient
}

func NewNeo4jRepository(neo4jClient *dbdriver.Neo4jClient) Repository {
	return &Neo4jRepository{
		neo4jClient: neo4jClient,
	}
}

func (rep *Neo4jRepository) MergePrecedingKMers(kPlus1Mer []byte) error {
	session := rep.neo4jClient.CreateSession()
	defer session.Close()

	_, err := session.WriteTransaction(
		func(tx neo4j.Transaction) (interface{}, error) {
			_, err := tx.Run(
				"MERGE (prefix:KMer { value: $prefixValue }) MERGE (suffix:KMer { value: $suffixValue }) MERGE (prefix)-[r:Precedes]->(suffix)",
				map[string]interface{}{
					"prefixValue": string(kPlus1Mer[:len(kPlus1Mer)-1]),
					"suffixValue": string(kPlus1Mer[1:]),
				},
			)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
	return err
}

func (rep *Neo4jRepository) AddGenomeMeta(genome *mdl.Genome) error {
	session := rep.neo4jClient.CreateSession()
	defer session.Close()

	_, err := session.WriteTransaction(
		func(tx neo4j.Transaction) (interface{}, error) {
			for _, seq := range genome.Sequences {
				_, err := tx.Run(
					"MERGE (genome:Genome { name: $genomeName }) MERGE (sequence:Sequence { name: $sequenceName }) MERGE (sequence)-[r:Belongs]->(genome)",
					map[string]interface{}{
						"genomeName":   genome.Name,
						"sequenceName": seq.Name,
					},
				)
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		},
	)
	return err
}

func (rep *Neo4jRepository) AddSequenceKMer(sequence *mdl.Sequence, KMer string, cnt int64) error {
	session := rep.neo4jClient.CreateSession()
	defer session.Close()

	_, err := session.WriteTransaction(
		func(tx neo4j.Transaction) (interface{}, error) {
			_, err := tx.Run(
				"MATCH (sequence:Sequence { name: $sequenceName }) MATCH (n:KMer { value: $kmer }) MERGE (n)-[r:Belongs{count: $count}]->(sequence)",
				map[string]interface{}{
					"sequenceName": sequence.Name,
					"kmer":         KMer,
					"count":        cnt,
				},
			)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
	return err
}
