package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/manifest-it/go-monday"
)

const (
	EnvVarMondayToken   = "MONDAY_TOKEN"
	EnvVarMondayBoardId = "MONDAY_BOARD_ID"
)

func main() {
	loaded, err := config.LoadDotEnv([]string{".env", os.Getenv("ENV_PATH")}, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loaded [%s]\n", strings.Join(loaded, ","))

	tok := os.Getenv(EnvVarMondayToken)
	ids := []string{os.Getenv(EnvVarMondayBoardId)}

	cl := monday.NewClient(tok)
	resp, data, err := cl.GetItemsBetween(&ids, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("STATUS [%d]\n", resp.StatusCode)

	fmtutil.PrintJSON(data)

	for _, item := range data {
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(item)
	}

	fmt.Println("DONE")
}
