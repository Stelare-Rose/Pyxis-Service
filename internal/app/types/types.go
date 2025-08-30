package types

type Items struct {
	Item []Item `json:"items"`;
}

type Item struct {
	Id string `json:"id"`;
	Type string `json:"type"`;
	Name string `json:"name"`;
	Status string `json:"status"`;
	HardDeadline string `json:"hardDeadline"`;
	SoftDeadline string;
	AfterTask []string `json:"afterTask"`;
	Tags []string `json:"tags"`;
	Path string `json:"path"`;
	Fingerprint int64 `json:"fingerprint"`;
	IsArchived bool;
	IsVerified bool;
}

type Tags struct {
	Tags map[string]Tag `toml:"tags"`;
}

type Tag struct {
	Colors []string `toml:"colors"`;
}


