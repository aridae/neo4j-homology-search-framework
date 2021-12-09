package dataaccess

import (
	"github.com/aridae/neo4j-homology-search-framework/backend/v1/internal/dbdriver"
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
