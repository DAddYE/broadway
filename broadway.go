package broadway

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Example:
//
//   New("./my/static/site").
//		Use(Coffee("assets/**/*.coffee")).
//		Use(Sass("assets/**/*.{scss,sass}")).
//		Use(Sprockets("/assets/{js,css,vendor}")).
//		Use(Markdown("*.{md,markdown}")).
//		Build()
//

type Fn func(b *App) error

type File struct {
	Path     string
	Contents []byte
	Attrs    FileAttrs
	Mode     os.FileMode
}

type FileAttrs map[string]string

type App struct {
	Dir   string  // working directory
	Dest  string  // destination directory
	Files []*File // Files in the meet during the walk
	fns   []Fn    // list o functions to apply
}

// Initializes a new Broadway instance
func New(dir string) (b *App) {
	b = &App{}
	b.Dir = dir

	return b
}

func (b *App) Use(fn Fn) *App {
	b.fns = append(b.fns, fn)
	return b
}

func (b *App) Walk(ch chan *File) {
	filepath.Walk(b.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directory
		if info.IsDir() {
			return nil
		}

		buff, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		ch <- &File{path, buff, make(FileAttrs), info.Mode()}

		return nil

	})
	close(ch)
}

func (b *App) Build(dest string) error {
	b.Dest = dest

	ch := make(chan *File)
	go b.Walk(ch)

	for file := range ch {
		b.Files = append(b.Files, file)
	}

	for _, f := range b.fns {
		err := f(b)
		if err != nil {
			return err
		}
	}

	return nil
}
