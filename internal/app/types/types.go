package types

type Items struct {
	Item []Item `json:"items"`;
}

type Item struct {
	Id string;
	Type string;
	Name string;
	Status string;
	EndDate string;
	StartDate string;
	AfterTask []string;
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
	Name string `toml:"name"`;
	Colors []string `toml:"colors"`;
}


