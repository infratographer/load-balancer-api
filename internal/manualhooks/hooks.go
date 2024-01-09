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
	"go.infratographer.com/load-balancer-api/internal/ent/schema"
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

					// Ensure we have additional relevant subjects in the msg
					lb, err := m.Client().LoadBalancer.Query().WithPorts().Where(loadbalancer.IDEQ(objID)).Only(ctx)
					if err == nil {
						if !slices.Contains(msg.AdditionalSubjectIDs, lb.LocationID) {
							msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, lb.LocationID)
						}

						if !slices.Contains(msg.AdditionalSubjectIDs, lb.ProviderID) {
							msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, lb.ProviderID)
						}

						for _, p := range lb.Edges.Ports {
							if !slices.Contains(msg.AdditionalSubjectIDs, p.ID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, p.ID)
							}
						}
					}

					if len(relationships) != 0 && eventType(m.Op()) == string(events.CreateChangeType) {
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

					// Ensure we have additional relevant subjects in the msg
					addSubjPorts, err := m.Client().Port.Query().WithPools().WithLoadBalancer().Where(port.HasPoolsWith(pool.HasOriginsWith(origin.IDEQ(objID)))).All(ctx)
					if err == nil {
						for _, port := range addSubjPorts {
							if !slices.Contains(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.LocationID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.LocationID)
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.ProviderID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.ProviderID)
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, port.LoadBalancerID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.LoadBalancerID)
							}

							for _, pool := range port.Edges.Pools {
								if !slices.Contains(msg.AdditionalSubjectIDs, pool.ID) {
									msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, pool.ID)
								}

								if !slices.Contains(msg.AdditionalSubjectIDs, pool.OwnerID) {
									msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, pool.OwnerID)
								}
							}
						}
					}

					if len(relationships) != 0 && eventType(m.Op()) == string(events.CreateChangeType) {
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

					// Ensure we have additional relevant subjects in the msg
					addSubjPorts, err := m.Client().Port.Query().WithPools().WithLoadBalancer().Where(port.HasPoolsWith(pool.HasOriginsWith(origin.IDEQ(objID)))).All(ctx)
					if err == nil {
						for _, port := range addSubjPorts {
							for _, pool := range port.Edges.Pools {
								if !slices.Contains(additionalSubjects, pool.ID) {
									additionalSubjects = append(additionalSubjects, pool.ID)
								}

								if !slices.Contains(additionalSubjects, pool.OwnerID) {
									additionalSubjects = append(additionalSubjects, pool.OwnerID)
								}
							}

							if !slices.Contains(additionalSubjects, port.LoadBalancerID) {
								additionalSubjects = append(additionalSubjects, port.LoadBalancerID)
							}

							if !slices.Contains(additionalSubjects, port.Edges.LoadBalancer.LocationID) {
								additionalSubjects = append(additionalSubjects, port.Edges.LoadBalancer.LocationID)
							}

							if !slices.Contains(additionalSubjects, port.Edges.LoadBalancer.ProviderID) {
								additionalSubjects = append(additionalSubjects, port.Edges.LoadBalancer.ProviderID)
							}
						}
					}

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

					// Ensure we have additional relevant subjects in the msg
					addSubjPorts, err := m.Client().Port.Query().WithLoadBalancer().WithPools(func(q *generated.PoolQuery) {
						q.WithOrigins()
					}).Where(port.HasPoolsWith(pool.IDEQ(objID))).All(ctx)
					if err == nil {
						for _, port := range addSubjPorts {
							for _, pool := range port.Edges.Pools {
								for _, origin := range pool.Edges.Origins {
									if !slices.Contains(msg.AdditionalSubjectIDs, origin.ID) {
										msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, origin.ID)
									}
								}
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, port.ID) && objID != port.ID {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.ID)
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, port.LoadBalancerID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.LoadBalancerID)
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.LocationID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.LocationID)
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.ProviderID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, port.Edges.LoadBalancer.ProviderID)
							}
						}
					}

					if len(relationships) != 0 && eventType(m.Op()) == string(events.CreateChangeType) {
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

					additionalSubjects = append(additionalSubjects, dbObj.OwnerID)

					// Ensure we have additional relevant subjects in the msg
					addSubjPorts, err := m.Client().Port.Query().WithLoadBalancer().Where(port.HasPoolsWith(pool.IDEQ(objID))).All(ctx)
					if err == nil {
						for _, port := range addSubjPorts {
							if !slices.Contains(additionalSubjects, port.Edges.LoadBalancer.LocationID) {
								additionalSubjects = append(additionalSubjects, port.Edges.LoadBalancer.LocationID)
							}

							if !slices.Contains(additionalSubjects, port.Edges.LoadBalancer.ProviderID) {
								additionalSubjects = append(additionalSubjects, port.Edges.LoadBalancer.ProviderID)
							}

							if !slices.Contains(additionalSubjects, port.LoadBalancerID) {
								additionalSubjects = append(additionalSubjects, port.LoadBalancerID)
							}
						}
					}

					relationships = append(relationships, events.AuthRelationshipRelation{
						Relation:  "owner",
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

					// Ensure we have additional relevant subjects in the event msg
					addSubjPort, err := m.Client().Port.Query().WithPools().WithLoadBalancer().Where(port.IDEQ(objID)).Only(ctx)
					if err == nil {
						if !slices.Contains(msg.AdditionalSubjectIDs, addSubjPort.LoadBalancerID) {
							msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, addSubjPort.LoadBalancerID)
						}

						if !slices.Contains(msg.AdditionalSubjectIDs, addSubjPort.Edges.LoadBalancer.LocationID) {
							msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, addSubjPort.Edges.LoadBalancer.LocationID)
						}

						if !slices.Contains(msg.AdditionalSubjectIDs, addSubjPort.Edges.LoadBalancer.OwnerID) {
							msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, addSubjPort.Edges.LoadBalancer.OwnerID)
						}

						if !slices.Contains(msg.AdditionalSubjectIDs, addSubjPort.Edges.LoadBalancer.ProviderID) {
							msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, addSubjPort.Edges.LoadBalancer.ProviderID)
						}

						for _, pool := range addSubjPort.Edges.Pools {
							if !slices.Contains(msg.AdditionalSubjectIDs, pool.ID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, pool.ID)
							}

							if !slices.Contains(msg.AdditionalSubjectIDs, pool.OwnerID) {
								msg.AdditionalSubjectIDs = append(msg.AdditionalSubjectIDs, pool.OwnerID)
							}
						}
					}

					if len(relationships) != 0 && eventType(m.Op()) == string(events.CreateChangeType) {
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

					dbObj, err := m.Client().Port.Query().WithLoadBalancer().Where(port.IDEQ(objID)).Only(ctx)
					if err != nil {
						return nil, fmt.Errorf("failed to load object to get values for event, err %w", err)
					}

					// Ensure we have additional relevant subjects in the event msg
					additionalSubjects = append(additionalSubjects, dbObj.LoadBalancerID)
					additionalSubjects = append(additionalSubjects, dbObj.Edges.LoadBalancer.LocationID)
					additionalSubjects = append(additionalSubjects, dbObj.Edges.LoadBalancer.OwnerID)
					additionalSubjects = append(additionalSubjects, dbObj.Edges.LoadBalancer.ProviderID)

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

func getLoadBalancerIDs(ctx context.Context, id gidx.PrefixedID, addID []gidx.PrefixedID) []gidx.PrefixedID {
	lbIDs := []gidx.PrefixedID{}

	if id.Prefix() == schema.LoadBalancerPrefix {
		lbIDs = append(lbIDs, id)
	}

	for _, id := range addID {
		if id.Prefix() == schema.LoadBalancerPrefix {
			lbIDs = append(lbIDs, id)
		}
	}

	return lbIDs
}
