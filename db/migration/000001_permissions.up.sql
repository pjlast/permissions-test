CREATE TABLE namespaces (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE relations (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE groups (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE notebooks (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE codeinsights (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE resource_mapping (
    id serial PRIMARY KEY,
    notebooks_id int UNIQUE REFERENCES notebooks(id) ON DELETE CASCADE,
    codeinsights_id int UNIQUE REFERENCES codeinsights(id) ON DELETE CASCADE,
    groups_id int UNIQUE REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE usersets (
    id serial PRIMARY KEY,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    resource_id int REFERENCES resource_mapping(id) ON DELETE CASCADE,
    UNIQUE (relation, resource_id)
);

CREATE TABLE groups_namespace (
    id int REFERENCES resource_mapping(id) ON DELETE CASCADE,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    userset_id int REFERENCES usersets(id) ON DELETE CASCADE,
    UNIQUE (id, relation, user_id),
    UNIQUE (id, relation, userset_id)
);

CREATE TABLE notebooks_namespace (
    id int REFERENCES resource_mapping(id) ON DELETE CASCADE,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    userset_id int REFERENCES usersets(id) ON DELETE CASCADE,
    UNIQUE (id, relation, user_id),
    UNIQUE (id, relation, userset_id)
);

CREATE TABLE codeinsights_namespace (
    id int REFERENCES resource_mapping(id) ON DELETE CASCADE,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    userset_id int REFERENCES usersets(id) ON DELETE CASCADE,
    UNIQUE (id, relation, user_id),
    UNIQUE (id, relation, userset_id)
);