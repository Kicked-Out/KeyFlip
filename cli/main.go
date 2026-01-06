package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Kicked-Out/KeyFlip/core"
)

func main() {
	from := flag.String("from", "en", "source layout (en|ua)")
	to := flag.String("to", "ua", "target layout (en|ua)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: keyflip --from en --to ua \"text\"")
		os.Exit(1)
	}

	input := flag.Arg(0)

	var mapping map[rune]rune
	switch {
	case *from == "en" && *to == "ua":
		mapping = core.EnToUa
	case *from == "ua" && *to == "en":
		mapping = core.UaToEn
	default:
		fmt.Println("Unsupported layout combination")
		os.Exit(1)
	}

	fmt.Println(core.Transform(input, mapping))
}
