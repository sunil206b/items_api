package model

import "github.com/olivere/elastic"

type EsQuery struct {
	Equals []FieldValue `json:"equals"`
}

type FieldValue struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}


func (q EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()
	equalsQuery := make([]elastic.Query, 0)
	for _, eq := range q.Equals {
		equalsQuery = append(equalsQuery, elastic.NewMatchQuery(eq.Field, eq.Value))
	}
	query.Must(equalsQuery...)
	return query
}