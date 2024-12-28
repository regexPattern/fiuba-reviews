package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	bytes, _ := io.ReadAll(os.Stdin)
	infoSiu := scrapearSiu(string(bytes))

	bytes, err := json.Marshal(infoSiu)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
}
