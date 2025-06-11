package service

import (
	"collectionsservice/internal/models"
	proto "collectionsservice/internal/proto"
	"collectionsservice/internal/repository"
	"collectionsservice/internal/utils"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
)

type CollectionService struct {
	Repo repository.CollectionRepoInterface
	proto.CollectionServiceServer
}

type CollectionServiceInterface interface {
	CreateCollection(ctx context.Context, req *proto.CreateCollectionRequest) (*proto.CreateCollectionResponse, error)
	AddRequestToCollection(ctx context.Context, req *proto.AddRequestToCollectionRequest) (*proto.CollectionResponse, error)
	ListCollectionsAndRequests(ctx context.Context) ([]*models.Collection, error)
	UpdateCollection(ctx context.Context, req *proto.UpdateCollectionRequest) (*proto.CollectionResponse, error)
	UpdateRequestInCollection(ctx context.Context, req *proto.UpdateRequestInCollectionRequest) (*proto.CollectionResponse, error)
	DeleteCollection(ctx context.Context, collectionID string) (*proto.DeleteResponse, error)
	DeleteRequestFromCollection(ctx context.Context, collectionID, requestID string) (*proto.DeleteResponse, error)
}

func NewCollectionService(repo repository.CollectionRepoInterface) *CollectionService {
	return &CollectionService{
		Repo: repo,
	}
}

func (s *CollectionService) CreateCollection(ctx context.Context, req *proto.CreateCollectionRequest) (*proto.CreateCollectionResponse, error) {
	if req.GetName() == "" {
		return nil, errors.New("collection name cannot be empty")
	}

	collection := models.Collection{
		ID:          uuid.New().String(),
		Name:        req.GetName(),
		Description: &req.Description,
	}

	id, err := s.Repo.CreateCollection(ctx, collection)
	if err != nil {
		// Log error creating collection
		log.Error().Err(err).Str("collection_name", collection.Name).Msg("Failed to create collection")
		return nil, err
	}

	return &proto.CreateCollectionResponse{
		Id:   id,
		Name: collection.Name,
	}, nil
}

func (s *CollectionService) AddRequestToCollection(ctx context.Context, req *proto.AddRequestToCollectionRequest) (*proto.CollectionResponse, error) {
	if s.Repo == nil {
		return nil, fmt.Errorf("repository is not initialized")
	}

	reqModel, err := utils.ConvertProtoRequests([]*proto.CollectionRequestInput{req.Request}, req.CollectionName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert request input")
		return nil, fmt.Errorf("failed to convert request input: %w", err)
	}

	err = s.Repo.AddRequestToCollection(ctx, req.CollectionName, reqModel)
	if err != nil {
		log.Error().Err(err).Str("collection_name", req.CollectionName).Msg("Failed to add request to collection")
		return nil, fmt.Errorf("failed to add request to collection: %w", err)
	}

	updatedCollection, err := s.Repo.GetCollectionByName(ctx, req.CollectionName)
	if err != nil {
		log.Error().Err(err).Str("collection_name", req.CollectionName).Msg("Failed to fetch updated collection")
		return nil, fmt.Errorf("failed to fetch updated collection: %w", err)
	}

	protoCollection := utils.ConvertModelCollectionToProto(updatedCollection)

	return protoCollection, nil
}

func (s *CollectionService) ListCollectionsAndRequests(ctx context.Context, req *proto.ListCollectionsRequest) (*proto.ListCollectionsResponse, error) {
	if s.Repo == nil {
		return nil, fmt.Errorf("repository is not initialized")
	}

	collections, err := s.Repo.ListCollectionsAndRequests(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list collections and requests")
		return nil, fmt.Errorf("failed to list collections and requests: %w", err)
	}

	// no logs here; normal processing

	var protoCollections []*proto.CollectionResponse

	for _, col := range collections {
		pCol := &proto.CollectionResponse{
			Id:          col.ID,
			Name:        col.Name,
			Description: "",
		}
		if col.Description != nil {
			pCol.Description = *col.Description
		}

		for _, r := range col.Requests {
			var pReq *proto.CollectionRequest

			switch r.Kind {
			case models.RequestKindHTTP:
				method := proto.HTTPMethod_GET
				if r.HTTPMethod != nil {
					if val, ok := proto.HTTPMethod_value[*r.HTTPMethod]; ok {
						method = proto.HTTPMethod(val)
					}
				}
				url := ""
				if r.HTTPURL != nil {
					url = *r.HTTPURL
				}
				pReq = &proto.CollectionRequest{
					Request: &proto.CollectionRequest_HttpRequest{
						HttpRequest: &proto.HTTPRequest{
							Name:   r.Name,
							Method: method,
							Url:    url,
						},
					},
				}

			case models.RequestKindGraphQL:
				endpoint := ""
				if r.GraphQLEndpoint != nil {
					endpoint = *r.GraphQLEndpoint
				}
				query := ""
				if r.GraphQLQuery != nil {
					query = *r.GraphQLQuery
				}
				pReq = &proto.CollectionRequest{
					Request: &proto.CollectionRequest_GraphqlRequest{
						GraphqlRequest: &proto.GraphQLRequest{
							Name:     r.Name,
							Endpoint: endpoint,
							Query:    query,
						},
					},
				}

			default:
				continue
			}

			pCol.Requests = append(pCol.Requests, pReq)
		}

		pCol.RequestCount = int32(len(pCol.Requests))
		protoCollections = append(protoCollections, pCol)
	}

	return &proto.ListCollectionsResponse{
		Collections: protoCollections,
	}, nil
}

