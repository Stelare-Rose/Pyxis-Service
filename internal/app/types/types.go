package types

type Item struct {
	Id string;
	Type string;
	Name string;
	Status string;
	EndDate string;
	StartDate string;
	PriorityDate string;
	CompletedDate string;
	AfterTask []string;
	Tags []string;
	Path string;
	Fingerprint int64;
	IsArchived bool;
	IsVerified bool;
}

type Idea struct {
	Id string;
	Name string;
	Status string;
	CreatedDate string;
	PriorityDate string;
	CompletedDate string;
	Tags []string;
	Path string;
	Fingerprint int64;
	IsArchived bool;
	IsVerified bool;
}
type Tags struct {
	Tags map[string]Tag `toml:"tags"`;
}

type Tag struct {
	Tag string `toml:"tag"`;
	Color []string `toml:"color"`;
}


