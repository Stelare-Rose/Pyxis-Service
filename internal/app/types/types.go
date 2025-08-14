package types

type Items struct {
	Item []Item `json:"items"`;
}

type Item struct {
	Type string `json:"type"`;
	Id string `json:"id"`;
	Name string `json:"name"`;
	Status string `json:"status"`;
	HardDeadline string `json:"hardDeadline"`;
	AfterTasks []string `json:"AfterTasks"`;
	Tags []string `json:"tags"`;
	Path string `json:"path"`;
	Hash uint64 `json:"hash"`;
}
