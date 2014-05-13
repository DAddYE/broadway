package p

import (
	"path"

	"github.com/DAddYE/broadway"
)

// Giving a folder like:
// a.css	b.css	d.css
// broadway.New("./blog").Use(concat("app.css", ".css"))
// will generate in ./blog, one file app.css with inside the content of a,b,c.css
func Concat(filepath string, ext string) broadway.Fn {

	return func(b *broadway.App) error {
		dest := broadway.File{}
		dest.Path = filepath
		dup := []*broadway.File(nil)

		for _, file := range b.Files {

			// If we the ext is different appen the file to dups
			if path.Ext(file.Path) != ext {
				dup = append(dup, file)
				continue
			}

			// Concat the content
			dest.Contents = append(dest.Contents, file.Contents...)

			// Set the file mode if still empty
			if dest.Mode == 0 {
				dest.Mode = file.Mode
			}
		}

		// Basically if the mode is still 0 means we did not find any file with the same ext
		if dest.Mode != 0 {
			b.Files = append(dup, &dest)
		}

		return nil
	}
}
