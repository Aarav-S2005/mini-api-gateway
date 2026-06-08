package models

type Middleware struct {
	Name   string                 `bson:"name"`
	Config map[string]interface{} `bson:"config"`
}
