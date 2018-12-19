package utils

import "time"

// StandardResponse contains the standard fields included in most responses.
type StandardResponse struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Operation  string `json:"operation"`
	ErrorCode  int    `json:"error_code"`
	Error      string `json:"error"`
}

// ListContainerResponse is the basic response returned by LXD.
type ListContainerResponse struct {
	StandardResponse
	Metadata []string `json:"metadata"`
}

// ListOperationsResponse contains the list of running background operations.
type ListOperationsResponse struct {
	StandardResponse
	Metadata struct {
		Running []string `json:"running"`
	} `json:"metadata"`
}

// SnapshotRequest is the JSON body containing snapshot information.
type SnapshotRequest struct {
	Name     string `json:"name"`
	Stateful bool   `json:"stateful"`
}

// SnapshotResponse is the JSON response containing background operation information.
type SnapshotResponse struct {
	StandardResponse
	Metadata struct {
		ID          string    `json:"id"`
		Class       string    `json:"class"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Status      string    `json:"status"`
		StatusCode  int       `json:"status_code"`
		Resources   struct {
			Containers []string `json:"containers"`
		} `json:"resources"`
		Metadata  interface{} `json:"metadata"`
		MayCancel bool        `json:"may_cancel"`
		Err       string      `json:"err"`
	} `json:"metadata"`
}

// SnapshotBackgroundOperation contains the background information of the snapshot job.
type SnapshotBackgroundOperation struct {
	ID          string    `json:"id"`
	Class       string    `json:"class"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	StatusCode  int       `json:"status_code"`
	Resources   struct {
		Containers []string `json:"containers"`
	} `json:"resources"`
	Metadata  interface{} `json:"metadata"`
	MayCancel bool        `json:"may_cancel"`
	Err       string      `json:"err"`
}

// PublishRequest contains information to publish an image.
type PublishRequest struct {
	Filename string                  `json:"filename"`
	Public   bool                    `json:"public"`
	Aliases  []PublishRequestAliases `json:"aliases"`
	Source   PublishRequestSource    `json:"source"`
}

// PublishRequestAliases contains the aliases of an image.
type PublishRequestAliases struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PublishRequestSource contains the source of an image.
type PublishRequestSource struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ImagePropertiesResponse contains the properties of an image.
type ImagePropertiesResponse struct {
	StandardResponse
	Metadata struct {
		AutoUpdate bool `json:"auto_update"`
		Properties struct {
		} `json:"properties"`
		Public  bool `json:"public"`
		Aliases []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"aliases"`
		Architecture string    `json:"architecture"`
		Cached       bool      `json:"cached"`
		Filename     string    `json:"filename"`
		Fingerprint  string    `json:"fingerprint"`
		Size         int       `json:"size"`
		CreatedAt    time.Time `json:"created_at"`
		ExpiresAt    time.Time `json:"expires_at"`
		LastUsedAt   time.Time `json:"last_used_at"`
		UploadedAt   time.Time `json:"uploaded_at"`
	} `json:"metadata"`
}

// SnapshotNamingScheme contains the container name and snapshot name.
type SnapshotNamingScheme struct {
	ContainerName string
	SnapshotName  string
}
