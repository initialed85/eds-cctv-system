package motion_config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Find(path string) ([]Config, error) {
	var configs []Config

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_, file := filepath.Split(path)
		if !strings.HasSuffix(file, ".conf") {
			return nil
		}

		log.Printf("found config at %v", path)

		config, err := ParseFile(path)
		if err != nil {
			log.Printf("failed to parse because %v", err)

			return err
		}

		log.Printf("relevant contents are %+v", config)

		configs = append(configs, config)

		return nil
	})

	if err != nil {
		return []Config{}, err
	}

	return configs, nil
}
