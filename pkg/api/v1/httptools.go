package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// parseUUID parses and validates a UUID from the request path if the path param is found
func (r *Router) parseUUID(c echo.Context, path string) (string, error) {
	var id string
	if err := echo.PathParamsBinder(c).String(path, &id).BindError(); err != nil {
		return "", err
	}

	if id != "" {
		if _, err := uuid.Parse(id); err != nil {
			return "", ErrInvalidUUID
		}

		return id, nil
	}

	return "", ErrUUIDNotFound
}

// queryParamsToQueryMods is a helper function that takes a echo.ValueBinder, table name,
// and column name to a append a slice of query mods.
func queryParamsToQueryMods(qpb *echo.ValueBinder, column string, mods []qm.QueryMod) []qm.QueryMod {
	var value string

	_ = qpb.String(column, &value)

	if value != "" {
		mods = append(mods, qm.Where(column+" = ?", value))
	}

	return mods
}
