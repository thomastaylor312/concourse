package atc

type WorkerArtifact struct {
	ID        int    `json:"id"`
	Path      string `json:"path"`
	BuildID   int    `json:"build_id"`
	CreatedAt int64  `json:"created_at"`
}
