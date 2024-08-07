// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ethereum-optimism/optimism/op-proposer/proposer/db/ent/predicate"
	"github.com/ethereum-optimism/optimism/op-proposer/proposer/db/ent/proofrequest"
)

// ProofRequestUpdate is the builder for updating ProofRequest entities.
type ProofRequestUpdate struct {
	config
	hooks    []Hook
	mutation *ProofRequestMutation
}

// Where appends a list predicates to the ProofRequestUpdate builder.
func (pru *ProofRequestUpdate) Where(ps ...predicate.ProofRequest) *ProofRequestUpdate {
	pru.mutation.Where(ps...)
	return pru
}

// SetType sets the "type" field.
func (pru *ProofRequestUpdate) SetType(pr proofrequest.Type) *ProofRequestUpdate {
	pru.mutation.SetType(pr)
	return pru
}

// SetNillableType sets the "type" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableType(pr *proofrequest.Type) *ProofRequestUpdate {
	if pr != nil {
		pru.SetType(*pr)
	}
	return pru
}

// SetStartBlock sets the "start_block" field.
func (pru *ProofRequestUpdate) SetStartBlock(u uint64) *ProofRequestUpdate {
	pru.mutation.ResetStartBlock()
	pru.mutation.SetStartBlock(u)
	return pru
}

// SetNillableStartBlock sets the "start_block" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableStartBlock(u *uint64) *ProofRequestUpdate {
	if u != nil {
		pru.SetStartBlock(*u)
	}
	return pru
}

// AddStartBlock adds u to the "start_block" field.
func (pru *ProofRequestUpdate) AddStartBlock(u int64) *ProofRequestUpdate {
	pru.mutation.AddStartBlock(u)
	return pru
}

// SetEndBlock sets the "end_block" field.
func (pru *ProofRequestUpdate) SetEndBlock(u uint64) *ProofRequestUpdate {
	pru.mutation.ResetEndBlock()
	pru.mutation.SetEndBlock(u)
	return pru
}

// SetNillableEndBlock sets the "end_block" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableEndBlock(u *uint64) *ProofRequestUpdate {
	if u != nil {
		pru.SetEndBlock(*u)
	}
	return pru
}

// AddEndBlock adds u to the "end_block" field.
func (pru *ProofRequestUpdate) AddEndBlock(u int64) *ProofRequestUpdate {
	pru.mutation.AddEndBlock(u)
	return pru
}

// SetStatus sets the "status" field.
func (pru *ProofRequestUpdate) SetStatus(pr proofrequest.Status) *ProofRequestUpdate {
	pru.mutation.SetStatus(pr)
	return pru
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableStatus(pr *proofrequest.Status) *ProofRequestUpdate {
	if pr != nil {
		pru.SetStatus(*pr)
	}
	return pru
}

// SetProverRequestID sets the "prover_request_id" field.
func (pru *ProofRequestUpdate) SetProverRequestID(s string) *ProofRequestUpdate {
	pru.mutation.SetProverRequestID(s)
	return pru
}

// SetNillableProverRequestID sets the "prover_request_id" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableProverRequestID(s *string) *ProofRequestUpdate {
	if s != nil {
		pru.SetProverRequestID(*s)
	}
	return pru
}

// ClearProverRequestID clears the value of the "prover_request_id" field.
func (pru *ProofRequestUpdate) ClearProverRequestID() *ProofRequestUpdate {
	pru.mutation.ClearProverRequestID()
	return pru
}

// SetProofRequestTime sets the "proof_request_time" field.
func (pru *ProofRequestUpdate) SetProofRequestTime(i int64) *ProofRequestUpdate {
	pru.mutation.ResetProofRequestTime()
	pru.mutation.SetProofRequestTime(i)
	return pru
}

// SetNillableProofRequestTime sets the "proof_request_time" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableProofRequestTime(i *int64) *ProofRequestUpdate {
	if i != nil {
		pru.SetProofRequestTime(*i)
	}
	return pru
}

// AddProofRequestTime adds i to the "proof_request_time" field.
func (pru *ProofRequestUpdate) AddProofRequestTime(i int64) *ProofRequestUpdate {
	pru.mutation.AddProofRequestTime(i)
	return pru
}

