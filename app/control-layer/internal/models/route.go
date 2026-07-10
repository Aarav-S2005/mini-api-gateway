package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthMode string
type PathType string

const (
	AuthNone     AuthMode = "none"
	AuthRequired AuthMode = "required"

	PathExact  PathType = "exact"
	PathPrefix PathType = "prefix"
	PathRegex  PathType = "regex"
)

type Route struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	ProjectID    bson.ObjectID `bson:"project_id"`
	Path         string        `bson:"path"`
	PathType     PathType      `bson:"path_type"`
	Method       string        `bson:"method"`
	UpstreamName string        `bson:"upstream_name"`
	AuthMode     AuthMode      `bson:"auth_mode"`
	Enabled      bool          `bson:"enabled"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}
