package db

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
)


var requiredPathMappings = []PathMapping{
	{Path: "/config/path-mappings", Table: "path_mappings"},
	{Path: "/config/key-mappings", Table: "key_mappings"},
	{Path: "/config/behaviors", Table: "behaviors"},
}
var requiredKeyMappings = []KeyMapping{
	{Key: "id", Column: "path_mapping_id"},
	{Key: "id", Column: "key_mapping_id"},
	{Key: "id", Column: "behavior_id"},
	{Key: "key", Column: "_key"},
	{Key: "column", Column: "_column"},
	{Key: "path", Column: "path"},
	{Key: "table", Column: "_table"},
	{Key: "path_mapping_id", Column: "path_mapping_id"},
	{Key: "key_mapping_id", Column: "key_mapping_id"},
}
var requiredBehaviors = []Behavior{
	{PathMapping: PathMapping{Path: "/config/path-mappings", Table: "path_mappings"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "path_mapping_id"},
		{Key: "path", Column: "path"},
		{Key: "table", Column: "_table"},
	},
	},
	{PathMapping: PathMapping{Path: "/config/key-mappings", Table: "key_mappings"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "key_mapping_id"},
		{Key: "key", Column: "_key"},
		{Key: "column", Column: "_column"},
	},
	},
	{PathMapping: PathMapping{Path: "/config/behaviors", Table: "behaviors"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "behavior_id"},
		{Key: "path_mapping_id", Column: "path_mapping_id"},
		{Key: "key_mapping_id", Column: "key_mapping_id"},
	},
	},
}

var LocalConn *sqlx.DB

func init() {
	// Remote connection.
	LocalConn = LocalDB()
}
