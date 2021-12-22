package segment

import (
	"fmt"

	"github.com/platform9/appctl/pkg/constants"
	"gopkg.in/segmentio/analytics-go.v3"
)

var APPCTL_SEGMENT_WRITE_KEY = "uDPDiaRE8jHI6NJKsQsXYWFyNGyw5iZj"

func SegmentClient() (analytics.Client, error) {
	client := analytics.New(APPCTL_SEGMENT_WRITE_KEY)
	return client, nil
}

func SendGroupTraits(c analytics.Client, id string, data map[string]interface{}) error {
	userID := fmt.Sprintf("appctl-%s", id)

	//fmt.Printf("Sending group traits to segment with groupID: %s, userID: %s", id, userID)

	if err := c.Enqueue(analytics.Group{
		UserId:  userID,
		GroupId: id,
		Traits:  data,
	}); err != nil {
		return err
	}

	return nil
}

//Sent event for appctl commands specific to an app.
func SendEvent(c analytics.Client, name string, id string, status string, loginType string, errMessage string, data []constants.ListAppInfo) error {
	userID := fmt.Sprintf("appctl-%s", id)
	var data_str constants.ListAppInfo
	if data != nil {
		data_str = data[0]
	}

	// Should be as a log message.
	//fmt.Printf("Sending event to segment with title: %s, userID: %s\n", name, userID)
	if err := c.Enqueue(analytics.Track{
		UserId: userID,
		Event:  name,
		Properties: analytics.NewProperties().
			Set("appname", data_str.Name).
			Set("image", data_str.Image).
			Set("url", data_str.URL).
			Set("port", data_str.Port).
			Set("appstatus", data_str.ReadyStatus).
			Set("id", id).
			Set("logintype", loginType).
			Set("eventstatus", status).
			Set("error", errMessage),
	}); err != nil {
		return err
	}
	return nil
}
func SendEventList(c analytics.Client, name string, id string, status string, loginType string, errMessage string, data interface{}) error {
	userID := fmt.Sprintf("appctl-%s", id)
	// Should be as a log message.
	//fmt.Printf("Sending event to segment with title: %s, userID: %s\n", name, userID)
	if err := c.Enqueue(analytics.Track{
		UserId: userID,
		Event:  name,
		Properties: analytics.NewProperties().
			Set("appinfo", data).
			Set("id", id).
			Set("logintype", loginType).
			Set("status", status).
			Set("error", errMessage),
	}); err != nil {
		return err
	}

	return nil
}

// To close the event.
func Close(c analytics.Client) {
	c.Close()
}
