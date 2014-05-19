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

func printOutput() broadway.Fn {
	return func(app *broadway.App) error {
		for _, f := range app.Files {
			fmt.Printf("%s: %s\n", f.Path, f.Contents)
		}
		return nil
	}
}

func Example() {

	broadway.New("./test").
		Use(p.Matilde("*.js")).
		// Use(p.Concat("dest.go", ".go")).
		Use(printOutput()).
		Build("./test/out")

	// Output:
	// README.md
	// dest.go
}
