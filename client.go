package monday

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httputilmore"
)

const (
	MondayAPIURL = "https://api.monday.com/v2"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(token string) Client {
	return Client{httpClient: authutil.NewClientAuthzTokenSimple("", token)}
}

func (c *Client) DoJSON(data []byte) (*http.Response, error) {
	if c.httpClient == nil {
		return nil, errors.New("no auth token")
	}
	req, err := http.NewRequest(
		http.MethodPost,
		MondayAPIURL,
		bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)
	return c.httpClient.Do(req)
}

func (c *Client) DoGraphQLString(gql string) (*http.Response, error) {
	req := QueryRequest{Query: gql}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.DoJSON(data)
}

func (c *Client) DoGraphQL(gql Query) (*http.Response, error) {
	return c.DoGraphQLString(gql.String())
}

type QueryRequest struct {
	Query string `json:"query"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
}

func (c *Client) GetItemsBetween(boardID string, startTime, endTime time.Time, limit int) (*http.Response, *ItemsPage, error) {
	q := GetItemsQuery(boardID, startTime, endTime, limit)
	log.Println("GetItemsBetween ", startTime.Format("2006-01-02 15:04"), endTime.Format("2006-01-02 15:04"))

	resp, err := c.DoGraphQL(q)
	if err != nil {
		return nil, nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(data, &errResp); err == nil {
		if errResp.StatusCode == http.StatusTooManyRequests {
			return &http.Response{StatusCode: errResp.StatusCode}, nil, errors.New("too many requests")
		}
	}

	var boardsItems ItemsResponse
	err = json.Unmarshal(data, &boardsItems)
	if err != nil {
		return nil, nil, err
	}

	var itemsPage ItemsPage

	if len(boardsItems.Data.Boards) > 0 {
		itemsPage = boardsItems.Data.Boards[0].ItemsPage
	}

	return resp, &itemsPage, nil
}

func (c *Client) GetNextItems(cursor string, limit int) (*http.Response, *ItemsPage, error) {
	q := GetNextItemsQuery(cursor, limit)
	log.Println("GetNextItems GQL query:", q)

	resp, err := c.DoGraphQL(q)
	if err != nil {
		return nil, nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var NextItems NextItemsResponse
	err = json.Unmarshal(data, &NextItems)
	if err != nil {
		return nil, nil, err
	}

	return resp, &NextItems.Data.ItemsPage, nil
}

func (c *Client) GetAllUsers() (*http.Response, []User, error) {
	q := GetAllUsersQuery()
	log.Println("GetAllUsers GQL query:", q)

	resp, err := c.DoGraphQL(q)
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var usersResponse UsersResponse
	err = json.Unmarshal(data, &usersResponse)
	if err != nil {
		return nil, nil, err
	}

	return resp, usersResponse.Data.Users, nil
}

func (c *Client) GetAllBoards() (*http.Response, []BoardGroups, error) {
	q := GetAllBoardsQuery()
	log.Println("GetAllBoards GQL query:", q)

	resp, err := c.DoGraphQL(q)

	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var boardsResponse BoardsResponse
	err = json.Unmarshal(data, &boardsResponse)
	if err != nil {
		return nil, nil, err
	}

	return resp, boardsResponse.Data.Boards, nil
}

func (c *Client) GetUserById(userId string) (*http.Response, *User, error) {
	q := GetUserByIdQuery(userId)
	log.Println("GetUserById GQL query:", q)

	resp, err := c.DoGraphQL(q)
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var usersResponse UsersResponse
	err = json.Unmarshal(data, &usersResponse)
	if err != nil {
		return nil, nil, err
	}

	var user User
	if len(usersResponse.Data.Users) > 0 {
		user = usersResponse.Data.Users[0]
	}

	return resp, &user, nil
}

func (c *Client) CreateItem(item CreateItemPayload) (*http.Response, *CreateItemResponse, error) {
	boardId, groupId, itemName := item.BoardId, item.GroupId, item.ItemName
	item.BoardId, item.GroupId, item.ItemName = "", "", ""

	columnValuesJSON, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("Error marshalling column values: %v", err)
	}

	q := fmt.Sprintf(`
	mutation {
		create_item(
			board_id: "%s"
			group_id: "%s"
			item_name: "%s"
			column_values: %q
		) {
			id
			url
		}
	}`, boardId, groupId, itemName, columnValuesJSON)

	log.Println("CreateItem GQL query:", q)

	resp, err := c.DoGraphQLString(q)
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var itemResp CreateItemResponse
	err = json.Unmarshal(data, &itemResp)
	if err != nil {
		return nil, nil, err
	}

	return resp, &itemResp, nil
}

func (c *Client) UpdateItem(item UpdateItemPayload) (*http.Response, *UpdateItemResponse, error) {
	boardId, itemId := item.BoardId, item.ID
	item.BoardId, item.ID = "", ""

	columnValuesJSON, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("Error marshalling column values: %v", err)
	}

	q := fmt.Sprintf(`
	mutation {
		change_multiple_column_values(
			board_id: "%s"
			item_id: "%s"
			column_values: %q
		) {
			id
			url
		}
	}`, boardId, itemId, columnValuesJSON)

	log.Println("UpdateItem GQL query:", q)

	resp, err := c.DoGraphQLString(q)
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var itemResp UpdateItemResponse
	err = json.Unmarshal(data, &itemResp)
	if err != nil {
		return nil, nil, err
	}

	return resp, &itemResp, nil
}
