package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ojalmeida/GREST/src/db/connection"
)

var Conn *sqlx.DB

var requiredPathMappings = []PathMapping{
	{Path: "/config/path-mappings", Table: "path_mappings"},
	{Path: "/config/key-mappings", Table: "key_mappings"},
	{Path: "/config/behaviors", Table: "behaviors"},
}
var requiredKeyMappings = []KeyMapping{
	{Key: "id", Column: "path_mapping_id"},
	{Key: "id", Column: "key_mapping_id"},
	{Key: "id", Column: "behavior_id"},
	{Key: "key", Column: "key"},
	{Key: "column", Column: "column"},
	{Key: "path", Column: "path"},
	{Key: "table", Column: "table"},
	{Key: "path_mapping_id", Column: "path_mapping_id"},
	{Key: "key_mapping_id", Column: "key_mapping_id"},
}
var requiredBehaviors = []Behavior{
	{PathMapping: PathMapping{Path: "/config/path-mappings", Table: "path_mappings"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "path_mapping_id"},
		{Key: "path", Column: "path"},
		{Key: "table", Column: "table"},
	},
	},
	{PathMapping: PathMapping{Path: "/config/key-mappings", Table: "key_mappings"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "key_mapping_id"},
		{Key: "key", Column: "key"},
		{Key: "column", Column: "column"},
	},
	},
	{PathMapping: PathMapping{Path: "/config/behaviors", Table: "behaviors"}, KeyMappings: []KeyMapping{
		{Key: "id", Column: "behavior_id"},
		{Key: "path_mapping_id", Column: "path_mapping_id"},
		{Key: "key_mapping_id", Column: "key_mapping_id"},
	},
	},
}

func init() {

	Conn = connection.GetConnection()

}
