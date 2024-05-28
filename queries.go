package monday

import (
	"time"
)

func GetItemsQuery(boardId string, startDate, endDate time.Time, limit int) Query {
	return Query{
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
					{Object: "cursor"},
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
					"limit": limit,
				},
			},
		},
		Where: map[string]any{
			"ids": boardId,
		},
	}
}

func GetNextItemsQuery(cursor string, limit int) Query {
	return Query{
		Wrap:   true,
		Object: "next_items_page",
		Select: Queries{
			{Object: "cursor"},
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
			"cursor": cursor,
			"limit": limit,
		},
	}
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

func GetAllBoardsQuery() Query {
	return Query{
		Wrap:   true,
		Object: "boards",
		Select: Queries{
			{Object: "id"},
			{Object: "name"},
		},
	}
}

func GetUserByIdQuery(userId string) Query {
	return Query{
		Wrap:   true,
		Object: "users",
		Select: Queries{
			{Object: "id"},
			{Object: "name"},
			{Object: "email"},
			{Object: "url"},
			{Object: "is_admin"},
			{Object: "phone"},
			{Object: "title"},
			{Object: "location"},
		},
		Where: map[string]any{
			"ids": userId,
		},
	}
}