// ClearProofRequestTime clears the value of the "proof_request_time" field.
func (pru *ProofRequestUpdate) ClearProofRequestTime() *ProofRequestUpdate {
	pru.mutation.ClearProofRequestTime()
	return pru
}

// SetL1BlockNumber sets the "l1_block_number" field.
func (pru *ProofRequestUpdate) SetL1BlockNumber(u uint64) *ProofRequestUpdate {
	pru.mutation.ResetL1BlockNumber()
	pru.mutation.SetL1BlockNumber(u)
	return pru
}

// SetNillableL1BlockNumber sets the "l1_block_number" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableL1BlockNumber(u *uint64) *ProofRequestUpdate {
	if u != nil {
		pru.SetL1BlockNumber(*u)
	}
	return pru
}

// AddL1BlockNumber adds u to the "l1_block_number" field.
func (pru *ProofRequestUpdate) AddL1BlockNumber(u int64) *ProofRequestUpdate {
	pru.mutation.AddL1BlockNumber(u)
	return pru
}

// ClearL1BlockNumber clears the value of the "l1_block_number" field.
func (pru *ProofRequestUpdate) ClearL1BlockNumber() *ProofRequestUpdate {
	pru.mutation.ClearL1BlockNumber()
	return pru
}

// SetL1BlockHash sets the "l1_block_hash" field.
func (pru *ProofRequestUpdate) SetL1BlockHash(s string) *ProofRequestUpdate {
	pru.mutation.SetL1BlockHash(s)
	return pru
}

// SetNillableL1BlockHash sets the "l1_block_hash" field if the given value is not nil.
func (pru *ProofRequestUpdate) SetNillableL1BlockHash(s *string) *ProofRequestUpdate {
	if s != nil {
		pru.SetL1BlockHash(*s)
	}
	return pru
}

// ClearL1BlockHash clears the value of the "l1_block_hash" field.
func (pru *ProofRequestUpdate) ClearL1BlockHash() *ProofRequestUpdate {
	pru.mutation.ClearL1BlockHash()
	return pru
}

// Mutation returns the ProofRequestMutation object of the builder.
func (pru *ProofRequestUpdate) Mutation() *ProofRequestMutation {
	return pru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pru *ProofRequestUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, pru.sqlSave, pru.mutation, pru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pru *ProofRequestUpdate) SaveX(ctx context.Context) int {
	affected, err := pru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pru *ProofRequestUpdate) Exec(ctx context.Context) error {
	_, err := pru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pru *ProofRequestUpdate) ExecX(ctx context.Context) {
	if err := pru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pru *ProofRequestUpdate) check() error {
	if v, ok := pru.mutation.GetType(); ok {
		if err := proofrequest.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "ProofRequest.type": %w`, err)}
		}
	}
	if v, ok := pru.mutation.Status(); ok {
		if err := proofrequest.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "ProofRequest.status": %w`, err)}
		}
	}
	return nil
}

