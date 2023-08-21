package manualhooks

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"
	"golang.org/x/exp/slices"

	"go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/hook"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/loadbalancer"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/origin"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/pool"
	"go.infratographer.com/load-balancer-api/internal/ent/generated/port"
)

func LoadBalancerHooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.LoadBalancerFunc(func(ctx context.Context, m *generated.LoadBalancerMutation) (ent.Value, error) {
					var err error
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					additionalSubjects = append(additionalSubjects, objID)

					addSubjPort, err := m.Client().Port.Query().Where(port.HasLoadBalancerWith(loadbalancer.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPort.ID) {
							additionalSubjects = append(additionalSubjects, addSubjPort.ID)
						}
					}

					changeset := []events.FieldChange{}
					cv_created_at := ""
					created_at, ok := m.CreatedAt()

					if ok {
						cv_created_at = created_at.Format(time.RFC3339)
						pv_created_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldCreatedAt(ctx)
							if err != nil {
								pv_created_at = "<unknown>"
							} else {
								pv_created_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "created_at",
							PreviousValue: pv_created_at,
							CurrentValue:  cv_created_at,
						})
					}

					cv_updated_at := ""
					updated_at, ok := m.UpdatedAt()

					if ok {
						cv_updated_at = updated_at.Format(time.RFC3339)
						pv_updated_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldUpdatedAt(ctx)
							if err != nil {
								pv_updated_at = "<unknown>"
							} else {
								pv_updated_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "updated_at",
							PreviousValue: pv_updated_at,
							CurrentValue:  cv_updated_at,
						})
					}

					cv_name := ""
					name, ok := m.Name()

					if ok {
						cv_name = fmt.Sprintf("%s", fmt.Sprint(name))
						pv_name := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldName(ctx)
							if err != nil {
								pv_name = "<unknown>"
							} else {
								pv_name = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "name",
							PreviousValue: pv_name,
							CurrentValue:  cv_name,
						})
					}

					cv_owner_id := ""
					owner_id, ok := m.OwnerID()
					if !ok && !m.Op().Is(ent.OpCreate) {
						// since we are doing an update or delete and these fields didn't change, load the "old" value
						owner_id, err = m.OldOwnerID(ctx)
						if err != nil {
							return nil, err
						}
					}
					additionalSubjects = append(additionalSubjects, owner_id)

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "owner",
						SubjectID: owner_id,
					})

					if ok {
						cv_owner_id = fmt.Sprintf("%s", fmt.Sprint(owner_id))
						pv_owner_id := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldOwnerID(ctx)
							if err != nil {
								pv_owner_id = "<unknown>"
							} else {
								pv_owner_id = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "owner_id",
							PreviousValue: pv_owner_id,
							CurrentValue:  cv_owner_id,
						})
					}

					cv_location_id := ""
					location_id, ok := m.LocationID()
					if !ok && !m.Op().Is(ent.OpCreate) {
						// since we are doing an update or delete and these fields didn't change, load the "old" value
						location_id, err = m.OldLocationID(ctx)
						if err != nil {
							return nil, err
						}
					}
					additionalSubjects = append(additionalSubjects, location_id)

					if ok {
						cv_location_id = fmt.Sprintf("%s", fmt.Sprint(location_id))
						pv_location_id := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldLocationID(ctx)
							if err != nil {
								pv_location_id = "<unknown>"
							} else {
								pv_location_id = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "location_id",
							PreviousValue: pv_location_id,
							CurrentValue:  cv_location_id,
						})
					}

					cv_provider_id := ""
					provider_id, ok := m.ProviderID()
					if !ok && !m.Op().Is(ent.OpCreate) {
						// since we are doing an update or delete and these fields didn't change, load the "old" value
						provider_id, err = m.OldProviderID(ctx)
						if err != nil {
							return nil, err
						}
					}
					additionalSubjects = append(additionalSubjects, provider_id)

					if ok {
						cv_provider_id = fmt.Sprintf("%s", fmt.Sprint(provider_id))
						pv_provider_id := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldProviderID(ctx)
							if err != nil {
								pv_provider_id = "<unknown>"
							} else {
								pv_provider_id = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "provider_id",
							PreviousValue: pv_provider_id,
							CurrentValue:  cv_provider_id,
						})
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
						FieldChanges:         changeset,
					}

					// complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.CreateAuthRelationships(ctx, "load-balancer", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),

		// Delete Hook
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.LoadBalancerFunc(func(ctx context.Context, m *generated.LoadBalancerMutation) (ent.Value, error) {
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					dbObj, err := m.Client().LoadBalancer.Get(ctx, objID)
					if err != nil {
						return nil, fmt.Errorf("failed to load object to get values for event, err %w", err)
					}

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "owner",
						SubjectID: dbObj.OwnerID,
					})

					additionalSubjects = append(additionalSubjects, dbObj.OwnerID)
					additionalSubjects = append(additionalSubjects, dbObj.LocationID)
					additionalSubjects = append(additionalSubjects, dbObj.ProviderID)

					// we have all the info we need, now complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.DeleteAuthRelationships(ctx, "load-balancer", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpDelete|ent.OpDeleteOne,
		),
	}
}

func OriginHooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.OriginFunc(func(ctx context.Context, m *generated.OriginMutation) (ent.Value, error) {
					var err error
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					// addSubjPool, err := m.Client().Pool.Get(ctx, objID)
					addSubjPool, err := m.Client().Pool.Query().Where(pool.HasOriginsWith(origin.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPool.ID) && objID != addSubjPool.ID {
							additionalSubjects = append(additionalSubjects, addSubjPool.ID)
						}

						if !slices.Contains(additionalSubjects, addSubjPool.OwnerID) {
							additionalSubjects = append(additionalSubjects, addSubjPool.OwnerID)
						}
					}

					addSubjPort, err := m.Client().Port.Query().Where(port.HasPoolsWith(pool.HasOriginsWith(origin.IDEQ(objID)))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPort.LoadBalancerID) {
							additionalSubjects = append(additionalSubjects, addSubjPort.LoadBalancerID)
						}
					}

					changeset := []events.FieldChange{}
					cv_created_at := ""
					created_at, ok := m.CreatedAt()

					if ok {
						cv_created_at = created_at.Format(time.RFC3339)
						pv_created_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldCreatedAt(ctx)
							if err != nil {
								pv_created_at = "<unknown>"
							} else {
								pv_created_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "created_at",
							PreviousValue: pv_created_at,
							CurrentValue:  cv_created_at,
						})
					}

					cv_updated_at := ""
					updated_at, ok := m.UpdatedAt()

					if ok {
						cv_updated_at = updated_at.Format(time.RFC3339)
						pv_updated_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldUpdatedAt(ctx)
							if err != nil {
								pv_updated_at = "<unknown>"
							} else {
								pv_updated_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "updated_at",
							PreviousValue: pv_updated_at,
							CurrentValue:  cv_updated_at,
						})
					}

					cv_name := ""
					name, ok := m.Name()

					if ok {
						cv_name = fmt.Sprintf("%s", fmt.Sprint(name))
						pv_name := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldName(ctx)
							if err != nil {
								pv_name = "<unknown>"
							} else {
								pv_name = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "name",
							PreviousValue: pv_name,
							CurrentValue:  cv_name,
						})
					}

					cv_target := ""
					target, ok := m.Target()

					if ok {
						cv_target = fmt.Sprintf("%s", fmt.Sprint(target))
						pv_target := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldTarget(ctx)
							if err != nil {
								pv_target = "<unknown>"
							} else {
								pv_target = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "target",
							PreviousValue: pv_target,
							CurrentValue:  cv_target,
						})
					}

					cv_port_number := ""
					port_number, ok := m.PortNumber()

					if ok {
						cv_port_number = fmt.Sprintf("%s", fmt.Sprint(port_number))
						pv_port_number := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldPortNumber(ctx)
							if err != nil {
								pv_port_number = "<unknown>"
							} else {
								pv_port_number = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "port_number",
							PreviousValue: pv_port_number,
							CurrentValue:  cv_port_number,
						})
					}

					cv_active := ""
					active, ok := m.Active()

					if ok {
						cv_active = fmt.Sprintf("%s", fmt.Sprint(active))
						pv_active := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldActive(ctx)
							if err != nil {
								pv_active = "<unknown>"
							} else {
								pv_active = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "active",
							PreviousValue: pv_active,
							CurrentValue:  cv_active,
						})
					}

					cv_pool_id := ""
					pool_id, ok := m.PoolID()
					if !ok && !m.Op().Is(ent.OpCreate) {
						// since we are doing an update or delete and these fields didn't change, load the "old" value
						pool_id, err = m.OldPoolID(ctx)
						if err != nil {
							return nil, err
						}
					}
					additionalSubjects = append(additionalSubjects, pool_id)

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "loadbalancerorigin",
						SubjectID: pool_id,
					})

					if ok {
						cv_pool_id = fmt.Sprintf("%s", fmt.Sprint(pool_id))
						pv_pool_id := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldPoolID(ctx)
							if err != nil {
								pv_pool_id = "<unknown>"
							} else {
								pv_pool_id = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "pool_id",
							PreviousValue: pv_pool_id,
							CurrentValue:  cv_pool_id,
						})
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
						FieldChanges:         changeset,
					}

					// complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.CreateAuthRelationships(ctx, "load-balancer-origin", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer-origin", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),

		// Delete Hook
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.OriginFunc(func(ctx context.Context, m *generated.OriginMutation) (ent.Value, error) {
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					dbObj, err := m.Client().Origin.Get(ctx, objID)
					if err != nil {
						return nil, fmt.Errorf("failed to load object to get values for event, err %w", err)
					}

					additionalSubjects = append(additionalSubjects, dbObj.PoolID)

					addSubjPort, err := m.Client().Port.Query().Where(port.HasPoolsWith(pool.HasOriginsWith(origin.IDEQ(objID)))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPort.LoadBalancerID) {
							additionalSubjects = append(additionalSubjects, addSubjPort.LoadBalancerID)
						}
					}

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "loadbalancerorigin",
						SubjectID: dbObj.PoolID,
					})

					// we have all the info we need, now complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.DeleteAuthRelationships(ctx, "load-balancer-origin", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer-origin", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpDelete|ent.OpDeleteOne,
		),
	}
}

func PoolHooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.PoolFunc(func(ctx context.Context, m *generated.PoolMutation) (ent.Value, error) {
					var err error
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					addSubjPort, err := m.Client().Port.Query().Where(port.HasPoolsWith(pool.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPort.ID) && objID != addSubjPort.ID {
							additionalSubjects = append(additionalSubjects, addSubjPort.ID)
						}

						if !slices.Contains(additionalSubjects, addSubjPort.LoadBalancerID) {
							additionalSubjects = append(additionalSubjects, addSubjPort.LoadBalancerID)
						}
					}

					addSubjOrigin, err := m.Client().Origin.Query().Where(origin.HasPoolWith(pool.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjOrigin.ID) && objID != addSubjOrigin.ID {
							additionalSubjects = append(additionalSubjects, addSubjOrigin.ID)
						}

						if !slices.Contains(additionalSubjects, addSubjOrigin.PoolID) {
							additionalSubjects = append(additionalSubjects, addSubjOrigin.PoolID)
						}
					}

					changeset := []events.FieldChange{}
					cv_created_at := ""
					created_at, ok := m.CreatedAt()

					if ok {
						cv_created_at = created_at.Format(time.RFC3339)
						pv_created_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldCreatedAt(ctx)
							if err != nil {
								pv_created_at = "<unknown>"
							} else {
								pv_created_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "created_at",
							PreviousValue: pv_created_at,
							CurrentValue:  cv_created_at,
						})
					}

					cv_updated_at := ""
					updated_at, ok := m.UpdatedAt()

					if ok {
						cv_updated_at = updated_at.Format(time.RFC3339)
						pv_updated_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldUpdatedAt(ctx)
							if err != nil {
								pv_updated_at = "<unknown>"
							} else {
								pv_updated_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "updated_at",
							PreviousValue: pv_updated_at,
							CurrentValue:  cv_updated_at,
						})
					}

					cv_name := ""
					name, ok := m.Name()

					if ok {
						cv_name = fmt.Sprintf("%s", fmt.Sprint(name))
						pv_name := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldName(ctx)
							if err != nil {
								pv_name = "<unknown>"
							} else {
								pv_name = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "name",
							PreviousValue: pv_name,
							CurrentValue:  cv_name,
						})
					}

					cv_protocol := ""
					protocol, ok := m.Protocol()

					if ok {
						cv_protocol = fmt.Sprintf("%s", fmt.Sprint(protocol))
						pv_protocol := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldProtocol(ctx)
							if err != nil {
								pv_protocol = "<unknown>"
							} else {
								pv_protocol = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "protocol",
							PreviousValue: pv_protocol,
							CurrentValue:  cv_protocol,
						})
					}

					cv_owner_id := ""
					owner_id, ok := m.OwnerID()
					if !ok && !m.Op().Is(ent.OpCreate) {
						// since we are doing an update or delete and these fields didn't change, load the "old" value
						owner_id, err = m.OldOwnerID(ctx)
						if err != nil {
							return nil, err
						}
					}
					additionalSubjects = append(additionalSubjects, owner_id)

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "loadbalancerpool",
						SubjectID: owner_id,
					})

					if ok {
						cv_owner_id = fmt.Sprintf("%s", fmt.Sprint(owner_id))
						pv_owner_id := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldOwnerID(ctx)
							if err != nil {
								pv_owner_id = "<unknown>"
							} else {
								pv_owner_id = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "owner_id",
							PreviousValue: pv_owner_id,
							CurrentValue:  cv_owner_id,
						})
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
						FieldChanges:         changeset,
					}

					// complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.CreateAuthRelationships(ctx, "load-balancer-pool", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer-pool", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),

		// Delete Hook
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.PoolFunc(func(ctx context.Context, m *generated.PoolMutation) (ent.Value, error) {
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					dbObj, err := m.Client().Pool.Get(ctx, objID)
					if err != nil {
						return nil, fmt.Errorf("failed to load object to get values for event, err %w", err)
					}

					addSubjPort, err := m.Client().Port.Query().Where(port.HasPoolsWith(pool.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPort.LoadBalancerID) {
							additionalSubjects = append(additionalSubjects, addSubjPort.LoadBalancerID)
						}
					}

					additionalSubjects = append(additionalSubjects, dbObj.OwnerID)

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "loadbalancerpool",
						SubjectID: dbObj.OwnerID,
					})

					// we have all the info we need, now complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.DeleteAuthRelationships(ctx, "load-balancer-pool", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer-pool", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpDelete|ent.OpDeleteOne,
		),
	}
}

func PortHooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.PortFunc(func(ctx context.Context, m *generated.PortMutation) (ent.Value, error) {
					var err error
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					addSubjPool, err := m.Client().Pool.Query().Where(pool.HasPortsWith(port.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjPool.ID) && objID != addSubjPool.ID {
							additionalSubjects = append(additionalSubjects, addSubjPool.ID)
						}

						if !slices.Contains(additionalSubjects, addSubjPool.OwnerID) {
							additionalSubjects = append(additionalSubjects, addSubjPool.OwnerID)
						}
					}
					addSubjLoadBalancer, err := m.Client().LoadBalancer.Query().Where(loadbalancer.HasPortsWith(port.IDEQ(objID))).Only(ctx)
					if err == nil {
						if !slices.Contains(additionalSubjects, addSubjLoadBalancer.LocationID) {
							additionalSubjects = append(additionalSubjects, addSubjLoadBalancer.LocationID)
						}

						if !slices.Contains(additionalSubjects, addSubjLoadBalancer.OwnerID) {
							additionalSubjects = append(additionalSubjects, addSubjLoadBalancer.OwnerID)
						}

						if !slices.Contains(additionalSubjects, addSubjLoadBalancer.ProviderID) {
							additionalSubjects = append(additionalSubjects, addSubjLoadBalancer.ProviderID)
						}
					}

					changeset := []events.FieldChange{}
					cv_created_at := ""
					created_at, ok := m.CreatedAt()

					if ok {
						cv_created_at = created_at.Format(time.RFC3339)
						pv_created_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldCreatedAt(ctx)
							if err != nil {
								pv_created_at = "<unknown>"
							} else {
								pv_created_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "created_at",
							PreviousValue: pv_created_at,
							CurrentValue:  cv_created_at,
						})
					}

					cv_updated_at := ""
					updated_at, ok := m.UpdatedAt()

					if ok {
						cv_updated_at = updated_at.Format(time.RFC3339)
						pv_updated_at := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldUpdatedAt(ctx)
							if err != nil {
								pv_updated_at = "<unknown>"
							} else {
								pv_updated_at = ov.Format(time.RFC3339)
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "updated_at",
							PreviousValue: pv_updated_at,
							CurrentValue:  cv_updated_at,
						})
					}

					cv_number := ""
					number, ok := m.Number()

					if ok {
						cv_number = fmt.Sprintf("%s", fmt.Sprint(number))
						pv_number := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldNumber(ctx)
							if err != nil {
								pv_number = "<unknown>"
							} else {
								pv_number = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "number",
							PreviousValue: pv_number,
							CurrentValue:  cv_number,
						})
					}

					cv_name := ""
					name, ok := m.Name()

					if ok {
						cv_name = fmt.Sprintf("%s", fmt.Sprint(name))
						pv_name := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldName(ctx)
							if err != nil {
								pv_name = "<unknown>"
							} else {
								pv_name = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "name",
							PreviousValue: pv_name,
							CurrentValue:  cv_name,
						})
					}

					cv_load_balancer_id := ""
					load_balancer_id, ok := m.LoadBalancerID()
					if !ok && !m.Op().Is(ent.OpCreate) {
						// since we are doing an update or delete and these fields didn't change, load the "old" value
						load_balancer_id, err = m.OldLoadBalancerID(ctx)
						if err != nil {
							return nil, err
						}
					}
					additionalSubjects = append(additionalSubjects, load_balancer_id)

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "loadbalancer",
						SubjectID: load_balancer_id,
					})

					if ok {
						cv_load_balancer_id = fmt.Sprintf("%s", fmt.Sprint(load_balancer_id))
						pv_load_balancer_id := ""
						if !m.Op().Is(ent.OpCreate) {
							ov, err := m.OldLoadBalancerID(ctx)
							if err != nil {
								pv_load_balancer_id = "<unknown>"
							} else {
								pv_load_balancer_id = fmt.Sprintf("%s", fmt.Sprint(ov))
							}
						}

						changeset = append(changeset, events.FieldChange{
							Field:         "load_balancer_id",
							PreviousValue: pv_load_balancer_id,
							CurrentValue:  cv_load_balancer_id,
						})
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
						FieldChanges:         changeset,
					}

					// complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.CreateAuthRelationships(ctx, "load-balancer-port", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer-port", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),

		// Delete Hook
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.PortFunc(func(ctx context.Context, m *generated.PortMutation) (ent.Value, error) {
					additionalSubjects := []gidx.PrefixedID{}
					relationships := []events.AuthRelationshipRelation{}

					objID, ok := m.ID()
					if !ok {
						return nil, fmt.Errorf("object doesn't have an id %s", objID)
					}

					dbObj, err := m.Client().Port.Get(ctx, objID)
					if err != nil {
						return nil, fmt.Errorf("failed to load object to get values for event, err %w", err)
					}

					additionalSubjects = append(additionalSubjects, dbObj.LoadBalancerID)

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "loadbalancer",
						SubjectID: dbObj.LoadBalancerID,
					})

					addSubjLoadBalancer, err := m.Client().LoadBalancer.Get(ctx, dbObj.LoadBalancerID)
					if err != nil {
						if !slices.Contains(additionalSubjects, addSubjLoadBalancer.LocationID) {
							additionalSubjects = append(additionalSubjects, addSubjLoadBalancer.LocationID)
						}

						if !slices.Contains(additionalSubjects, addSubjLoadBalancer.OwnerID) {
							additionalSubjects = append(additionalSubjects, addSubjLoadBalancer.OwnerID)
						}

						if !slices.Contains(additionalSubjects, addSubjLoadBalancer.ProviderID) {
							additionalSubjects = append(additionalSubjects, addSubjLoadBalancer.ProviderID)
						}
					}

					// we have all the info we need, now complete the mutation before we process the event
					retValue, err := next.Mutate(ctx, m)
					if err != nil {
						return retValue, err
					}

					if len(relationships) != 0 {
						if err := permissions.DeleteAuthRelationships(ctx, "load-balancer-port", objID, relationships...); err != nil {
							return nil, fmt.Errorf("relationship request failed with error: %w", err)
						}
					}

					msg := events.ChangeMessage{
						EventType:            eventType(m.Op()),
						SubjectID:            objID,
						AdditionalSubjectIDs: additionalSubjects,
						Timestamp:            time.Now().UTC(),
					}

					if _, err := m.EventsPublisher.PublishChange(ctx, "load-balancer-port", msg); err != nil {
						return nil, fmt.Errorf("failed to publish change: %w", err)
					}

					return retValue, nil
				})
			},
			ent.OpDelete|ent.OpDeleteOne,
		),
	}
}

// PubsubHooks registers our hooks with the ent client
func PubsubHooks(c *generated.Client) {
	c.LoadBalancer.Use(LoadBalancerHooks()...)

	c.Origin.Use(OriginHooks()...)

	c.Pool.Use(PoolHooks()...)

	c.Port.Use(PortHooks()...)
}

func eventType(op ent.Op) string {
	switch op {
	case ent.OpCreate:
		return string(events.CreateChangeType)
	case ent.OpUpdate, ent.OpUpdateOne:
		return string(events.UpdateChangeType)
	case ent.OpDelete, ent.OpDeleteOne:
		return string(events.DeleteChangeType)
	default:
		return "unknown"
	}
}
