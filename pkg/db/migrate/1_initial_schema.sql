-- +migrate Up

/* enable CI Text extension */
CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

CREATE TABLE packages (
  "id" serial primary key,
  "name" citext NOT NULL
);

CREATE UNIQUE INDEX packages_name_key ON packages(name);

CREATE TABLE installed_packages (
  "package_id" int primary key
);
ALTER TABLE installed_packages
add constraint package_id_fkey
    foreign key (package_id)
    references packages(id)
    on delete cascade;



CREATE TABLE package_dependencies (
  "package_id" int,
  "needed_package_id" int,
  PRIMARY KEY(package_id, needed_package_id)
);


ALTER TABLE package_dependencies
add constraint package_id_fkey
    foreign key (package_id)
    references packages(id)
    on delete cascade,
add constraint needed_package_id_fkey
    foreign key (needed_package_id)
    references packages(id)
    on delete cascade;


-- +migrate Down
DROP TABLE package_dependencies;
DROP TABLE installed_packages;
DROP TABLE packages;

