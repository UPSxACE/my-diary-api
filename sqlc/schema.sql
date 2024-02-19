CREATE TABLE "user" (
    id int GENERATED ALWAYS AS IDENTITY NOT NULL,
    username varchar(62) NOT NULL,
    "password" varchar(254) NOT NULL,
    email varchar(254) NOT NULL,
    avatar_url varchar(254) NULL,
    full_name varchar(254) NULL,
    skip_tutorials boolean DEFAULT false NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    deleted_at timestamp NULL,
    CONSTRAINT user_pk PRIMARY KEY (id),
    CONSTRAINT user_unique UNIQUE (username),
    CONSTRAINT user_unique_1 UNIQUE (email)
);
COMMENT ON TABLE "user" IS 'Table to store user accounts';

CREATE TABLE user_tutorial (
    id int GENERATED ALWAYS AS IDENTITY NOT NULL,
    user_id int NOT NULL,
    welcome boolean DEFAULT false NOT NULL,
    CONSTRAINT user_tutorial_pk PRIMARY KEY (id),
    CONSTRAINT user_tutorial_unique UNIQUE (user_id),
    CONSTRAINT user_tutorial_user_fk FOREIGN KEY (user_id) REFERENCES "user"(id)
);
COMMENT ON TABLE user_tutorial IS 'Table that stores whether or not an user has completed certain tutorials';

CREATE TABLE note (
    id int GENERATED ALWAYS AS IDENTITY NOT NULL,
    author_id int NULL,
    title varchar(254) NOT NULL,
    "content" varchar(131070) NOT NULL,
    content_raw varchar(131070) NOT NULL,
    "views" int DEFAULT 0 NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    lastread_at timestamp NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    deleted_at timestamp NULL,
    CONSTRAINT note_pk PRIMARY KEY (id),
    CONSTRAINT note_user_fk FOREIGN KEY (author_id) REFERENCES "user"(id)
);
COMMENT ON TABLE note IS 'Table to store user''s notes';

CREATE TABLE note_change (
    id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
    note_id int NOT NULL,
    title varchar(254) NULL,
    "content" varchar(131070) NULL,
    valid_until timestamp NOT NULL,
    CONSTRAINT note_change_note_fk FOREIGN KEY (note_id) REFERENCES note(id)
);
COMMENT ON TABLE note_change IS 'Table that stores the previous versions of the user''s notes and when they got changed';