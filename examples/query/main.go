package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/grokify/go-monday"
	"github.com/grokify/go-monday/simpleitem"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
)

func main() {
	loaded, err := config.LoadDotEnv(".env", os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loaded [%s]\n", strings.Join(loaded, ","))

	tok := os.Getenv("MONDAY_TOKEN")

	gql := monday.BoardQuery(os.Getenv("MONDAY_BOARD_ID"))

	cl := monday.NewClient(tok)
	resp, err := cl.DoGraphQL(gql)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("STATUS [%d]\n", resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	var brds monday.Response

	err = json.Unmarshal(data, &brds)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(brds)

	for _, b := range brds.Data.Boards {
		si, err := simpleitem.BoardSimpleItems(b)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(si)
	}

	fmt.Println("DONE")
}
