package output

import (
	"fmt"
)

type ResultOutput struct {
	Added []string
	Deled []string
	All   []string
	Code  map[string]int
}

func OutResult(results map[string]ResultOutput) error {
	//@title OutResult
	//@param
	//Return
	for domain, info := range results {
		fmt.Println("\n[*]Domain:", domain)
		for _, added := range info.Added {
			fmt.Printf("\n   [+]Added: %v", added)
		}
		for _, deled := range info.Deled {
			fmt.Printf("\n   [-]Deled: %v", deled)
		}
		fmt.Println("\n   [MORE]Alive Check Result:")
		for subd, code := range info.Code {
			if code != -1 {
				fmt.Printf("      - %v ( %v )\n", subd, code)
			}
		}

	}
	return nil
}
