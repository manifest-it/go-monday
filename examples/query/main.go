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

	cl := monday.NewClient(tok)
	resp, data, err := cl.GetItemsBetween(os.Getenv(EnvVarMondayBoardId), time.Now(), time.Now(), 25)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("STATUS [%d]\n", resp.StatusCode)

	fmtutil.PrintJSON(data)

	for _, item := range data.Items {
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(item)
	}

	fmt.Println("DONE")
}
