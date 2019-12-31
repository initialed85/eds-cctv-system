package motion_config_segment_recorder

import (
	"github.com/initialed85/eds-cctv-system/internal/common"
	"github.com/initialed85/eds-cctv-system/internal/motion_config"
	"github.com/initialed85/eds-cctv-system/internal/segment_recorder"
	"log"
)

type Recorder struct {
	destinationPath string
	duration        int
	configs         []motion_config.Config
	processes       []*common.BackgroundProcess
	started         bool
}

func New(configPath, destinationPath string, duration int) (Recorder, error) {
	configs, err := motion_config.Find(configPath)
	if err != nil {
		return Recorder{}, err
	}

	baseConfig := configs[len(configs)-1]

	cameraConfigs := make([]motion_config.Config, 0)

	for _, config := range configs[0 : len(configs)-1] {
		cameraConfigs = append(cameraConfigs, motion_config.MergeConfigs(baseConfig, config))
	}

	return Recorder{
		destinationPath: destinationPath,
		duration:        duration,
		configs:         cameraConfigs,
		processes:       make([]*common.BackgroundProcess, 0),
		started:         false,
	}, nil
}

func (r *Recorder) Start() error {
	if r.started {
		return nil
	}

	for _, config := range r.configs {
		log.Printf("starting segment recorder for %+v", config)

		process, err := segment_recorder.RecordSegments(config.NetCamURL, r.destinationPath, config.CameraName, r.duration)
		if err != nil {
			return err
		}

		r.processes = append(r.processes, process)
	}

	r.started = true

	return nil
}

func (r *Recorder) Stop() {
	if !r.started {
		return
	}

	for _, process := range r.processes {
		log.Printf("starting segment recorder at %+v", process)

		process.Stop()
	}

	r.started = false
}
