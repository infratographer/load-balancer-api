// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Pool is an object representing the database table.
type Pool struct {
	CreatedAt   time.Time `query:"created_at" param:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `query:"updated_at" param:"updated_at" boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	DeletedAt   null.Time `query:"deleted_at" param:"deleted_at" boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	PoolID      string    `query:"pool_id" param:"pool_id" boil:"pool_id" json:"pool_id" toml:"pool_id" yaml:"pool_id"`
	Protocol    string    `query:"protocol" param:"protocol" boil:"protocol" json:"protocol" toml:"protocol" yaml:"protocol"`
	DisplayName string    `query:"display_name" param:"display_name" boil:"display_name" json:"display_name" toml:"display_name" yaml:"display_name"`
	Slug        string    `query:"slug" param:"slug" boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	TenantID    string    `query:"tenant_id" param:"tenant_id" boil:"tenant_id" json:"tenant_id" toml:"tenant_id" yaml:"tenant_id"`

	R *poolR `query:"-" param:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L poolL  `query:"-" param:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PoolColumns = struct {
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	PoolID      string
	Protocol    string
	DisplayName string
	Slug        string
	TenantID    string
}{
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	PoolID:      "pool_id",
	Protocol:    "protocol",
	DisplayName: "display_name",
	Slug:        "slug",
	TenantID:    "tenant_id",
}

var PoolTableColumns = struct {
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
	PoolID      string
	Protocol    string
	DisplayName string
	Slug        string
	TenantID    string
}{
	CreatedAt:   "pools.created_at",
	UpdatedAt:   "pools.updated_at",
	DeletedAt:   "pools.deleted_at",
	PoolID:      "pools.pool_id",
	Protocol:    "pools.protocol",
	DisplayName: "pools.display_name",
	Slug:        "pools.slug",
	TenantID:    "pools.tenant_id",
}

// Generated where

var PoolWhere = struct {
	CreatedAt   whereHelpertime_Time
	UpdatedAt   whereHelpertime_Time
	DeletedAt   whereHelpernull_Time
	PoolID      whereHelperstring
	Protocol    whereHelperstring
	DisplayName whereHelperstring
	Slug        whereHelperstring
	TenantID    whereHelperstring
}{
	CreatedAt:   whereHelpertime_Time{field: "\"pools\".\"created_at\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"pools\".\"updated_at\""},
	DeletedAt:   whereHelpernull_Time{field: "\"pools\".\"deleted_at\""},
	PoolID:      whereHelperstring{field: "\"pools\".\"pool_id\""},
	Protocol:    whereHelperstring{field: "\"pools\".\"protocol\""},
	DisplayName: whereHelperstring{field: "\"pools\".\"display_name\""},
	Slug:        whereHelperstring{field: "\"pools\".\"slug\""},
	TenantID:    whereHelperstring{field: "\"pools\".\"tenant_id\""},
}

// PoolRels is where relationship names are stored.
var PoolRels = struct {
	Assignments string
	Origins     string
}{
	Assignments: "Assignments",
	Origins:     "Origins",
}

// poolR is where relationships are stored.
type poolR struct {
	Assignments AssignmentSlice `query:"Assignments" param:"Assignments" boil:"Assignments" json:"Assignments" toml:"Assignments" yaml:"Assignments"`
	Origins     OriginSlice     `query:"Origins" param:"Origins" boil:"Origins" json:"Origins" toml:"Origins" yaml:"Origins"`
}

// NewStruct creates a new relationship struct
func (*poolR) NewStruct() *poolR {
	return &poolR{}
}

func (r *poolR) GetAssignments() AssignmentSlice {
	if r == nil {
		return nil
	}
	return r.Assignments
}

func (r *poolR) GetOrigins() OriginSlice {
	if r == nil {
		return nil
	}
	return r.Origins
}

// poolL is where Load methods for each relationship are stored.
type poolL struct{}

var (
	poolAllColumns            = []string{"created_at", "updated_at", "deleted_at", "pool_id", "protocol", "display_name", "slug", "tenant_id"}
	poolColumnsWithoutDefault = []string{"protocol", "display_name", "slug", "tenant_id"}
	poolColumnsWithDefault    = []string{"created_at", "updated_at", "deleted_at", "pool_id"}
	poolPrimaryKeyColumns     = []string{"pool_id"}
	poolGeneratedColumns      = []string{}
)

type (
	// PoolSlice is an alias for a slice of pointers to Pool.
	// This should almost always be used instead of []Pool.
	PoolSlice []*Pool
	// PoolHook is the signature for custom Pool hook methods
	PoolHook func(context.Context, boil.ContextExecutor, *Pool) error

	poolQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	poolType                 = reflect.TypeOf(&Pool{})
	poolMapping              = queries.MakeStructMapping(poolType)
	poolPrimaryKeyMapping, _ = queries.BindMapping(poolType, poolMapping, poolPrimaryKeyColumns)
	poolInsertCacheMut       sync.RWMutex
	poolInsertCache          = make(map[string]insertCache)
	poolUpdateCacheMut       sync.RWMutex
	poolUpdateCache          = make(map[string]updateCache)
	poolUpsertCacheMut       sync.RWMutex
	poolUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var poolAfterSelectHooks []PoolHook

var poolBeforeInsertHooks []PoolHook
var poolAfterInsertHooks []PoolHook

var poolBeforeUpdateHooks []PoolHook
var poolAfterUpdateHooks []PoolHook

var poolBeforeDeleteHooks []PoolHook
var poolAfterDeleteHooks []PoolHook

var poolBeforeUpsertHooks []PoolHook
var poolAfterUpsertHooks []PoolHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Pool) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Pool) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Pool) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Pool) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Pool) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Pool) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Pool) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Pool) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Pool) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range poolAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPoolHook registers your hook function for all future operations.
func AddPoolHook(hookPoint boil.HookPoint, poolHook PoolHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		poolAfterSelectHooks = append(poolAfterSelectHooks, poolHook)
	case boil.BeforeInsertHook:
		poolBeforeInsertHooks = append(poolBeforeInsertHooks, poolHook)
	case boil.AfterInsertHook:
		poolAfterInsertHooks = append(poolAfterInsertHooks, poolHook)
	case boil.BeforeUpdateHook:
		poolBeforeUpdateHooks = append(poolBeforeUpdateHooks, poolHook)
	case boil.AfterUpdateHook:
		poolAfterUpdateHooks = append(poolAfterUpdateHooks, poolHook)
	case boil.BeforeDeleteHook:
		poolBeforeDeleteHooks = append(poolBeforeDeleteHooks, poolHook)
	case boil.AfterDeleteHook:
		poolAfterDeleteHooks = append(poolAfterDeleteHooks, poolHook)
	case boil.BeforeUpsertHook:
		poolBeforeUpsertHooks = append(poolBeforeUpsertHooks, poolHook)
	case boil.AfterUpsertHook:
		poolAfterUpsertHooks = append(poolAfterUpsertHooks, poolHook)
	}
}

