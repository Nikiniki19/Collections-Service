package utils

import (
	"collectionsservice/internal/models"
	proto "collectionsservice/internal/proto"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func ConvertProtoRequests(protoReqs []*proto.CollectionRequestInput, collectionID string) ([]models.Request, error) {
	var requests []models.Request

	for _, r := range protoReqs {
		req := models.Request{
			Kind:         models.RequestKind(r.Kind.String()),
			Name:         r.Name,
			ID:           uuid.New().String(),
			CollectionID: collectionID,
		}

		switch r.Kind {
		case proto.RequestKind_HTTP:
			if r.Http != nil {
				method := r.Http.Method.String()
				url := r.Http.Url
				req.HTTPMethod = &method
				req.HTTPURL = &url

				headers, err := json.Marshal(r.Http.Headers)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal HTTP headers: %w", err)
				}
				req.HTTPHeaders = datatypes.JSON(headers)

				queryParams, err := json.Marshal(r.Http.QueryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal HTTP query params: %w", err)
				}
				req.HTTPQueryParams = datatypes.JSON(queryParams)

				if r.Http.Body != nil {
					jsonBody, err := r.Http.Body.MarshalJSON()
					if err != nil {
						return nil, fmt.Errorf("failed to marshal HTTP body: %w", err)
					}
					bodyStr := string(jsonBody)
					req.HTTPBody = &bodyStr
				}
			}

		case proto.RequestKind_GRAPHQL:
			if r.Graphql != nil {
				endpoint := r.Graphql.Endpoint
				query := r.Graphql.Query
				req.GraphQLEndpoint = &endpoint
				req.GraphQLQuery = &query

				if r.Graphql.Variables != nil {
					jsonBytes, err := r.Graphql.Variables.MarshalJSON()
					if err != nil {
						return nil, fmt.Errorf("failed to marshal GraphQL variables: %w", err)
					}
					req.GraphQLVariables = datatypes.JSON(jsonBytes)
				}

				headers, err := json.Marshal(r.Graphql.Headers)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal GraphQL headers: %w", err)
				}
				req.GraphQLHeaders = datatypes.JSON(headers)
			}
		}

		requests = append(requests, req)
	}

	return requests, nil
}

func ConvertModelCollectionToProto(col *models.Collection) *proto.CollectionResponse {
	pCol := &proto.CollectionResponse{
		Id:           col.ID,
		Name:         col.Name,
		Description:  "",
		RequestCount: int32(len(col.Requests)),
		Requests:     []*proto.CollectionRequest{},
	}

	if col.Description != nil {
		pCol.Description = *col.Description
	}

	for _, r := range col.Requests {
		pReq := &proto.CollectionRequest{}

		switch r.Kind {
		case models.RequestKindHTTP:
			pReq.Request = &proto.CollectionRequest_HttpRequest{
				HttpRequest: &proto.HTTPRequest{
					Name:   r.Name,
					Method: proto.HTTPMethod(proto.HTTPMethod_value[*r.HTTPMethod]),
					Url:    *r.HTTPURL,
				},
			}
		case models.RequestKindGraphQL:
			pReq.Request = &proto.CollectionRequest_GraphqlRequest{
				GraphqlRequest: &proto.GraphQLRequest{
					Name:     r.Name,
					Endpoint: *r.GraphQLEndpoint,
					Query:    *r.GraphQLQuery,
				},
			}
		}

		pCol.Requests = append(pCol.Requests, pReq)
	}

	return pCol
}




