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

CREATE TABLE usersets (
    id serial PRIMARY KEY,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    notebooks_id int REFERENCES notebooks(id) ON DELETE CASCADE,
    codeinsights_id int REFERENCES codeinsights(id) ON DELETE CASCADE,
    groups_id int REFERENCES groups(id) ON DELETE CASCADE,
    UNIQUE (relation, notebooks_id),
    UNIQUE (relation, codeinsights_id),
    UNIQUE (relation, groups_id)
);

CREATE TABLE groups_namespace (
    id int REFERENCES groups(id) ON DELETE CASCADE,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    userset_id int REFERENCES usersets(id) ON DELETE CASCADE,
    UNIQUE (id, relation, user_id),
    UNIQUE (id, relation, userset_id)
);

CREATE TABLE notebooks_namespace (
    id int REFERENCES notebooks(id) ON DELETE CASCADE,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    userset_id int REFERENCES usersets(id) ON DELETE CASCADE,
    UNIQUE (id, relation, user_id),
    UNIQUE (id, relation, userset_id)
);

CREATE TABLE codeinsights_namespace (
    id int REFERENCES codeinsights(id) ON DELETE CASCADE,
    relation int REFERENCES relations(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    userset_id int REFERENCES usersets(id) ON DELETE CASCADE,
    UNIQUE (id, relation, user_id),
    UNIQUE (id, relation, userset_id)
);