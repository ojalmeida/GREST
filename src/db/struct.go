package db

type Behavior struct {
	PathMapping PathMapping
	KeyMappings []KeyMapping
}

type KeyMapping struct {
	Key    string
	Column string
}

type PathMapping struct {
	Path  string
	Table string
}