func (pru *ProofRequestUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(proofrequest.Table, proofrequest.Columns, sqlgraph.NewFieldSpec(proofrequest.FieldID, field.TypeInt))
	if ps := pru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pru.mutation.GetType(); ok {
		_spec.SetField(proofrequest.FieldType, field.TypeEnum, value)
	}
	if value, ok := pru.mutation.StartBlock(); ok {
		_spec.SetField(proofrequest.FieldStartBlock, field.TypeUint64, value)
	}
	if value, ok := pru.mutation.AddedStartBlock(); ok {
		_spec.AddField(proofrequest.FieldStartBlock, field.TypeUint64, value)
	}
	if value, ok := pru.mutation.EndBlock(); ok {
		_spec.SetField(proofrequest.FieldEndBlock, field.TypeUint64, value)
	}
	if value, ok := pru.mutation.AddedEndBlock(); ok {
		_spec.AddField(proofrequest.FieldEndBlock, field.TypeUint64, value)
	}
	if value, ok := pru.mutation.Status(); ok {
		_spec.SetField(proofrequest.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := pru.mutation.ProverRequestID(); ok {
		_spec.SetField(proofrequest.FieldProverRequestID, field.TypeString, value)
	}
	if pru.mutation.ProverRequestIDCleared() {
		_spec.ClearField(proofrequest.FieldProverRequestID, field.TypeString)
	}
	if value, ok := pru.mutation.ProofRequestTime(); ok {
		_spec.SetField(proofrequest.FieldProofRequestTime, field.TypeInt64, value)
	}
	if value, ok := pru.mutation.AddedProofRequestTime(); ok {
		_spec.AddField(proofrequest.FieldProofRequestTime, field.TypeInt64, value)
	}
	if pru.mutation.ProofRequestTimeCleared() {
		_spec.ClearField(proofrequest.FieldProofRequestTime, field.TypeInt64)
	}
	if value, ok := pru.mutation.L1BlockNumber(); ok {
		_spec.SetField(proofrequest.FieldL1BlockNumber, field.TypeUint64, value)
	}
	if value, ok := pru.mutation.AddedL1BlockNumber(); ok {
		_spec.AddField(proofrequest.FieldL1BlockNumber, field.TypeUint64, value)
	}
	if pru.mutation.L1BlockNumberCleared() {
		_spec.ClearField(proofrequest.FieldL1BlockNumber, field.TypeUint64)
	}
	if value, ok := pru.mutation.L1BlockHash(); ok {
		_spec.SetField(proofrequest.FieldL1BlockHash, field.TypeString, value)
	}
	if pru.mutation.L1BlockHashCleared() {
		_spec.ClearField(proofrequest.FieldL1BlockHash, field.TypeString)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{proofrequest.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pru.mutation.done = true
	return n, nil
}

// ProofRequestUpdateOne is the builder for updating a single ProofRequest entity.
type ProofRequestUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ProofRequestMutation
}

// SetType sets the "type" field.
func (pruo *ProofRequestUpdateOne) SetType(pr proofrequest.Type) *ProofRequestUpdateOne {
	pruo.mutation.SetType(pr)
	return pruo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableType(pr *proofrequest.Type) *ProofRequestUpdateOne {
	if pr != nil {
		pruo.SetType(*pr)
	}
	return pruo
}

// SetStartBlock sets the "start_block" field.
func (pruo *ProofRequestUpdateOne) SetStartBlock(u uint64) *ProofRequestUpdateOne {
	pruo.mutation.ResetStartBlock()
	pruo.mutation.SetStartBlock(u)
	return pruo
}

// SetNillableStartBlock sets the "start_block" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableStartBlock(u *uint64) *ProofRequestUpdateOne {
	if u != nil {
		pruo.SetStartBlock(*u)
	}
	return pruo
}

// AddStartBlock adds u to the "start_block" field.
func (pruo *ProofRequestUpdateOne) AddStartBlock(u int64) *ProofRequestUpdateOne {
	pruo.mutation.AddStartBlock(u)
	return pruo
}

// SetEndBlock sets the "end_block" field.
func (pruo *ProofRequestUpdateOne) SetEndBlock(u uint64) *ProofRequestUpdateOne {
	pruo.mutation.ResetEndBlock()
	pruo.mutation.SetEndBlock(u)
	return pruo
}

// SetNillableEndBlock sets the "end_block" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableEndBlock(u *uint64) *ProofRequestUpdateOne {
	if u != nil {
		pruo.SetEndBlock(*u)
	}
	return pruo
}

// AddEndBlock adds u to the "end_block" field.
func (pruo *ProofRequestUpdateOne) AddEndBlock(u int64) *ProofRequestUpdateOne {
	pruo.mutation.AddEndBlock(u)
	return pruo
}

// SetStatus sets the "status" field.
func (pruo *ProofRequestUpdateOne) SetStatus(pr proofrequest.Status) *ProofRequestUpdateOne {
	pruo.mutation.SetStatus(pr)
	return pruo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableStatus(pr *proofrequest.Status) *ProofRequestUpdateOne {
	if pr != nil {
		pruo.SetStatus(*pr)
	}
	return pruo
}

// SetProverRequestID sets the "prover_request_id" field.
func (pruo *ProofRequestUpdateOne) SetProverRequestID(s string) *ProofRequestUpdateOne {
	pruo.mutation.SetProverRequestID(s)
	return pruo
}

// SetNillableProverRequestID sets the "prover_request_id" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableProverRequestID(s *string) *ProofRequestUpdateOne {
	if s != nil {
		pruo.SetProverRequestID(*s)
	}
	return pruo
}

// ClearProverRequestID clears the value of the "prover_request_id" field.
func (pruo *ProofRequestUpdateOne) ClearProverRequestID() *ProofRequestUpdateOne {
	pruo.mutation.ClearProverRequestID()
	return pruo
}

// SetProofRequestTime sets the "proof_request_time" field.
func (pruo *ProofRequestUpdateOne) SetProofRequestTime(i int64) *ProofRequestUpdateOne {
	pruo.mutation.ResetProofRequestTime()
	pruo.mutation.SetProofRequestTime(i)
	return pruo
}

// SetNillableProofRequestTime sets the "proof_request_time" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableProofRequestTime(i *int64) *ProofRequestUpdateOne {
	if i != nil {
		pruo.SetProofRequestTime(*i)
	}
	return pruo
}

// AddProofRequestTime adds i to the "proof_request_time" field.
func (pruo *ProofRequestUpdateOne) AddProofRequestTime(i int64) *ProofRequestUpdateOne {
	pruo.mutation.AddProofRequestTime(i)
	return pruo
}

// ClearProofRequestTime clears the value of the "proof_request_time" field.
func (pruo *ProofRequestUpdateOne) ClearProofRequestTime() *ProofRequestUpdateOne {
	pruo.mutation.ClearProofRequestTime()
	return pruo
}

// SetL1BlockNumber sets the "l1_block_number" field.
func (pruo *ProofRequestUpdateOne) SetL1BlockNumber(u uint64) *ProofRequestUpdateOne {
	pruo.mutation.ResetL1BlockNumber()
	pruo.mutation.SetL1BlockNumber(u)
	return pruo
}

// SetNillableL1BlockNumber sets the "l1_block_number" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableL1BlockNumber(u *uint64) *ProofRequestUpdateOne {
	if u != nil {
		pruo.SetL1BlockNumber(*u)
	}
	return pruo
}

