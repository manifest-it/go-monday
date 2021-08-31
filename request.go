package monday

import (
	"fmt"
	"strings"
)

type Query struct {
	Object string
	Where  Where
	Select Queries
	Wrap   bool
}

type Queries []Query

func (q Query) String() string {
	qParts := []string{q.Object}
	if len(q.Where) > 0 {
		qParts = append(qParts, q.Where.String())
	}
	if len(q.Select) == 0 {
		return WrapQuery(strings.Join(qParts, " "), q.Wrap)
	}
	subfields := []string{}
	for _, subfield := range q.Select {
		subfields = append(subfields, subfield.String())
	}
	if len(subfields) > 0 {
		qParts = append(qParts,
			"{"+strings.Join(subfields, " ")+"}")
	}
	return WrapQuery(strings.Join(qParts, " "), q.Wrap)
}

type Where map[string]string

func (w Where) String() string {
	if len(w) == 0 {
		return ""
	}
	parts := []string{}
	for k, v := range w {
		parts = append(parts, k+":"+v)
	}
	return "(" + strings.Join(parts, ",") + ")"
}

func WrapQuery(gql string, wrap bool) string {
	if wrap {
		return fmt.Sprintf("query {%s}", gql)
	}
	return gql
}

func BoardQuery(boardId string) Query {
	return Query{
		Wrap:   true,
		Object: "boards",
		Where: map[string]string{
			"ids": boardId},
		Select: Queries{
			{Object: "name"},
			{Object: "state"},
			{Object: "columns", Select: Queries{
				{Object: "id"},
				{Object: "title"},
				{Object: "type"},
			}},
			{Object: "owner", Select: Queries{
				{Object: "id"},
			}},
			{Object: "items", Select: Queries{
				{Object: "id"},
				{Object: "name"},
				{Object: "state"},
				{Object: "column_values", Select: Queries{
					{Object: "title"},
					{Object: "id"},
					{Object: "value"},
					{Object: "text"},
				}},
			}},
		},
	}
}

/*

gql := "query {
	boards (ids: 12345) {
		name
		columns { id title type }
		owner {id}
		items{id name state column_values {title id value text } } state
	}
}"

*/