// One returns a single pool record from the query.
func (q poolQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Pool, error) {
	o := &Pool{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to execute a one query for pools")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Pool records from the query.
func (q poolQuery) All(ctx context.Context, exec boil.ContextExecutor) (PoolSlice, error) {
	var o []*Pool

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Pool slice")
	}

	if len(poolAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Pool records in the query.
func (q poolQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count pools rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q poolQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if pools exists")
	}

	return count > 0, nil
}

// Assignments retrieves all the assignment's Assignments with an executor.
func (o *Pool) Assignments(mods ...qm.QueryMod) assignmentQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"assignments\".\"pool_id\"=?", o.PoolID),
	)

	return Assignments(queryMods...)
}

// Origins retrieves all the origin's Origins with an executor.
func (o *Pool) Origins(mods ...qm.QueryMod) originQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"origins\".\"pool_id\"=?", o.PoolID),
	)

	return Origins(queryMods...)
}

// LoadAssignments allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (poolL) LoadAssignments(ctx context.Context, e boil.ContextExecutor, singular bool, maybePool interface{}, mods queries.Applicator) error {
	var slice []*Pool
	var object *Pool

	if singular {
		var ok bool
		object, ok = maybePool.(*Pool)
		if !ok {
			object = new(Pool)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePool)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePool))
			}
		}
	} else {
		s, ok := maybePool.(*[]*Pool)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePool)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePool))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &poolR{}
		}
		args = append(args, object.PoolID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &poolR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.PoolID) {
					continue Outer
				}
			}

			args = append(args, obj.PoolID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`assignments`),
		qm.WhereIn(`assignments.pool_id in ?`, args...),
		qmhelper.WhereIsNull(`assignments.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load assignments")
	}

	var resultSlice []*Assignment
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice assignments")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on assignments")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for assignments")
	}

	if len(assignmentAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Assignments = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &assignmentR{}
			}
			foreign.R.Pool = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.PoolID, foreign.PoolID) {
				local.R.Assignments = append(local.R.Assignments, foreign)
				if foreign.R == nil {
					foreign.R = &assignmentR{}
				}
				foreign.R.Pool = local
				break
			}
		}
	}

	return nil
}

// LoadOrigins allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (poolL) LoadOrigins(ctx context.Context, e boil.ContextExecutor, singular bool, maybePool interface{}, mods queries.Applicator) error {
	var slice []*Pool
	var object *Pool

	if singular {
		var ok bool
		object, ok = maybePool.(*Pool)
		if !ok {
			object = new(Pool)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePool)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePool))
			}
		}
	} else {
		s, ok := maybePool.(*[]*Pool)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePool)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePool))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &poolR{}
		}
		args = append(args, object.PoolID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &poolR{}
			}

			for _, a := range args {
				if a == obj.PoolID {
					continue Outer
				}
			}

			args = append(args, obj.PoolID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`origins`),
		qm.WhereIn(`origins.pool_id in ?`, args...),
		qmhelper.WhereIsNull(`origins.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load origins")
	}

	var resultSlice []*Origin
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice origins")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on origins")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for origins")
	}

	if len(originAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Origins = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &originR{}
			}
			foreign.R.Pool = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.PoolID == foreign.PoolID {
				local.R.Origins = append(local.R.Origins, foreign)
				if foreign.R == nil {
					foreign.R = &originR{}
				}
				foreign.R.Pool = local
				break
			}
		}
	}

	return nil
}

