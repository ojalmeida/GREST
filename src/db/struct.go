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

/*
	Profile User definition.
	Perm Model:
	1 Read
	2 Create
	4 Edit/Delete

	1+2   = 3 Read/Create
	1+2+4 = 7 ALL
 */
type Profile struct {
	uid      	int    // Not needed
	Username 	string
	Password 	string
	Perm		int
	Token		string

	Authenticated bool
}