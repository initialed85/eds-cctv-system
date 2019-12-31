package duration_finder

import (
	"fmt"
	"github.com/initialed85/eds-cctv-system/internal/common"
	"strings"
	"time"
)

func FindDuration(path string) (time.Duration, error) {
	_, stderr, err := common.RunCommand(
		"ffprobe",
		"-sexagesimal",
		path,
	)

	if err != nil {
		return time.Duration(0), fmt.Errorf("%v; stderr=%v", err, stderr)
	}

	for _, line := range strings.Split(stderr, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "Duration: ") {
			continue
		}

		durationString := strings.Split(strings.Split(line, "Duration: ")[1], ",")[0]

		durationParts := strings.Split(durationString, ":")
		if len(durationParts) != 3 {
			return time.Duration(0), fmt.Errorf("expected 3 colon-separated parts; instead had %v", durationParts)
		}

		duration, err := time.ParseDuration(
			fmt.Sprintf("%vh%vm%vs", durationParts[0], durationParts[1], durationParts[2]),
		)

		if err != nil {
			return time.Duration(0), fmt.Errorf("failed to build duration from %v because: %v", duration, err)
		}

		return duration, err
	}

	return time.Duration(0), fmt.Errorf("failed to find duration in %v", stderr)
}
