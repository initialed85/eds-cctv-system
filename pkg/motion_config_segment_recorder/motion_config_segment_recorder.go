package motion_config_segment_recorder

import (
	"eds-cctv-system/internal/motion_config"
	"eds-cctv-system/internal/segment_recorder"
	"log"
	"os"
)

type MotionConfigSegmentRecorder struct {
	destinationPath string
	duration        int
	configs         []motion_config.Config
	processes       []*os.Process
	started         bool
}

func New(configPath, destinationPath string, duration int) (MotionConfigSegmentRecorder, error) {
	configs, err := motion_config.Find(configPath)
	if err != nil {
		return MotionConfigSegmentRecorder{}, err
	}

	baseConfig := configs[len(configs)-1]

	cameraConfigs := make([]motion_config.Config, 0)

	for _, config := range configs[0 : len(configs)-1] {
		cameraConfigs = append(cameraConfigs, motion_config.MergeConfigs(baseConfig, config))
	}

	return MotionConfigSegmentRecorder{
		destinationPath: destinationPath,
		duration:        duration,
		configs:         cameraConfigs,
		processes:       make([]*os.Process, 0),
		started:         false,
	}, nil
}

func (m *MotionConfigSegmentRecorder) Start() error {
	if m.started {
		return nil
	}

	for _, config := range m.configs {
		log.Printf("starting segment recorder for %+v", config)

		process, err := segment_recorder.RecordSegments(config.NetCamURL, m.destinationPath, config.CameraName, m.duration)
		if err != nil {
			return err
		}

		m.processes = append(m.processes, process)
	}

	m.started = true

	return nil
}

func (m *MotionConfigSegmentRecorder) Stop() {
	if !m.started {
		return
	}

	for _, process := range m.processes {
		log.Printf("starting segment recorder at %+v", process)

		_ = process.Kill()

		_, _ = process.Wait()
	}

	m.started = false
}
