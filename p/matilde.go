package p

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/DAddYE/broadway"
)

var (
	MATILDE_DIRECTIVE = `(?m)^\W*=\s*import\s*"(.*)"\s*$`
	MATILDE_RE        = regexp.MustCompile(MATILDE_DIRECTIVE)
)

/*
	Process deps of a file like:

				A
			   / \
			  B   C
			 /
			C

	A depends on: B and C
	B depends on: C

	Parsing B, means we have already parsed C as well.


*/
func Matilde(filetypes string, paths ...string) broadway.Fn {
	return func(b *broadway.App) error {
		paths = append(paths, b.Dir)
		conf := &Config{b, paths, make(map[*broadway.File][]*broadway.File, 1)}

		for _, file := range b.Files {
			matched, err := filepath.Match(filetypes, file.Name)
			if !matched || err != nil {
				continue
			}
			resolve(conf, file)
		}

		// If print dependency tree
		// for k, v := range conf.deps {
		// 	fmt.Println(k.Path)
		// 	for _, d := range v {
		// 		fmt.Printf("\t%s\n", d.Path)
		// 	}
		// }

		return nil
	}
}

type Config struct {
	app   *broadway.App
	paths []string
	deps  map[*broadway.File][]*broadway.File
}

func resolve(conf *Config, file *broadway.File) {
	file.Contents = MATILDE_RE.ReplaceAllFunc(file.Contents, func(directive []byte) []byte {
		match := MATILDE_RE.FindStringSubmatch(string(directive))
		name := string(match[1]) // dependency name

		// Add the extension of the parent if not present
		if filepath.Ext(name) == "" {
			name = name + file.Ext
		}

		// Compute the correct filename (if any)
		filename, stat := fileExistsIn(name, conf.paths...)
		if stat == nil {
			fmt.Errorf("Dependency `%s` of file `%s` not found!", name, file.Path)
			return nil
		}

		// Try to see if the dependency is tracked by broadway
		found := conf.app.FindByStat(stat)
		if found == nil {
			buff, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Errorf("Error reading file: %s", err)
				return nil
			}

			found = &broadway.File{
				Path:     filename,
				Name:     stat.Name(),
				Dir:      filepath.Dir(filename),
				Ext:      filepath.Ext(filename),
				Contents: buff,
				Mode:     stat.Mode(),
				Stat:     stat,
			}
		}

		// If we already resolved it
		if contains(conf.deps[file], found) > -1 {
			return nil
		}

		// If we already resolved it (recursive)
		for _, dep := range conf.deps[file] {
			if contains(conf.deps[dep], found) > -1 {
				return nil
			}
		}

		// Track dependency of the file
		conf.deps[file] = append(conf.deps[file], found)

		// Check circular dependencies
		if contains(conf.deps[file], file) > -1 {
			fmt.Errorf("Circular dependency detected between %s and %s", file.Path, found.Path)
			return nil
		}

		// Check again if we have new dependencies
		resolve(conf, found)

		// Remove the dependency from broadway.Files.
		if i := contains(conf.app.Files, found); i > -1 {
			conf.app.Files = append(conf.app.Files[:i], conf.app.Files[i+1:]...)
		}

		// Now we can return back the result (if any)
		return found.Contents
	})
}

// Helpers
func fileExistsIn(file string, paths ...string) (string, os.FileInfo) {
	for _, path := range paths {
		fullpath := filepath.Join(path, file)
		stat, err := os.Stat(fullpath)
		if err == nil {
			return fullpath, stat
		}
	}
	return "", nil
}

func contains(stack []*broadway.File, file *broadway.File) int {
	for i, item := range stack {
		if file.Path == item.Path {
			return i
		}
	}
	return -1
}
