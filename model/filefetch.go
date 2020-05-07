package model

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/errors"
)

type parser interface {
	parse(file *os.File) (interface{}, error)
}

type defaultParser struct{}

func (d defaultParser) parse() (interface{}, error) {
	return nil, nil
}

func fetch(filePath string, parser parser) (interface{}, error) {
	file, err := loadFile(filePath)
	if err != nil {
		return nil, errors.InternalServerError.Wrapf(err, "failed to open file: %s", filePath)
	}
	log.Printf("file: %s opened successfully", file.Name())
	return parser.parse(file)
}

func loadFile(filePath string) (*os.File, error) {
	exec, _ := os.Executable()
	return os.Open(path.Join(path.Dir(exec), filePath))
}

func listFiles(dir string) []string {
	files := []string{}
	var walkFn filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.InternalServerError.Wrapf(err, "path=%s walking failed", path)
		}
		if filepath.Ext(path) != ".csv" {
			log.Printf("file=%s is not .csv", info.Name())
			return nil
		}
		files = append(files, path)
		return nil
	}
	filepath.Walk(dir, walkFn)
	return files
}
