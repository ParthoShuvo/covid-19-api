package db

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
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
		log.Fatalf("failed to open file at %s", filePath)
		return nil, err
	}
	log.Printf("file: %s opened successfully", file.Name())
	return parser.parse(file)
}

func loadFile(filePath string) (*os.File, error) {
	exec, _ := os.Executable()
	return os.Open(path.Join(path.Dir(exec), filePath))
}
