package repository

import (
	"collectionsservice/internal/models"
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	DB *gorm.DB
}

type CollectionRepoInterface interface {
	CreateCollection(ctx context.Context, collection models.Collection) (string, error)
	AddRequestToCollection(ctx context.Context, collectionName string, req []models.Request) error
	GetCollectionByName(ctx context.Context, name string) (*models.Collection, error)
	ListCollectionsAndRequests(ctx context.Context) ([]*models.Collection, error)
	GetByID(ctx context.Context, id string) (*models.Collection, error)
	Update(ctx context.Context, collection *models.Collection) (*models.Collection, error)
	UpdateRequestInCollection(ctx context.Context, collectionID, requestID string, input *models.Request) (*models.UpdateRequestInCollectionResponse, error)
	DeleteCollection(collectionID string) error
	RemoveRequestFromCollection(collectionID, requestID string) error
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{
		DB: db,
	}
}

func (r *CollectionRepository) CreateCollection(ctx context.Context, collection models.Collection) (string, error) {
	log.Info().Str("name", collection.Name).Msg("Creating collection")
	if err := r.DB.WithContext(ctx).Create(&collection).Error; err != nil {
		log.Error().Err(err).Str("name", collection.Name).Msg("Failed to create collection")
		return "", err
	}
	log.Info().Str("id", collection.ID).Msg("Collection created")
	return collection.ID, nil
}

func (r *CollectionRepository) AddRequestToCollection(ctx context.Context, collectionName string, reqs []models.Request) error {
	var collection models.Collection
	if err := r.DB.WithContext(ctx).Where("name = ?", collectionName).First(&collection).Error; err != nil {
		log.Error().Err(err).Str("collection_name", collectionName).Msg("Collection not found")
		return err
	}

	for i := range reqs {
		reqs[i].CollectionID = collection.ID
	}

	if err := r.DB.WithContext(ctx).Create(&reqs).Error; err != nil {
		log.Error().Err(err).Str("collection_id", collection.ID).Msg("Failed to add requests")
		return err
	}

	log.Info().Str("collection_id", collection.ID).Int("count", len(reqs)).Msg("Requests added")
	return nil
}

func (r *CollectionRepository) GetCollectionByName(ctx context.Context, name string) (*models.Collection, error) {
	var collection models.Collection
	err := r.DB.WithContext(ctx).Preload("Requests").Where("name = ?", name).First(&collection).Error
	if err != nil {
		log.Error().Err(err).Str("collection_name", name).Msg("Failed to get collection")
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepository) ListCollectionsAndRequests(ctx context.Context) ([]*models.Collection, error) {
	var collections []*models.Collection
	err := r.DB.WithContext(ctx).Preload("Requests").Find(&collections).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to list collections")
		return nil, err
	}
	return collections, nil
}

func (r *CollectionRepository) GetByID(ctx context.Context, id string) (*models.Collection, error) {
	var collection models.Collection
	if err := r.DB.WithContext(ctx).First(&collection, "id = ?", id).Error; err != nil {
		log.Error().Err(err).Str("collection_id", id).Msg("Failed to fetch collection")
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepository) Update(ctx context.Context, collection *models.Collection) (*models.Collection, error) {
	log.Info().Str("collection_id", collection.ID).Msg("Updating collection")
	if err := r.DB.WithContext(ctx).Save(collection).Error; err != nil {
		log.Error().Err(err).Str("collection_id", collection.ID).Msg("Update failed")
		return nil, err
	}
	log.Info().Str("collection_id", collection.ID).Msg("Collection updated")
	return collection, nil
}

func (r *CollectionRepository) UpdateRequestInCollection(ctx context.Context, collectionID, requestID string, input *models.Request) (*models.UpdateRequestInCollectionResponse, error) {
	var request models.Request
	if err := r.DB.Where("id = ? AND collection_id = ?", requestID, collectionID).First(&request).Error; err != nil {
		log.Error().Err(err).Str("request_id", requestID).Msg("Request not found")
		return nil, err
	}

	updates := map[string]interface{}{
		"kind":               input.Kind,
		"name":               input.Name,
		"http_method":        input.HTTPMethod,
		"http_url":           input.HTTPURL,
		"http_headers":       input.HTTPHeaders,
		"http_query_params":  input.HTTPQueryParams,
		"http_body":          input.HTTPBody,
		"graph_ql_endpoint":  input.GraphQLEndpoint,
		"graph_ql_query":     input.GraphQLQuery,
		"graph_ql_variables": input.GraphQLVariables,
		"graph_ql_headers":   input.GraphQLHeaders,
	}

	if err := r.DB.Model(&request).Updates(updates).Error; err != nil {
		log.Error().Err(err).Str("request_id", requestID).Msg("Failed to update request")
		return nil, err
	}

	log.Info().Str("request_id", requestID).Msg("Request updated")
	return &models.UpdateRequestInCollectionResponse{
		Message:   "Request Updated successfully",
		RequestID: requestID,
	}, nil
}

func (r *CollectionRepository) DeleteCollection(collectionID string) error {
	log.Info().Str("collection_id", collectionID).Msg("Deleting collection")
	err := r.DB.Delete(&models.Collection{}, "id = ?", collectionID).Error
	if err != nil {
		log.Error().Err(err).Str("collection_id", collectionID).Msg("Delete failed")
		return err
	}
	log.Info().Str("collection_id", collectionID).Msg("Collection deleted")
	return nil
}

func (r *CollectionRepository) RemoveRequestFromCollection(collectionID, requestID string) error {
	var request models.Request

	err := r.DB.First(&request, "id = ? AND collection_id = ?", requestID, collectionID).Error
	if err != nil {
		log.Error().Err(err).Str("request_id", requestID).Str("collection_id", collectionID).Msg("Request not found in collection")
		return errors.New("request not found in collection")
	}

	if err := r.DB.Delete(&request).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete request")
		return err
	}

	log.Info().Str("request_id", requestID).Msg("Request deleted from collection")
	return nil
}

