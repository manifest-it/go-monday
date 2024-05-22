package monday

import (
	"time"
)

func GetItemsQuery(boardIDs *[]string, startDate, endDate time.Time) Query {
	q := Query{
		Wrap:   true,
		Object: "boards",
		Select: Queries{
			{Object: "id"},
			{Object: "name"},
			{Object: "state"},
			{Object: "columns", Select: Queries{
				{Object: "id"},
				{Object: "title"},
				{Object: "type"},
			}},
			{Object: "owners", Select: Queries{
				{Object: "id"},
			}},
			{
				Object: "items_page", Select: Queries{
					{Object: "items", Select: Queries{
						{Object: "id"},
						{Object: "name"},
						{Object: "url"},
						{Object: "created_at"},
						{Object: "updated_at"},
						{Object: "creator", Select: Queries{
							{Object: "name"},
						}},
						{Object: "board", Select: Queries{
							{Object: "name"},
						}},
						{Object: "updates", Select: Queries{
							{Object: "id"},
							{Object: "text_body"},
							{Object: "creator", Select: Queries{
								{Object: "name"},
							}},
						}},
						{Object: "subscribers", Select: Queries{
							{Object: "id"},
							{Object: "name"},
						}},
						{Object: "created_at"},
						{Object: "column_values", Select: Queries{
							{Object: "id"},
							{Object: "value"},
							{Object: "text"},
						}},
					}},
				},
				Where: map[string]any{
					"query_params": map[string]any{
						"rules": []map[string]any{
							{
								"column_id":         "__creation_log__",
								"compare_value":     []string{startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)},
								"operator":          "between",
								"compare_attribute": "CREATED_AT",
							},
						},
					},
				},
			},
		},
	}

	if boardIDs != nil {
		q.Where = map[string]any{
			"ids": *boardIDs,
		}
	}

	return q
}

func GetAllUsersQuery() Query {
	return Query{
		Wrap:   true,
		Object: "users",
		Select: Queries{
			{Object: "id"},
			{Object: "name"},
			{Object: "email"},
			{Object: "url"},
			{Object: "is_admin"},
			{Object: "is_guest"},
			{Object: "is_view_only"},
			{Object: "created_at"},
		},
	}
}

func GetAllBoardsAndGroupsQuery() Query {
	return Query{
		Wrap:   true,
		Object: "boards",
		Select: Queries{
			{Object: "id"},
			{Object: "name"},
			{Object: "groups", Select: Queries{
				{Object: "id"},
				{Object: "title"},
			}},
		},
	}
}
