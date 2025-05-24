package github

type Event struct {
	HeadCommit HeadCommit `json:"head_commit"`
}

type HeadCommit struct {
	Message string `json:"message"`
}
