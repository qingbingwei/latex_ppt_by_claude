package vectordb

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusClient struct {
	client         client.Client
	collectionName string
}

func NewMilvusClient(host, port string) (*MilvusClient, error) {
	c, err := client.NewClient(context.Background(), client.Config{
		Address: fmt.Sprintf("%s:%s", host, port),
	})
	if err != nil {
		return nil, err
	}

	return &MilvusClient{
		client:         c,
		collectionName: "document_chunks",
	}, nil
}

func (m *MilvusClient) CreateCollection(ctx context.Context) error {
	// Check if collection exists
	has, err := m.client.HasCollection(ctx, m.collectionName)
	if err != nil {
		return err
	}

	if has {
		log.Printf("Collection %s already exists", m.collectionName)
		return nil
	}

	// Create schema
	schema := &entity.Schema{
		CollectionName: m.collectionName,
		AutoID:         true,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "chunk_id",
				DataType: entity.FieldTypeInt64,
			},
			{
				Name:     "document_id",
				DataType: entity.FieldTypeInt64,
			},
			{
				Name:     "content",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "65535",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "1536", // OpenAI ada-002 dimension
				},
			},
		},
	}

	err = m.client.CreateCollection(ctx, schema, 2)
	if err != nil {
		return err
	}

	// Create index
	idx, err := entity.NewIndexAUTOINDEX(entity.L2)
	if err != nil {
		return err
	}

	return m.client.CreateIndex(ctx, m.collectionName, "embedding", idx, false)
}

func (m *MilvusClient) Insert(ctx context.Context, chunkID, documentID int64, content string, embedding []float32) (string, error) {
	// Load collection
	err := m.client.LoadCollection(ctx, m.collectionName, false)
	if err != nil {
		return "", err
	}

	chunkIDColumn := entity.NewColumnInt64("chunk_id", []int64{chunkID})
	documentIDColumn := entity.NewColumnInt64("document_id", []int64{documentID})
	contentColumn := entity.NewColumnVarChar("content", []string{content})
	embeddingColumn := entity.NewColumnFloatVector("embedding", 1536, [][]float32{embedding})

	_, err = m.client.Insert(ctx, m.collectionName, "", chunkIDColumn, documentIDColumn, contentColumn, embeddingColumn)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", chunkID), nil
}

func (m *MilvusClient) Search(ctx context.Context, embedding []float32, topK int) ([]SearchResult, error) {
	err := m.client.LoadCollection(ctx, m.collectionName, false)
	if err != nil {
		return nil, err
	}

	sp, _ := entity.NewIndexAUTOINDEXSearchParam(1)
	
	results, err := m.client.Search(
		ctx,
		m.collectionName,
		nil,
		"",
		[]string{"chunk_id", "document_id", "content"},
		[]entity.Vector{entity.FloatVector(embedding)},
		"embedding",
		entity.L2,
		topK,
		sp,
	)

	if err != nil {
		return nil, err
	}

	var searchResults []SearchResult
	if len(results) > 0 {
		for i := 0; i < results[0].ResultCount; i++ {
			chunkID, _ := results[0].Fields.GetColumn("chunk_id").Get(i)
			documentID, _ := results[0].Fields.GetColumn("document_id").Get(i)
			content, _ := results[0].Fields.GetColumn("content").Get(i)

			searchResults = append(searchResults, SearchResult{
				ChunkID:    chunkID.(int64),
				DocumentID: documentID.(int64),
				Content:    content.(string),
				Score:      results[0].Scores[i],
			})
		}
	}

	return searchResults, nil
}

func (m *MilvusClient) Close() {
	if m.client != nil {
		m.client.Close()
	}
}

type SearchResult struct {
	ChunkID    int64
	DocumentID int64
	Content    string
	Score      float32
}
