// Package locations provides a the CRUD operations for locations
package locations

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.infratographer.sh/loadbalancerapi/internal/models"
)

var logger *zap.SugaredLogger

// SetLogger sets the logger for this package
func SetLogger(l *zap.SugaredLogger) {
	logger = l
}

func init() {
	logger = zap.NewNop().Sugar()
}

func qmName(name string) qm.QueryMod {
	return qm.Where("display_name = ?", name)
}

func qmTenantID(tenantID uuid.UUID) qm.QueryMod {
	return qm.Where("tenant_id = ?", tenantID.String())
}

func qmCombine(mods ...qm.QueryMod) qm.QueryMod {
	return qm.Expr(mods...)
}

// GetLocations returns all locations from the database
func GetLocations(ctx context.Context, db *sqlx.DB, tenant uuid.UUID) (Locations, error) {
	qm := qmTenantID(tenant)

	dbms, err := models.Locations(qm).All(ctx, db)
	if err != nil {
		return nil, err
	}

	ls := make(Locations, len(dbms))

	for i, dbm := range dbms {
		l := &Location{}
		if err := l.FromDBModel(ctx, db, dbm); err != nil {
			return nil, err
		}

		ls[i] = l
	}

	return ls, nil
}

// Create inserts a new location into the database
func (l *Location) Create(ctx context.Context, db *sqlx.DB) error {
	dbm, err := l.ToDBModel()
	if err != nil {
		return err
	}

	if err := dbm.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return l.FromDBModel(ctx, db, dbm)
}

// Find returns a location from the database
func (l *Location) Find(ctx context.Context, db *sqlx.DB) error {
	if err := l.validate(); err != nil {
		return err
	}

	mods := qmCombine(qmTenantID(l.TenantID), qmName(l.Name))

	dbm, err := models.Locations(mods).One(ctx, db)
	if err != nil {
		return err
	}

	if err := l.FromDBModel(ctx, db, dbm); err != nil {
		return err
	}

	if l.ID.String() == uuid.Nil.String() {
		return ErrNullUUID
	}

	return nil
}

// Delete removes a location from the database
func (l *Location) Delete(ctx context.Context, db *sqlx.DB) error {
	err := l.Find(ctx, db)
	if err != nil {
		return err
	}

	dbm, err := l.ToDBModel()
	if err != nil {
		return err
	}

	if dbm.LocationID == "" {
		return ErrNullUUID
	}

	// soft delete
	_, err = dbm.Delete(ctx, db, false)
	if err != nil {
		return err
	}

	return nil
}

// FromDBModel converts a DB model to a Location
func (l *Location) FromDBModel(ctx context.Context, db *sqlx.DB, dbm *models.Location) error {
	l.Name = dbm.DisplayName
	l.CreatedAt = dbm.CreatedAt
	l.UpdatedAt = dbm.UpdatedAt

	if dbm.DeletedAt.IsZero() {
		l.DeletedAt = nil
	} else {
		l.DeletedAt = &dbm.DeletedAt
	}

	var err error

	l.ID, err = uuid.Parse(dbm.LocationID)
	if err != nil {
		return err
	}

	l.TenantID, err = uuid.Parse(dbm.TenantID)
	if err != nil {
		return err
	}

	if l.ID.String() == uuid.Nil.String() {
		return ErrNullUUID
	}

	logger.Debugw("location from db model", "api-model", l, "db-model", dbm)

	return nil
}

// ToDBModel converts a Location to a DB model
func (l *Location) ToDBModel() (*models.Location, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	dbm := &models.Location{
		DisplayName: l.Name,
		TenantID:    l.TenantID.String(),
	}

	if l.ID.String() != uuid.Nil.String() {
		dbm.LocationID = l.ID.String()
	}

	return dbm, nil
}

// NewLocation creates a new location from a gin context
func NewLocation(c *gin.Context) (*Location, error) {
	l := &Location{}

	if err := c.ShouldBindJSON(l); err != nil {
		return nil, err
	}

	if err := l.validate(); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *Location) validate() error {
	if l.Name == "" {
		return ErrNameRequired
	}

	if l.TenantID.String() == uuid.Nil.String() {
		return ErrTenantIDRequired
	}

	return nil
}
