package motion_config

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	CameraName string
	CameraId   int64
	NetCamURL  string
	Width      int64
	Height     int64
}

func ParseFile(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := Config{}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		words := strings.Split(strings.TrimSpace(line), " ")
		if len(words) == 0 {
			continue
		}

		key := words[0]
		value := strings.Join(words[1:], " ")

		if key == "camera_name" {
			config.CameraName = value
		} else if key == "camera_id" {
			config.CameraId, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return Config{}, err
			}
		} else if key == "netcam_url" {
			config.NetCamURL = value
		} else if key == "width" {
			config.Width, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return Config{}, err
			}
		} else if key == "height" {
			config.Height, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return Config{}, err
			}
		}
	}

	return config, nil
}

func MergeConfigs(leftConfig, rightConfig Config) Config {
	config := leftConfig

	if rightConfig.CameraName != "" {
		config.CameraName = rightConfig.CameraName
	}

	if rightConfig.CameraId != 0 {
		config.CameraId = rightConfig.CameraId
	}

	if rightConfig.NetCamURL != "" {
		config.NetCamURL = rightConfig.NetCamURL
	}

	if rightConfig.Width != 0 {
		config.Width = rightConfig.Width
	}

	if rightConfig.Height != 0 {
		config.Height = rightConfig.Height
	}

	log.Printf("left config = %+v", leftConfig)
	log.Printf("right config = %+v", rightConfig)
	log.Printf("merged config = %+v", config)

	return config
}
