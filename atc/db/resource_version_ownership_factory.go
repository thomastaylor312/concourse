package db

import sq "github.com/Masterminds/squirrel"

type ResourceVersionOwnershipFactory interface {
	FindOrCreateResourceVersionOwnership(ResourceConfig, Resource) (ResourceVersionOwnership, error)
}

type resourceVersionOwnershipFactory struct {
	conn Conn
}

func NewResourceVersionOwnershipFactory(conn Conn) ResourceVersionOwnershipFactory {
	return &resourceVersionOwnershipFactory{
		conn: conn,
	}
}

func (r *resourceVersionOwnershipFactory) FindOrCreateResourceVersionOwnership(resourceConfig ResourceConfig, resource Resource) (ResourceVersionOwnership, error) {
	//Select and Delete
	var resourceID *int
	// if unique
	if resourceConfig.UniqueVersionHistory() {
		rid := resource.ID()
		resourceID = &rid
		// check for existing RVO entry
		var rcID int
		rows, err := psql.Select("id, resource_config_id").
			From("resource_version_ownership").
			Where(sq.Eq{
				"resource_id": resourceID,
			}).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, err
		}

		if rows.Next() {
			var ownerID int
			rows.Scan(&ownerID, &rcID)
			// ensure that there is an entry for the matching resource_config_id
			if rcID != resourceConfig.ID() {
				_, err := psql.Delete("resource_version_ownership").
					Where(sq.And{
						sq.Eq{
							"resource_id": resourceID,
						},
						sq.NotEq{
							"resource_config_id": resourceConfig.ID(),
						},
					}).
					RunWith(tx).
					Exec()
				if err != nil {
					return nil, err
				}
			} else {
				return &resourceVersionOwnership{ownerID, resourceConfig, r.conn}, nil
			}
		}
	}

	var ownerID int
	// -> create a new entry for the resource_id, resource_config_id
	err = psql.Insert("resource_version_ownership").
		Columns("resource_id", "resource_config_id").
		Values(resourceID, resourceConfig.ID()).
		Suffix(`
			ON CONFLICT (resource_id, resource_config_id) DO UPDATE SET
				resource_id = ?,
				resource_config_id = ?,
			RETURNING id
		`, resourceID, resourceConfig.ID()).
		RunWith(tx).
		QueryRow().
		Scan(&ownerID)
	if err != nil {
		return nil, err
	}
}
