package image

import "go.mongodb.org/mongo-driver/v2/bson"

type Filter map[string]any

type FilterBuilder struct {
	filter Filter
	err    error
}

func NewFilter() *FilterBuilder {
	return &FilterBuilder{
		filter: make(map[string]any),
		err:    nil,
	}
}

func (f *FilterBuilder) Build() (Filter, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.filter, nil
}

func (f *FilterBuilder) SetId(id string) *FilterBuilder {
	if f.err != nil {
		return f
	}

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		f.err = err
		return f
	}

	f.filter["_id"] = objId
	return f
}
