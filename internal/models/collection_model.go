package models

import (
	"gorm.io/datatypes"
)

type RequestKind string

const (
	RequestKindHTTP    RequestKind = "HTTP"
	RequestKindGraphQL RequestKind = "GRAPHQL"
)

type Collection struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description *string   `gorm:"type:text"`
	Requests    []Request `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE"`
}

type Request struct {
	ID           string      `gorm:"type:uuid;primaryKey"`
	CollectionID string      `gorm:"type:uuid;not null;index"`
	Kind         RequestKind `gorm:"type:text;not null"`
	Name         string      `gorm:"not null"`

	HTTPMethod      *string        `gorm:"type:text"`
	HTTPURL         *string        `gorm:"type:text"`
	HTTPHeaders     datatypes.JSON `gorm:"type:jsonb"`
	HTTPQueryParams datatypes.JSON `gorm:"type:jsonb"`
	HTTPBody        *string        `gorm:"type:text"`

	GraphQLEndpoint  *string        `gorm:"type:text"`
	GraphQLQuery     *string        `gorm:"type:text"`
	GraphQLVariables datatypes.JSON `gorm:"type:jsonb"`
	GraphQLHeaders   datatypes.JSON `gorm:"type:jsonb"`
}

type UpdateRequestInCollectionParams struct {
	CollectionID string `gorm:"-"`
	RequestID    string `gorm:"-"`

	Name        *string `gorm:"column:name"`
	Method      *string `gorm:"column:method"`
	URL         *string `gorm:"column:url"`
	Body        *string `gorm:"column:body"`
	Description *string `gorm:"column:description"`
}

type UpdateRequestInCollectionResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"requestID"`
}
