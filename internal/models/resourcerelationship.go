package models

type ResourceRelationship struct {
	Id               string `json:"id"`
	ParentResourceId string `json:"parent_resource_id"`
	ChildResourceId  string `json:"child_resource_id"`
}
