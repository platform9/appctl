package segment

import (
	"fmt"

	"gopkg.in/segmentio/analytics-go.v3"
)

var APPCTL_SEGMENT_WRITE_KEY = "***REMOVED***"

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

func SendEvent(c analytics.Client, name string, id string, status string, data interface{}) error {
	userID := fmt.Sprintf("appctl-%s", id)
	// Should be as a log message.
	//fmt.Printf("Sending event to segment with title: %s, userID: %s\n", name, userID)
	if err := c.Enqueue(analytics.Track{
		UserId: userID,
		Event:  name,
		Properties: analytics.NewProperties().
			Set("AppInfo", data).
			Set("Status", status),
	}); err != nil {
		return err
	}

	return nil
}

// To close the event.
func Close(c analytics.Client) {
	c.Close()
}
