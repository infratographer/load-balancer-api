package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// parseTenantID parses the tenant_id from the request path and returns an error if the tenant_id
// is not present or an invalid uuid is provided.
func parseTenantID(c echo.Context) (string, error) {
	var tenantID string

	err := echo.PathParamsBinder(c).MustString("tenant_id", &tenantID).BindError()
	if err != nil {
		return "", ErrTenantIDRequired
	}

	if tenantID != "" {
		if _, err := uuid.Parse(tenantID); err != nil {
			return "", ErrInvalidUUID
		}
	}

	return tenantID, nil
}

// queryParamsToQueryMods is a helper function that takes a echo.ValueBinder, table name,
// and column name to a append a slice of query mods.
func queryParamsToQueryMods(qpb *echo.ValueBinder, table, column string, mods []qm.QueryMod) []qm.QueryMod {
	var value string

	_ = qpb.String(column, &value)

	if value != "" {
		mods = append(mods, qm.Where(column+" = ?", value))
	}

	return mods
}
