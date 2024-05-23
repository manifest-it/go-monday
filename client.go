package monday

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (c *Client) GetItemsBetween(boardID string, startTime, endTime time.Time, limit int) (*http.Response, *ItemsPage, error) {
	q := GetItemsQuery(boardID, startTime, endTime, limit)
	log.Println("GetItemsBetween GQL query:", q)

	resp, err := c.DoGraphQL(q)
	if err != nil {
		return nil, nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var boardsItems ItemsResponse
	err = json.Unmarshal(data, &boardsItems)
	if err != nil {
		return nil, nil, err
	}

	return resp, &boardsItems.Data.Boards[0].ItemsPage, nil
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
