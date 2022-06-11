package stream

import (
	"time"
)

type Message struct {
	Stream    string `json:"stream"`
	Partition string `json:"partition"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Offset    string `json:"offset"`
	Timestamp string `json:"timestamp"`
}

type Events struct {
	EventType          string    `json:"eventType"`
	CloudEventsVersion string    `json:"cloudEventsVersion"`
	EventTypeVersion   string    `json:"eventTypeVersion"`
	Source             string    `json:"source"`
	EventTime          time.Time `json:"eventTime"`
	ContentType        string    `json:"contentType"`
	Data               struct {
		CompartmentID      string `json:"compartmentId"`
		CompartmentName    string `json:"compartmentName"`
		ResourceName       string `json:"resourceName"`
		ResourceID         string `json:"resourceId"`
		AvailabilityDomain string `json:"availabilityDomain"`
		AdditionalDetails  struct {
			BucketName    string `json:"bucketName"`
			VersionID     string `json:"versionId"`
			ArchivalState string `json:"archivalState"`
			Namespace     string `json:"namespace"`
			BucketID      string `json:"bucketId"`
			ETag          string `json:"eTag"`
		} `json:"additionalDetails"`
	} `json:"data"`
	EventID    string `json:"eventID"`
	Extensions struct {
		CompartmentID string `json:"compartmentId"`
	} `json:"extensions"`
}
