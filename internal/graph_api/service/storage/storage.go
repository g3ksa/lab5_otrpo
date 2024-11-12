package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/g3ksa/lab5_otrpo/internal/graph_api/service/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log/slog"
	"strings"
)

type DBStorage struct {
	db neo4j.DriverWithContext
}

func NewDBStorage(db neo4j.DriverWithContext) *DBStorage {
	return &DBStorage{db: db}
}

func (s *DBStorage) GetAllNodes(ctx context.Context) (*[]model.Node, error) {
	return nil, errors.New("not implemented")
	query := "MATCH (n) RETURN id(n) AS id, labels(n) AS label"
	result, err := neo4j.ExecuteQuery(ctx, s.db, query, nil, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	nodes := make([]model.Node, 0, len(result.Records))

	for _, record := range result.Records {
		nodes = append(nodes, model.Node{
			ID:    int(record.Values[0].(int64)),
			Label: record.Values[1].([]interface{})[0].(string),
		})
	}

	return &nodes, nil
}

func (s *DBStorage) GetNodeWithRelations(ctx context.Context, id int) (*[]model.Relation, error) {
	query := `
		MATCH (n)-[r]->(m) WHERE id(n) = $id
		RETURN n {.*}, id(n), type(r) AS relationship_type, m {.*} AS end_node, id(m)`
	result, err := neo4j.ExecuteQuery(ctx, s.db, query, map[string]interface{}{"id": id}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	relations := make([]model.Relation, 0, len(result.Records))
	for _, record := range result.Records {
		relations = append(relations, model.Relation{
			Node: model.Node{
				ID:         int(record.Values[1].(int64)),
				City:       getString(record.Values[0].(map[string]interface{}), "home_town", ""),
				Name:       getString(record.Values[0].(map[string]interface{}), "name", ""),
				Sex:        record.Values[0].(map[string]interface{})["sex"].(int64),
				ScreenName: getString(record.Values[0].(map[string]interface{}), "screen_name", ""),
			},
			RelType: record.Values[2].(string),
			EndNode: model.Node{
				ID:         int(record.Values[4].(int64)),
				Name:       getString(record.Values[3].(map[string]interface{}), "name", ""),
				ScreenName: getString(record.Values[3].(map[string]interface{}), "screen_name", ""),
			},
		})
	}
	return &relations, nil
}

func (s *DBStorage) Insert(ctx context.Context, data model.InsertRequest) error {
	session := s.db.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MERGE (n:User {id: $id, name: $name, screen_name: $screen_name, sex: $sex, home_town: $city})"
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"id":          data.Node.ID,
			"name":        data.Node.Name,
			"screen_name": data.Node.ScreenName,
			"sex":         data.Node.Sex,
			"city":        data.Node.City,
		})
		if err != nil {
			slog.Error("error", err)
			return nil, err
		}
		for _, rel := range data.Relations {
			relType := strings.ToUpper(rel.RelType)
			if relType != "FOLLOWS" && relType != "SUBSCRIBES" {
				continue
			}
			relNode, err := s.GetOneById(ctx, rel.EndNodeID)
			if err != nil {
				return nil, err
			}
			fmt.Println(relNode)
			query := `
				MATCH (n:User {id: $id})
				MATCH (m:` + relNode.Label + `{id: $end_node_id})
				MERGE (n)-[r:` + relType + `]->(m)`
			_, err = tx.Run(ctx, query, map[string]interface{}{
				"id":          data.Node.ID,
				"end_node_id": rel.EndNodeID,
			})
			if err != nil {
				slog.Error("error", err)
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) Delete(ctx context.Context, id int) error {
	query := `MATCH (n) WHERE id(n)=$id DETACH DELETE n`
	_, err := neo4j.ExecuteQuery(ctx, s.db, query, map[string]interface{}{"id": id}, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) GetOneById(ctx context.Context, id int) (model.Node, error) {
	query := "MATCH (n { id: $id }) RETURN id(n) AS id, labels(n) AS label"
	result, err := neo4j.ExecuteQuery(ctx, s.db, query, map[string]interface{}{"id": id}, neo4j.EagerResultTransformer)
	if err != nil {
		return model.Node{}, err
	}

	return model.Node{
		ID:    int(result.Records[0].Values[0].(int64)),
		Label: result.Records[0].Values[1].([]interface{})[0].(string),
	}, nil
}

func getString(data map[string]interface{}, key string, defaultValue string) string {
	if data == nil {
		return defaultValue
	}
	if value, ok := data[key]; ok && value != nil {
		return value.(string)
	}
	return defaultValue
}
