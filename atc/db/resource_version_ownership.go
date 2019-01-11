package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/concourse/atc"
)

type ResourceVersionOwnership interface {
	ID() int
	ResourceConfig() ResourceConfig

	SaveUncheckedVersion(version atc.Version, metadata ResourceConfigMetadataFields) (bool, error)
	SaveVersions(versions []atc.Version) error
	FindVersion(atc.Version) (ResourceConfigVersion, bool, error)
	LatestVersion() (ResourceConfigVersion, bool, error)
}

type resourceVersionOwnership struct {
	id             int
	resourceConfig ResourceConfig

	conn Conn
}

func (r *resourceVersionOwnership) ID() int                        { return r.id }
func (r *resourceVersionOwnership) ResourceConfig() ResourceConfig { return r.resourceConfig }

// SaveUncheckedVersion is used by the "get" and "put" step to find or create of a
// resource config version. We want to do an upsert because there will be cases
// where resource config versions can become outdated while the versions
// associated to it are still valid. This will be special case where we save
// the version with a check order of 0 in order to avoid using this version
// until we do a proper check. Note that this method will not bump the cache
// index for the pipeline because we want to ignore these versions until the
// check orders get updated. The bumping of the index will be done in
// SaveOutput for the put step.
func (r *resourceVersionOwnership) SaveUncheckedVersion(version atc.Version, metadata ResourceConfigMetadataFields) (bool, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return false, err
	}

	defer Rollback(tx)

	newVersion, err := saveResourceVersion(tx, r, version, metadata)
	if err != nil {
		return false, err
	}

	return newVersion, tx.Commit()
}

// SaveVersions stores a list of version in the db for a resource config
// Each version will also have its check order field updated and the
// Cache index for pipelines using the resource config will be bumped.
//
// In the case of a check resource from an older version, the versions
// that already exist in the DB will be re-ordered using
// incrementCheckOrderWhenNewerVersion to input the correct check order
func (r *resourceVersionOwnership) SaveVersions(versions []atc.Version) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	defer Rollback(tx)

	for _, version := range versions {
		_, err = saveResourceVersion(tx, r, version, nil)
		if err != nil {
			return err
		}

		versionJSON, err := json.Marshal(version)
		if err != nil {
			return err
		}

		err = incrementCheckOrder(tx, r, string(versionJSON))
		if err != nil {
			return err
		}
	}

	err = bumpCacheIndexForPipelinesUsingResourceConfig(tx, r.id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *resourceVersionOwnership) FindVersion(v atc.Version) (ResourceConfigVersion, bool, error) {
	rcv := &resourceConfigVersion{
		resourceVersionOwnership: r,
		conn:                     r.conn,
	}

	versionByte, err := json.Marshal(v)
	if err != nil {
		return nil, false, err
	}

	row := resourceConfigVersionQuery.
		Where(sq.Eq{
			"v.resource_config_id": r.id,
		}).
		Where(sq.Expr(fmt.Sprintf("v.version_md5 = md5('%s')", versionByte))).
		RunWith(r.conn).
		QueryRow()

	err = scanResourceConfigVersion(rcv, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return rcv, true, nil
}

func (r *resourceVersionOwnership) LatestVersion() (ResourceConfigVersion, bool, error) {
	rcv := &resourceConfigVersion{
		conn:                     r.conn,
		resourceVersionOwnership: r,
	}

	row := resourceConfigVersionQuery.
		Where(sq.Eq{"v.resource_config_id": r.id}).
		OrderBy("v.check_order DESC").
		Limit(1).
		RunWith(r.conn).
		QueryRow()

	err := scanResourceConfigVersion(rcv, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return rcv, true, nil
}

func saveResourceVersion(tx Tx, r ResourceVersionOwnership, version atc.Version, metadata ResourceConfigMetadataFields) (bool, error) {
	versionJSON, err := json.Marshal(version)
	if err != nil {
		return false, err
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return false, err
	}

	var checkOrder int
	err = tx.QueryRow(`
		INSERT INTO resource_config_versions (resource_version_ownership_id, version, version_md5, metadata)
		SELECT $1, $2, md5($3), $4
		ON CONFLICT (resource_version_ownership_id, version_md5) DO UPDATE SET metadata = $4
		RETURNING check_order
		`, r.ID(), string(versionJSON), string(versionJSON), string(metadataJSON)).Scan(&checkOrder)
	if err != nil {
		return false, err
	}

	return checkOrder == 0, nil
}

// increment the check order if the version's check order is less than the
// current max. This will fix the case of a check from an old version causing
// the desired order to change; existing versions will be re-ordered since
// we add them in the desired order.
func incrementCheckOrder(tx Tx, r ResourceVersionOwnership, version string) error {
	_, err := tx.Exec(`
		WITH max_checkorder AS (
			SELECT max(check_order) co
			FROM resource_config_versions
			WHERE resource_version_ownership_id = $1
		)

		UPDATE resource_config_versions
		SET check_order = mc.co + 1
		FROM max_checkorder mc
		WHERE resource_version_ownership_id = $1
		AND version = $2
		AND check_order <= mc.co;`, r.ID(), version)
	return err
}
