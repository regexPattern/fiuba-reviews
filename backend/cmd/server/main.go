package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/regexPattern/fiuba-reviews/pkg/scraper_siu"
)

func main() {
	bytes, _ := io.ReadAll(os.Stdin)
	infoSiu := scraper_siu.ScrapearSiu(string(bytes))

	bytes, err := json.Marshal(infoSiu)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
}
