package broadway_test

import (
	"fmt"

	"github.com/DAddYE/broadway"
	"github.com/DAddYE/broadway/p"
)

func printFiles() broadway.Fn {
	return func(app *broadway.App) error {
		for _, f := range app.Files {
			fmt.Println(f.Path)
		}
		return nil
	}
}

func Example() {

	broadway.New(".").
		Use(p.Concat("dest.go", ".go")).
		Use(printFiles()).
		Build("./out")

	// Output:
	// README.md
	// dest.go
}
