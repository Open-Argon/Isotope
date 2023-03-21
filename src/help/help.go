package help

import (
	"fmt"
)

type Options [][2]string

func Help(usage string, choices Options) {
	fmt.Println()
	fmt.Println("Usage:", usage)
	fmt.Println()
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	longest := 0
	for _, choice := range choices {
		if len(choice[0]) > longest {
			longest = len(choice[0])
		}
	}
	for _, choice := range choices {
		fmt.Printf(" %s", choice[0])
		for i := len(choice[0]); i < longest; i++ {
			fmt.Print(" ")
		}
		fmt.Printf(" : %s", choice[1])
		fmt.Println()
	}
	fmt.Println()
}
