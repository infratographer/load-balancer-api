// Package locations provides a the CRUD operations for locations
package locations

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"go.uber.org/zap"

	"go.infratographer.sh/loadbalancerapi/internal/dbtools"
	"go.infratographer.sh/loadbalancerapi/internal/models"
)

func init() {
	logger = zap.NewNop().Sugar()
}

func fakeBody(payload string) io.ReadCloser {
	return io.NopCloser(bytes.NewBufferString(payload))
}

var (
	errInvalidUUID = errors.New("invalid UUID length: 7")
)

func TestNewLocation(t *testing.T) {
	tenantID := uuid.New()
	happyPathBody := `{"display_name": "Nemo", "tenant_id": "` + tenantID.String() + `"}`
	missingTenantIDBody := `{"display_name": "Nemo"}`
	missingDisplayNameBody := `{"tenant_id": "` + tenantID.String() + `"}`
	invalidUUIDBody := `{"display_name": "Nemo", "tenant_id": "1234567"}`

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name    string
		args    args
		want    *Location
		wantErr bool
		err     error
	}{
		{
			name: "Happy path",
			args: args{
				c: &[]gin.Context{
					{
						Request: &http.Request{
							Body: fakeBody(happyPathBody),
						},
					},
				}[0],
			},
			want: &Location{
				TenantID: tenantID,
				Name:     "Nemo",
			},
			wantErr: false,
		},
		{
			name: "Missing tenant ID",
			args: args{
				c: &[]gin.Context{
					{
						Request: &http.Request{
							Body: fakeBody(missingTenantIDBody),
						},
					},
				}[0],
			},
			wantErr: true,
			err:     ErrTenantIDRequired,
		},
		{
			name: "Missing display name",
			args: args{
				c: &[]gin.Context{
					{
						Request: &http.Request{
							Body: fakeBody(missingDisplayNameBody),
						},
					},
				}[0],
			},
			wantErr: true,
			err:     ErrNameRequired,
		},
		{
			name: "Invalid tenant ID",
			args: args{
				c: &[]gin.Context{
					{
						Request: &http.Request{
							Body: fakeBody(invalidUUIDBody),
						},
					},
				}[0],
			},
			wantErr: true,
			err:     errInvalidUUID,
		},
	}

	for _, tt := range tests {
		got, err := NewLocation(tt.args.c)
		if tt.wantErr {
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.want.TenantID, got.TenantID)
			assert.Equal(t, tt.want.Name, got.Name)
		}
	}
}

func TestLocation_ToDBModel(t *testing.T) {
	locationID := uuid.New()
	now := time.Now()
	tenantID := uuid.New()

	type fields struct {
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *null.Time
		ID        uuid.UUID
		TenantID  uuid.UUID
		Name      string
	}

	tests := []struct {
		name    string
		fields  fields
		want    *models.Location
		wantErr bool
		err     error
	}{
		{
			name: "Happy path",
			fields: fields{
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: &null.Time{},
				TenantID:  tenantID,
				Name:      "Nemo",
			},
			want: &models.Location{
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   null.Time{},
				LocationID:  locationID.String(),
				TenantID:    tenantID.String(),
				DisplayName: "Nemo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		l := &Location{
			CreatedAt: tt.fields.CreatedAt,
			UpdatedAt: tt.fields.UpdatedAt,
			DeletedAt: tt.fields.DeletedAt,
			ID:        tt.fields.ID,
			TenantID:  tt.fields.TenantID,
			Name:      tt.fields.Name,
		}

		got, err := l.ToDBModel()
		if tt.wantErr {
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.want.TenantID, got.TenantID)
			assert.Equal(t, tt.want.DisplayName, got.DisplayName)
		}
	}
}

func TestLocation_DB(t *testing.T) {
	ctx := context.Background()
	db := dbtools.DatabaseTest(t)

	tenantID := uuid.New()

	SetLogger(zap.NewNop().Sugar())

	loc1 := &Location{
		Name:     "Nemo",
		TenantID: tenantID,
	}

	err := loc1.Create(ctx, db)
	assert.Nil(t, err)

	uuid1, err := uuid.Parse(loc1.ID.String())
	assert.Nil(t, err)

	assert.Len(t, uuid1.String(), 36)

	locArray, err := GetLocations(ctx, db, tenantID)
	assert.Nil(t, err)
	assert.Len(t, locArray, 1)

	err = loc1.Delete(ctx, db)
	assert.Nil(t, err)

	dbtools.CleanUpTables(t, tenantID, "locations")
}
