package page_renderer

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"
)

func getEventStoreAndLastTimeAndNow() (event_store.Store, time.Time, time.Time) {
	time1 := time.Time{}
	time2 := time1.Add(time.Hour * 24)
	time3 := time2.Add(time.Hour * 24)
	time4 := time1.Add(time.Minute * 24)
	time5 := time2.Add(time.Minute * 24)
	time6 := time3.Add(time.Minute * 24)

	event1 := event_store.NewEvent(time1, "camera1", "image1-hi", "image1-lo", "video1-hi", "video1-lo")
	event2 := event_store.NewEvent(time2, "camera2", "image2-hi", "image2-lo", "video2-hi", "video2-lo")
	event3 := event_store.NewEvent(time3, "camera3", "image3-hi", "image3-lo", "video3-hi", "video3-lo")
	event4 := event_store.NewEvent(time4, "camera4", "image4-hi", "image4-lo", "video4-hi", "video4-lo")
	event5 := event_store.NewEvent(time5, "camera5", "image5-hi", "image5-lo", "video5-hi", "video5-lo")
	event6 := event_store.NewEvent(time6, "camera6", "image6-hi", "image6-lo", "video6-hi", "video6-lo")

	event1.EventID = uuid.UUID{0}
	event2.EventID = uuid.UUID{1}
	event3.EventID = uuid.UUID{2}
	event4.EventID = uuid.UUID{3}
	event5.EventID = uuid.UUID{4}
	event6.EventID = uuid.UUID{5}

	dir, err := ioutil.TempDir("", "eds-cctv-system")
	if err != nil {
		log.Fatalf("during test: %v", err)
	}

	path := filepath.Join(dir, "some_file.jsonl")

	store := event_store.NewStore(path)

	store.Add(event1)
	store.Add(event2)
	store.Add(event3)
	store.Add(event4)
	store.Add(event5)
	store.Add(event6)

	return store, time3, time3.Add(time.Hour * 24)
}

func TestRenderSummary(t *testing.T) {
	store, _, now := getEventStoreAndLastTimeAndNow()

	data, err := RenderSummary("All events", store.GetAllByDate(), now)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(data)

	assert.Equal(
		t,
		"</html>\n<head>\n<title>All events as at 0001-01-04 00:00:00 +0000 UTC</title>\n<style type=\"text/css\">\nBODY {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n}\n\nTH {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: bold;\n    text-align: center;\n}\n\nTD {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n    border: 1px solid gray; \n}\n</style>\n</head>\n\n<body>\n<h2>All events as at 0001-01-04 00:00:00 +0000 UTC</h2>\n\n<center>\n<table width=\"90%\">\n\n\t<tr>\n\t\t<th>Date</th>\n\t\t<th>Events</th>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_03.html\">0001-01-03</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_02.html\">0001-01-02</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n\t<tr>\n\t\t<td><a target=\"event\" href=\"events_0001_01_01.html\">0001-01-01</a></td>\n\t\t<td>2</td>\n\t</tr>\n\n</table>\n<center>\n\n</body>\n</html>",
		data,
	)
}

func TestRenderPage(t *testing.T) {
	store, time3, now := getEventStoreAndLastTimeAndNow()

	data, err := RenderPage("Events", store.GetAllByDate()[time3], time3, now)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(data)

	assert.Equal(
		t,
		"</html>\n<head>\n<title>Events for 0001-01-03 as at 0001-01-04 00:00:00 +0000 UTC</title>\n<style type=\"text/css\">\nBODY {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n}\n\nTH {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: bold;\n    text-align: center;\n}\n\nTD {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n    border: 1px solid gray; \n}\n</style>\n</head>\n\n<body>\n<h1>Events for 0001-01-03 as at 0001-01-04 00:00:00 +0000 UTC</h1>\n\n<center>\n<table width=\"90%\">\n\n\t<tr>\n\t\t<th>Event ID</th>\n\t\t<th>Timestamp</th>\n\t\t<th>Camera</th>\n\t\t<th>Screenshot</th>\n\t\t<th>Download</th>\n\t</tr>\n\n\t<tr>\n\t\t<td>05000000-0000-0000-0000-000000000000</td>\n\t\t<td>0001-01-03 00:24:00 +0000 UTC</td>\n\t\t<td>camera6</td>\n\t\t<td style=\"width: 320px\";><a target=\"_blank\" href=\"image6-hi\"><img src=\"image6-lo\" width=\"320\" height=\"180\" /></a></td>\n\t\t<td>Download <a target=\"_blank\" href=\"video6-hi\">high-res</a> or <a target=\"_blank\" href=\"video6-lo\">low-res</a></td>\n\t</tr>\n\n\t<tr>\n\t\t<td>02000000-0000-0000-0000-000000000000</td>\n\t\t<td>0001-01-03 00:00:00 +0000 UTC</td>\n\t\t<td>camera3</td>\n\t\t<td style=\"width: 320px\";><a target=\"_blank\" href=\"image3-hi\"><img src=\"image3-lo\" width=\"320\" height=\"180\" /></a></td>\n\t\t<td>Download <a target=\"_blank\" href=\"video3-hi\">high-res</a> or <a target=\"_blank\" href=\"video3-lo\">low-res</a></td>\n\t</tr>\n\n</table>\n<center>\n\n</body>\n</html>",
		data,
	)
}