// AddL1BlockNumber adds u to the "l1_block_number" field.
func (pruo *ProofRequestUpdateOne) AddL1BlockNumber(u int64) *ProofRequestUpdateOne {
	pruo.mutation.AddL1BlockNumber(u)
	return pruo
}

// ClearL1BlockNumber clears the value of the "l1_block_number" field.
func (pruo *ProofRequestUpdateOne) ClearL1BlockNumber() *ProofRequestUpdateOne {
	pruo.mutation.ClearL1BlockNumber()
	return pruo
}

// SetL1BlockHash sets the "l1_block_hash" field.
func (pruo *ProofRequestUpdateOne) SetL1BlockHash(s string) *ProofRequestUpdateOne {
	pruo.mutation.SetL1BlockHash(s)
	return pruo
}

// SetNillableL1BlockHash sets the "l1_block_hash" field if the given value is not nil.
func (pruo *ProofRequestUpdateOne) SetNillableL1BlockHash(s *string) *ProofRequestUpdateOne {
	if s != nil {
		pruo.SetL1BlockHash(*s)
	}
	return pruo
}

// ClearL1BlockHash clears the value of the "l1_block_hash" field.
func (pruo *ProofRequestUpdateOne) ClearL1BlockHash() *ProofRequestUpdateOne {
	pruo.mutation.ClearL1BlockHash()
	return pruo
}

// Mutation returns the ProofRequestMutation object of the builder.
func (pruo *ProofRequestUpdateOne) Mutation() *ProofRequestMutation {
	return pruo.mutation
}

// Where appends a list predicates to the ProofRequestUpdate builder.
func (pruo *ProofRequestUpdateOne) Where(ps ...predicate.ProofRequest) *ProofRequestUpdateOne {
	pruo.mutation.Where(ps...)
	return pruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pruo *ProofRequestUpdateOne) Select(field string, fields ...string) *ProofRequestUpdateOne {
	pruo.fields = append([]string{field}, fields...)
	return pruo
}

