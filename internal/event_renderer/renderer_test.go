package event_renderer

import (
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func TestRenderEventsSummary(t *testing.T) {
	time1 := time.Time{}
	time2 := time1.Add(time.Hour * 24)
	time3 := time2.Add(time.Hour * 24)

	event1 := event_store.NewEvent(time1, "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := event_store.NewEvent(time1, "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := event_store.NewEvent(time2, "image3-hi", "image3-lo", "video3-hi", "video3-lo")
	event4 := event_store.NewEvent(time2, "image4-hi", "image4-lo", "video4-hi", "video4-lo")
	event5 := event_store.NewEvent(time3, "image5-hi", "image5-lo", "video5-hi", "video5-lo")
	event6 := event_store.NewEvent(time3, "image6-hi", "image6-lo", "video6-hi", "video6-lo")

	dir, err := ioutil.TempDir("", "watcher_test")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.txt")

	collection := event_store.NewStore(path)

	collection.Add(event1)
	collection.Add(event2)
	collection.Add(event3)
	collection.Add(event4)
	collection.Add(event5)
	collection.Add(event6)

	now := time3.Add(time.Hour * 24)

	data, err := RenderEventsSummary(collection.GetAllByDate(), now)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(
		t,
		"</html>\n<head>\n<title>All events as at 0001-01-04 00:00:00 +0000 UTC</title>\n<style type=\"text/css\">\nBODY {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n}\n\nTH {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: bold;\n    text-align: center;\n}\n\nTD {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n    border: 1px solid gray; \n}\n</style>\n</head>\n\n<body>\n<h2>All events as at 0001-01-04 00:00:00 +0000 UTC</h2>\n\n<center>\n<table width=\"90%\">\n\n\t<tr>\n\t\t<th>Date</th>\n\t\t<th>Events</th>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_01.html\">0001-01-01</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_02.html\">0001-01-02</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_03.html\">0001-01-03</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n</table>\n<center>\n\n</body>\n</html>",
		data,
	)
}