func (s *CollectionService) UpdateCollection(ctx context.Context, req *proto.UpdateCollectionRequest) (*proto.CollectionResponse, error) {
	existing, err := s.Repo.GetByID(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).Str("collection_id", req.Id).Msg("Failed to get collection by ID")
		return nil, err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = &req.Description
	}

	updated, err := s.Repo.Update(ctx, existing)
	if err != nil {
		log.Error().Err(err).Str("collection_id", req.Id).Msg("Failed to update collection")
		return nil, err
	}
	return &proto.CollectionResponse{
		Id:          updated.ID,
		Name:        updated.Name,
		Description: *updated.Description,
	}, nil
}

func (s *CollectionService) UpdateRequestInCollection(ctx context.Context, req *proto.UpdateRequestInCollectionRequest) (*proto.UpdateRequestInCollectionResponse, error) {
	input := &models.Request{}

	if req.Name != "" {
		input.Name = req.Name
	}
	if req.Kind != proto.RequestKind(0) {
		input.Kind = models.RequestKind(req.Kind.String())
	}
	if req.HttpMethod != "" {
		input.HTTPMethod = &req.HttpMethod
	}
	if req.HttpUrl != "" {
		input.HTTPURL = &req.HttpUrl
	}
	if req.HttpHeaders != "" {
		input.HTTPHeaders = datatypes.JSON([]byte(req.HttpHeaders))
	}
	if req.HttpQueryParams != "" {
		input.HTTPQueryParams = datatypes.JSON([]byte(req.HttpQueryParams))
	}
	if req.HttpBody != "" {
		input.HTTPBody = &req.HttpBody
	}
	if req.GraphqlEndpoint != "" {
		input.GraphQLEndpoint = &req.GraphqlEndpoint
	}
	if req.GraphqlQuery != "" {
		input.GraphQLQuery = &req.GraphqlQuery
	}
	if req.GraphqlVariables != "" {
		input.GraphQLVariables = datatypes.JSON([]byte(req.GraphqlVariables))
	}
	if req.GraphqlHeaders != "" {
		input.GraphQLHeaders = datatypes.JSON([]byte(req.GraphqlHeaders))
	}

	collectionAndRequests, err := s.Repo.UpdateRequestInCollection(ctx, req.CollectionId, req.RequestId, input)
	if err != nil {
		log.Error().Err(err).Str("request_id", req.RequestId).Msg("Failed to update request in collection")
		return nil, err
	}

	return &proto.UpdateRequestInCollectionResponse{
		Message:   collectionAndRequests.Message,
		RequestId: req.RequestId,
	}, nil
}

func (s *CollectionService) DeleteCollection(ctx context.Context, req *proto.DeleteCollectionRequest) (*proto.DeleteResponse, error) {
	err := s.Repo.DeleteCollection(req.Id)
	if err != nil {
		log.Error().Err(err).Str("collection_id", req.Id).Msg("Failed to delete collection")
		return &proto.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete collection: %v", err),
		}, nil
	}

	return &proto.DeleteResponse{
		Success: true,
		Message: "Collection deleted successfully",
	}, nil
}

func (s *CollectionService) DeleteRequestFromCollection(ctx context.Context, req *proto.DeleteRequestFromCollectionRequest) (*proto.DeleteResponse, error) {
	err := s.Repo.RemoveRequestFromCollection(req.CollectionId, req.RequestId)
	if err != nil {
		log.Error().Err(err).Str("request_id", req.RequestId).Str("collection_id", req.CollectionId).Msg("Failed to remove request from collection")
		return &proto.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to remove request: %v", err),
		}, nil
	}

	return &proto.DeleteResponse{
		Success: true,
		Message: "Request removed from collection successfully",
	}, nil
}
