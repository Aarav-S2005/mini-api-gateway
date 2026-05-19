package models

type Middleware struct {
	Name    string                 `bson:"name"`
	Enabled bool                   `bson:"enabled"`
	Config  map[string]interface{} `bson:"config"`
}
