package inputs

import (
	"net/url"
)

type QueryBuilder struct {
	values url.Values
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		values: url.Values{},
	}
}

func (b *QueryBuilder) Add(key, value string) *QueryBuilder {
	b.values.Add(key, value)
	return b
}

func (b *QueryBuilder) AddMapped(flagName, key, value string) *QueryBuilder {
	b.values.Add(key, value)
	return b
}

func (b *QueryBuilder) AddAll(key string, values []string) *QueryBuilder {
	for _, v := range values {
		b.values.Add(key, v)
	}
	return b
}

func (b *QueryBuilder) Encode() string {
	return b.values.Encode()
}