// AddAssignments adds the given related objects to the existing relationships
// of the pool, optionally inserting them as new records.
// Appends related to o.R.Assignments.
// Sets related.R.Pool appropriately.
func (o *Pool) AddAssignments(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Assignment) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.PoolID, o.PoolID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"assignments\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"pool_id"}),
				strmangle.WhereClause("\"", "\"", 2, assignmentPrimaryKeyColumns),
			)
			values := []interface{}{o.PoolID, rel.AssignmentID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.PoolID, o.PoolID)
		}
	}

	if o.R == nil {
		o.R = &poolR{
			Assignments: related,
		}
	} else {
		o.R.Assignments = append(o.R.Assignments, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &assignmentR{
				Pool: o,
			}
		} else {
			rel.R.Pool = o
		}
	}
	return nil
}

// SetAssignments removes all previously related items of the
// pool replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Pool's Assignments accordingly.
// Replaces o.R.Assignments with related.
// Sets related.R.Pool's Assignments accordingly.
func (o *Pool) SetAssignments(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Assignment) error {
	query := "update \"assignments\" set \"pool_id\" = null where \"pool_id\" = $1"
	values := []interface{}{o.PoolID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Assignments {
			queries.SetScanner(&rel.PoolID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Pool = nil
		}
		o.R.Assignments = nil
	}

	return o.AddAssignments(ctx, exec, insert, related...)
}

// RemoveAssignments relationships from objects passed in.
// Removes related items from R.Assignments (uses pointer comparison, removal does not keep order)
// Sets related.R.Pool.
func (o *Pool) RemoveAssignments(ctx context.Context, exec boil.ContextExecutor, related ...*Assignment) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.PoolID, nil)
		if rel.R != nil {
			rel.R.Pool = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("pool_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Assignments {
			if rel != ri {
				continue
			}

			ln := len(o.R.Assignments)
			if ln > 1 && i < ln-1 {
				o.R.Assignments[i] = o.R.Assignments[ln-1]
			}
			o.R.Assignments = o.R.Assignments[:ln-1]
			break
		}
	}

	return nil
}

// AddOrigins adds the given related objects to the existing relationships
// of the pool, optionally inserting them as new records.
// Appends related to o.R.Origins.
// Sets related.R.Pool appropriately.
func (o *Pool) AddOrigins(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Origin) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PoolID = o.PoolID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"origins\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"pool_id"}),
				strmangle.WhereClause("\"", "\"", 2, originPrimaryKeyColumns),
			)
			values := []interface{}{o.PoolID, rel.OriginID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PoolID = o.PoolID
		}
	}

	if o.R == nil {
		o.R = &poolR{
			Origins: related,
		}
	} else {
		o.R.Origins = append(o.R.Origins, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &originR{
				Pool: o,
			}
		} else {
			rel.R.Pool = o
		}
	}
	return nil
}

// Pools retrieves all the records using an executor.
func Pools(mods ...qm.QueryMod) poolQuery {
	mods = append(mods, qm.From("\"pools\""), qmhelper.WhereIsNull("\"pools\".\"deleted_at\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"pools\".*"})
	}

	return poolQuery{q}
}

