BEGIN;

  CREATE TABLE resource_version_ownership (
    "id" serial NOT NULL PRIMARY KEY,
    "resource_config_id" integer NOT NULL REFERENCES resource_configs (id) ON DELETE CASCADE,
    "resource_id" integer REFERENCES resources (id) ON DELETE CASCADE
  );

  CREATE UNIQUE INDEX resource_version_ownership_resource_id_uniq
  ON resource_version_ownership (resource_id);

  CREATE UNIQUE INDEX resource_version_ownership_global
  ON resource_version_ownership (resource_config_id)
  WHERE resource_id IS NULL;

  ALTER TABLE resource_configs
    DROP COLUMN unique_versions_resource_id;

  CREATE UNIQUE INDEX resource_configs_resource_cache_id_so_key
  ON resource_configs (resource_cache_id, source_hash);

  CREATE UNIQUE INDEX resource_configs_base_resource_type_id_so_key
  ON resource_configs (base_resource_type_id, source_hash);

COMMIT;
