package monday

import (
	"time"
)

const (
	ColumnValueTitleDate   = "Date"
	ColumnValueTitleStatus = "Status"
)

type ItemsResponse struct {
	AccountID int `json:"account_id"`
	Data      struct {
		Boards []BoardItems `json:"boards"`
	} `json:"data"`
}

type NextItemsResponse struct {
	AccountID int `json:"account_id"`
	Data      struct {
		ItemsPage ItemsPage `json:"next_items_page"`
	} `json:"data"`
}

type UsersResponse struct {
	AccountID int `json:"account_id"`
	Data      struct {
		Users []User `json:"users"`
	} `json:"data"`
}

type BoardsResponse struct {
	AccountID int `json:"account_id"`
	Data      struct {
		Boards []BoardGroups `json:"boards"`
	} `json:"data"`
}

type BoardItems struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
	Owners  []struct {
		ID string `json:"id"`
	} `json:"owners"`
	ItemsPage ItemsPage `json:"items_page"`
}

type BoardGroups struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Groups []Group `json:"groups"`
}
type Group struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	EmailAddress string     `json:"email"`
	URL          string     `json:"url"`
	Admin        bool       `json:"is_admin"`
	Guest        bool       `json:"is_guest"`
	ViewOnly     bool       `json:"is_view_only"`
	Phone        string     `json:"phone"`
	Title        string     `json:"title"`
	Location     string     `json:"location"`
	CreatedAt    *time.Time `json:"created_at"`
}

type ItemsPage struct {
	Items  []Item `json:"items"`
	Cursor string `json:"cursor"`
}

type Column struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type CreateItemPayload struct {
	ItemName string `json:"item_name,omitempty"`
	BoardId  string `json:"board_id,omitempty"`
	GroupId  string `json:"group_id,omitempty"`
	Assignee string `json:"person"`
	Status   string `json:"status"`
}

type UpdateItemPayload struct {
	ID       string `json:"item_id,omitempty"`
	BoardId  string `json:"board_id,omitempty"`
	Assignee string `json:"person"`
	Status   string `json:"status"`
}

type CreateItemResponse struct {
	Data struct {
		CreateItem struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"create_item"`
	} `json:"data"`
	AccountID int `json:"account_id"`
}

type UpdateItemResponse struct {
	Data struct {
		UpdateItem struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"change_multiple_column_values"`
	} `json:"data"`
	AccountID int `json:"account_id"`
}

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Board struct {
		Name string `json:"name"`
	} `json:"board"`
	Creator struct {
		Name string `json:"name"`
	} `json:"creator"`
	Updates      []Update      `json:"updates"`
	Subscribers  []Subscriber  `json:"subscribers"`
	ColumnValues []ColumnValue `json:"column_values"`
	CreatedAt    *time.Time    `json:"created_at"`
	UpdatedAt    *time.Time    `json:"updated_at"`
}

type Update struct {
	ID      string `json:"id"`
	Creator struct {
		Name string `json:"name"`
	} `json:"creator"`
	Body string `json:"text_body"`
}

type Subscriber struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const (
	ColumnValueIDLink = "link"
)

type ColumnValue struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Value *string `json:"value"`
	Text  *string `json:"text"`
}

type ColumnValueValue struct {
	URL       string     `json:"url,omitempty"`  // "title":"Link", "id":"link",
	Text      string     `json:"text,omitempty"` // "title":"Link", "id":"link",
	ChangedAt *time.Time `json:"changed_at"`
}