// FindPool retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPool(ctx context.Context, exec boil.ContextExecutor, poolID string, selectCols ...string) (*Pool, error) {
	poolObj := &Pool{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"pools\" where \"pool_id\"=$1 and \"deleted_at\" is null", sel,
	)

	q := queries.Raw(query, poolID)

	err := q.Bind(ctx, exec, poolObj)
	if err != nil {
		return nil, errors.Wrap(err, "models: unable to select from pools")
	}

	if err = poolObj.doAfterSelectHooks(ctx, exec); err != nil {
		return poolObj, err
	}

	return poolObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Pool) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no pools provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(poolColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	poolInsertCacheMut.RLock()
	cache, cached := poolInsertCache[key]
	poolInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			poolAllColumns,
			poolColumnsWithDefault,
			poolColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(poolType, poolMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(poolType, poolMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"pools\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"pools\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into pools")
	}

	if !cached {
		poolInsertCacheMut.Lock()
		poolInsertCache[key] = cache
		poolInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Pool.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Pool) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	poolUpdateCacheMut.RLock()
	cache, cached := poolUpdateCache[key]
	poolUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			poolAllColumns,
			poolPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update pools, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"pools\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, poolPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(poolType, poolMapping, append(wl, poolPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update pools row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for pools")
	}

	if !cached {
		poolUpdateCacheMut.Lock()
		poolUpdateCache[key] = cache
		poolUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q poolQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for pools")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for pools")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PoolSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), poolPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"pools\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, poolPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in pool slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all pool")
	}
	return rowsAff, nil
}

// Delete deletes a single Pool record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Pool) Delete(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Pool provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), poolPrimaryKeyMapping)
		sql = "DELETE FROM \"pools\" WHERE \"pool_id\"=$1"
	} else {
		currTime := time.Now().In(boil.GetLocation())
		o.DeletedAt = null.TimeFrom(currTime)
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"pools\" SET %s WHERE \"pool_id\"=$2",
			strmangle.SetParamNames("\"", "\"", 1, wl),
		)
		valueMapping, err := queries.BindMapping(poolType, poolMapping, append(wl, poolPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), valueMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from pools")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for pools")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q poolQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no poolQuery provided for delete all")
	}

	if hardDelete {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{"deleted_at": currTime})
	}

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from pools")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for pools")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PoolSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(poolBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), poolPrimaryKeyMapping)
			args = append(args, pkeyArgs...)
		}
		sql = "DELETE FROM \"pools\" WHERE " +
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, poolPrimaryKeyColumns, len(o))
	} else {
		currTime := time.Now().In(boil.GetLocation())
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), poolPrimaryKeyMapping)
			args = append(args, pkeyArgs...)
			obj.DeletedAt = null.TimeFrom(currTime)
		}
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"pools\" SET %s WHERE "+
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 2, poolPrimaryKeyColumns, len(o)),
			strmangle.SetParamNames("\"", "\"", 1, wl),
		)
		args = append([]interface{}{currTime}, args...)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from pool slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for pools")
	}

	if len(poolAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Pool) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPool(ctx, exec, o.PoolID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PoolSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PoolSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), poolPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"pools\".* FROM \"pools\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, poolPrimaryKeyColumns, len(*o)) +
		"and \"deleted_at\" is null"

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PoolSlice")
	}

	*o = slice

	return nil
}

// PoolExists checks if the Pool row exists.
func PoolExists(ctx context.Context, exec boil.ContextExecutor, poolID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"pools\" where \"pool_id\"=$1 and \"deleted_at\" is null limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, poolID)
	}
	row := exec.QueryRowContext(ctx, sql, poolID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if pools exists")
	}

	return exists, nil
}

// Exists checks if the Pool row exists.
func (o *Pool) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return PoolExists(ctx, exec, o.PoolID)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Pool) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no pools provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(poolColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	poolUpsertCacheMut.RLock()
	cache, cached := poolUpsertCache[key]
	poolUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			poolAllColumns,
			poolColumnsWithDefault,
			poolColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			poolAllColumns,
			poolPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert pools, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(poolPrimaryKeyColumns))
			copy(conflict, poolPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryCockroachDB(dialect, "\"pools\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(poolType, poolMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(poolType, poolMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		_, _ = fmt.Fprintln(boil.DebugWriter, cache.query)
		_, _ = fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // CockcorachDB doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert pools")
	}

	if !cached {
		poolUpsertCacheMut.Lock()
		poolUpsertCache[key] = cache
		poolUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
