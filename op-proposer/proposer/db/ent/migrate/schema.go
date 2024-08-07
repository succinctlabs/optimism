// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ProofRequestsColumns holds the columns for the "proof_requests" table.
	ProofRequestsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"SPAN", "AGG"}},
		{Name: "start_block", Type: field.TypeUint64},
		{Name: "end_block", Type: field.TypeUint64},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"UNREQ", "REQ", "FAILED", "COMPLETE"}},
		{Name: "prover_request_id", Type: field.TypeString, Nullable: true},
		{Name: "proof_request_time", Type: field.TypeInt64, Nullable: true},
		{Name: "l1_block_number", Type: field.TypeUint64, Nullable: true},
		{Name: "l1_block_hash", Type: field.TypeString, Nullable: true},
	}
	// ProofRequestsTable holds the schema information for the "proof_requests" table.
	ProofRequestsTable = &schema.Table{
		Name:       "proof_requests",
		Columns:    ProofRequestsColumns,
		PrimaryKey: []*schema.Column{ProofRequestsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ProofRequestsTable,
	}
)

func init() {
}
