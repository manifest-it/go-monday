package monday

import (
	"encoding/json"
	"fmt"
	"io"
)

func BoardQuery(boardID string) Query {
	return Query{
		Wrap:   true,
		Object: "boards",
		Where: map[string]string{
			"ids": boardID},
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

func QueryBoard(c Client, boardID string) (*Response, error) {
	gql := BoardQuery(boardID)
	httpResp, err := c.DoGraphQL(gql)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status code [%d]", httpResp.StatusCode)
	}

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var brds Response
	return &brds, json.Unmarshal(data, &brds)
}
