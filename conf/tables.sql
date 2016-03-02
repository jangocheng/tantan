CREATE TABLE user_info(
    user_id serial NOT NULL PRIMARY KEY,
    user_name VARCHAR(64) NOT NULL,
    type VARCHAR(32) NOT NULL,
    created_utc INT NOT NULL,
    status INT NOT NULL
);

select setval('user_info_user_id_seq',1000000);
CREATE UNIQUE INDEX uix_user_info_name ON user_info(user_name);


CREATE TABLE relationship(
    id serial NOT NULL PRIMARY KEY,
    master BIGINT NOT NULL,
    liker BIGINT NOT NULL,
    type INT NOT NULL,
    state INT NOT NULL,   -- 0代表dislike 1 代表liked, 2 代表match
    created_utc INT NOT NULL,
    status INT NOT NULL
);

CREATE UNIQUE INDEX uix_relationship_master_linker ON relationship(master,liker);