// Save executes the query and returns the updated ProofRequest entity.
func (pruo *ProofRequestUpdateOne) Save(ctx context.Context) (*ProofRequest, error) {
	return withHooks(ctx, pruo.sqlSave, pruo.mutation, pruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pruo *ProofRequestUpdateOne) SaveX(ctx context.Context) *ProofRequest {
	node, err := pruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pruo *ProofRequestUpdateOne) Exec(ctx context.Context) error {
	_, err := pruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pruo *ProofRequestUpdateOne) ExecX(ctx context.Context) {
	if err := pruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pruo *ProofRequestUpdateOne) check() error {
	if v, ok := pruo.mutation.GetType(); ok {
		if err := proofrequest.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "ProofRequest.type": %w`, err)}
		}
	}
	if v, ok := pruo.mutation.Status(); ok {
		if err := proofrequest.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "ProofRequest.status": %w`, err)}
		}
	}
	return nil
}

func (pruo *ProofRequestUpdateOne) sqlSave(ctx context.Context) (_node *ProofRequest, err error) {
	if err := pruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(proofrequest.Table, proofrequest.Columns, sqlgraph.NewFieldSpec(proofrequest.FieldID, field.TypeInt))
	id, ok := pruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ProofRequest.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, proofrequest.FieldID)
		for _, f := range fields {
			if !proofrequest.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != proofrequest.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pruo.mutation.GetType(); ok {
		_spec.SetField(proofrequest.FieldType, field.TypeEnum, value)
	}
	if value, ok := pruo.mutation.StartBlock(); ok {
		_spec.SetField(proofrequest.FieldStartBlock, field.TypeUint64, value)
	}
	if value, ok := pruo.mutation.AddedStartBlock(); ok {
		_spec.AddField(proofrequest.FieldStartBlock, field.TypeUint64, value)
	}
	if value, ok := pruo.mutation.EndBlock(); ok {
		_spec.SetField(proofrequest.FieldEndBlock, field.TypeUint64, value)
	}
	if value, ok := pruo.mutation.AddedEndBlock(); ok {
		_spec.AddField(proofrequest.FieldEndBlock, field.TypeUint64, value)
	}
	if value, ok := pruo.mutation.Status(); ok {
		_spec.SetField(proofrequest.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := pruo.mutation.ProverRequestID(); ok {
		_spec.SetField(proofrequest.FieldProverRequestID, field.TypeString, value)
	}
	if pruo.mutation.ProverRequestIDCleared() {
		_spec.ClearField(proofrequest.FieldProverRequestID, field.TypeString)
	}
	if value, ok := pruo.mutation.ProofRequestTime(); ok {
		_spec.SetField(proofrequest.FieldProofRequestTime, field.TypeInt64, value)
	}
	if value, ok := pruo.mutation.AddedProofRequestTime(); ok {
		_spec.AddField(proofrequest.FieldProofRequestTime, field.TypeInt64, value)
	}
	if pruo.mutation.ProofRequestTimeCleared() {
		_spec.ClearField(proofrequest.FieldProofRequestTime, field.TypeInt64)
	}
	if value, ok := pruo.mutation.L1BlockNumber(); ok {
		_spec.SetField(proofrequest.FieldL1BlockNumber, field.TypeUint64, value)
	}
	if value, ok := pruo.mutation.AddedL1BlockNumber(); ok {
		_spec.AddField(proofrequest.FieldL1BlockNumber, field.TypeUint64, value)
	}
	if pruo.mutation.L1BlockNumberCleared() {
		_spec.ClearField(proofrequest.FieldL1BlockNumber, field.TypeUint64)
	}
	if value, ok := pruo.mutation.L1BlockHash(); ok {
		_spec.SetField(proofrequest.FieldL1BlockHash, field.TypeString, value)
	}
	if pruo.mutation.L1BlockHashCleared() {
		_spec.ClearField(proofrequest.FieldL1BlockHash, field.TypeString)
	}
	_node = &ProofRequest{config: pruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{proofrequest.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	pruo.mutation.done = true
	return _node, nil
}
