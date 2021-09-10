package monday

import (
	"strconv"
	"testing"
)

var requestTests = []struct {
	boardId int
	gql     string
}{
	{12345, "query {boards (ids:12345) {name state columns {id title type} owner {id} items {id name state column_values {title id value text}}}}"},
}

func TestRequest(t *testing.T) {
	for _, tt := range requestTests {
		gqlTry := BoardQuery(strconv.Itoa(tt.boardId))
		if tt.gql != gqlTry.String() {
			t.Errorf("monday.BoardQuery(%d) Want: [%s] Got: [%s]",
				tt.boardId, tt.gql, gqlTry)
		}
	}
}
