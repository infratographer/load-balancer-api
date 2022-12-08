// Package loadbalancers provides a the CRUD operations for locations
package loadbalancers

import (
	"context"

	"github.com/dspinhirne/netaddr-go"
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

func qmIPAddress(ip string) qm.QueryMod {
	return qm.Where("ip_addr = ?", ip)
}

func qmTenantID(tenantID uuid.UUID) qm.QueryMod {
	return qm.Where("tenant_id = ?", tenantID.String())
}

func qmNotDeleted() qm.QueryMod {
	return qm.Where("deleted_at IS NULL")
}

func qmCombine(mods ...qm.QueryMod) qm.QueryMod {
	return qm.Expr(mods...)
}

// GetLoadBalancers returns all locations from the database
func GetLoadBalancers(ctx context.Context, db *sqlx.DB, tenant uuid.UUID) (LoadBalancers, error) {
	qm := qmTenantID(tenant)

	dbms, err := models.LoadBalancers(qm).All(ctx, db)
	if err != nil {
		return nil, err
	}

	lbs := make(LoadBalancers, len(dbms))

	for i, dbm := range dbms {
		lb := &LoadBalancer{}
		if err := lb.FromDBModel(ctx, db, dbm); err != nil {
			return nil, err
		}

		lbs[i] = lb
	}

	return lbs, nil
}

// Create inserts a new location into the database
func (lb *LoadBalancer) Create(ctx context.Context, db *sqlx.DB) error {
	dbm, err := lb.ToDBModel()
	if err != nil {
		return err
	}

	if err := dbm.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return lb.FromDBModel(ctx, db, dbm)
}

// Find returns a load balancers from the database
func (lb *LoadBalancer) Find(ctx context.Context, db *sqlx.DB) error {
	if lb.IPAddress == "" {
		return ErrIPAddressRequired
	}

	if lb.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}

	mods := qmCombine(qmTenantID(lb.TenantID), qmIPAddress(lb.IPAddress), qmNotDeleted())

	dbm, err := models.LoadBalancers(mods).One(ctx, db)
	if err != nil {
		logger.Debugw("failed to find load balancer", "error", err, "load_balancer", lb, "mods", mods)
		return err
	}

	logger.Debugw("found load balancer", "load_balancer", lb, "mods", mods)

	return lb.FromDBModel(ctx, db, dbm)
}

// Delete removes a location from the database
func (lb *LoadBalancer) Delete(ctx context.Context, db *sqlx.DB) error {
	err := lb.Find(ctx, db)
	if err != nil {
		return err
	}

	dbm, err := lb.ToDBModel()
	if err != nil {
		return err
	}

	boil.DebugMode = true

	// soft delete
	_, err = dbm.Delete(ctx, db, false)
	if err != nil {
		logger.Errorw("failed to delete load balancer", "error", err, "load_balancer", lb, "db_model", dbm)
		return err
	}

	return nil
}

// FromDBModel converts a DB model to a LoadBalancer
func (lb *LoadBalancer) FromDBModel(ctx context.Context, db *sqlx.DB, dbm *models.LoadBalancer) error {
	lb.Name = dbm.DisplayName
	lb.Size = dbm.LoadBalancerSize
	lb.Type = dbm.LoadBalancerType
	lb.CreatedAt = dbm.CreatedAt
	lb.UpdatedAt = dbm.UpdatedAt

	if dbm.DeletedAt.IsZero() {
		lb.DeletedAt = nil
	} else {
		lb.DeletedAt = &dbm.DeletedAt
	}

	var err error

	lb.ID, err = uuid.Parse(dbm.LoadBalancerID)
	if err != nil {
		return err
	}

	lb.LocationID, err = uuid.Parse(dbm.LocationID)
	if err != nil {
		return err
	}

	lb.TenantID, err = uuid.Parse(dbm.TenantID)
	if err != nil {
		return err
	}

	ip, err := netaddr.ParseIP(dbm.IPAddr)
	if err != nil {
		return err
	}

	lb.IPAddress = ip.String()

	if err := lb.validate(); err != nil {
		return err
	}

	logger.Debugw("location from db model", "api-model", lb, "db-model", dbm)

	return nil
}

// ToDBModel converts a Location to a DB model
func (lb *LoadBalancer) ToDBModel() (*models.LoadBalancer, error) {
	if err := lb.validate(); err != nil {
		return nil, err
	}

	dbm := &models.LoadBalancer{
		DisplayName:      lb.Name,
		IPAddr:           lb.IPAddress,
		LoadBalancerSize: lb.Size,
		LoadBalancerType: lb.Type,
		LocationID:       lb.LocationID.String(),
		TenantID:         lb.TenantID.String(),
	}

	if lb.ID != uuid.Nil {
		dbm.LoadBalancerID = lb.ID.String()
	}

	if err := lb.validate(); err != nil {
		return nil, err
	}

	return dbm, nil
}

// NewLoadBalancer creates a new location from a gin context
func NewLoadBalancer(c *gin.Context) (*LoadBalancer, error) {
	l := &LoadBalancer{}

	if err := c.BindJSON(l); err != nil {
		return nil, err
	}

	if err := l.validate(); err != nil {
		return nil, err
	}

	return l, nil
}

func (lb *LoadBalancer) validate() error {
	if lb.Name == "" {
		return ErrNameRequired
	}

	if lb.Size == "" {
		return ErrSizeRequired
	}

	if lb.Type == "" {
		return ErrTypeRequired
	}

	// Only allow layer-3 load balancers for now
	if lb.Type != "layer-3" {
		return ErrTypeInvalid
	}

	if lb.LocationID.String() == uuid.Nil.String() {
		return ErrLocationIDRequired
	}

	if ip, err := netaddr.ParseIP(lb.IPAddress); err != nil {
		if ip.Version() != 4 { //nolint:gomnd
			return ErrIPv4Required
		}

		return ErrIPAddressInvalid
	}

	if lb.TenantID.String() == uuid.Nil.String() {
		return ErrTenantIDRequired
	}

	if lb.DeletedAt != nil {
		if lb.DeletedAt.IsZero() {
			return ErrAlreadyDeleted
		}
	}

	return nil
}
