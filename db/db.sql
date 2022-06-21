CREATE UNLOGGED TABLE Users
(
    id       serial       NOT NULL,
    nickname varchar(255) NOT NULL UNIQUE,
    fullname varchar(255) NOT NULL,
    about    text         NOT NULL,
    email    varchar(255) NOT NULL UNIQUE,
    CONSTRAINT UserPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE Forum
(
    id            serial       NOT NULL,
    title         varchar(255) NOT NULL,
    user_nickname varchar(255) NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    slug          varchar(255) NOT NULL UNIQUE,
    posts         int          NOT NULL default 0,
    threads       int          NOT NULL default 0,
    CONSTRAINT ForumPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE Thread
(
    id            serial       NOT NULL,
    title         varchar(255) NOT NULL,
    user_nickname varchar(255) NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    forum         varchar(255) NOT NULL REFERENCES Forum (slug) ON DELETE CASCADE,
    message       text         NOT NULL,
    votes         int                      default 0,
    slug          varchar(255) NOT NULL,
    created       timestamp with time zone DEFAULT now(),
    CONSTRAINT ThreadPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE Post
(
    id        serial       NOT NULL,
    parent    int          NOT NULL    default 0,
    author    varchar(255) NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    message   text         NOT NULL,
    is_edited boolean      NOT NULL,
    forum     varchar(255) NOT NULL REFERENCES Forum (slug) ON DELETE CASCADE,
    thread    int          NOT NULL,
    created   timestamp with time zone DEFAULT now(),
    pathTree  bigint[]                 default array []::bigint[],
    CONSTRAINT PostPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE Vote
(
    threadId int          NOT NULL REFERENCES Thread (id) ON DELETE CASCADE,
    nickname varchar(255) NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    voice    int          NOT NULL default 0
);

CREATE UNLOGGED TABLE Stat
(
    posts int default 0,
    users int default 0,
    forums int default 0,
    threads int default 0
);

INSERT into Stat
VALUES (0, 0, 0, 0);

CREATE UNIQUE INDEX lowercase_emails on Users (lower(email));
CREATE UNIQUE INDEX lowercase_nicknames on Users (lower(nickname));
CREATE UNIQUE INDEX lowercase_slug_forum on Forum (lower(slug));
CREATE INDEX if not exists lowercase_slug_thread on Thread (lower(slug));
CREATE UNIQUE INDEX unique_vote on Vote (threadId, nickname);

CREATE INDEX IF NOT EXISTS thread_parentTree_post on post (thread,pathtree);
CREATE INDEX IF NOT EXISTS first_parent_post on post ((pathtree[1]),pathtree);

CREATE OR REPLACE FUNCTION insertPathTree() RETURNS trigger as
$insertPathTree$
Declare
    parent_path BIGINT[];
    parent_thread int;
    parent_post int;
--     post_author Users;
begin
    if (new.parent = 0) then
        new.pathtree := array_append(new.pathtree, new.id);
    else
--         select * from post where post.id = new.parent into parent_post;
--         if (parent_post is null) then
--             raise exception sqlstate '22003' using message='parent id is not exist';
--         end if;

--         select * from users where lower(users.nickname) = lower(new.author) into post_author;
--         if (post_author is null) then
--             raise exception sqlstate '22003' using message='author is not exist';
--         end if;

        select thread, pathTree, id from post where id = new.parent into parent_thread, parent_path, parent_post;
        if (parent_post is null) then
            raise exception sqlstate '22003' using message='parent id is not exist';
        end if;

        if (parent_thread != new.thread) then
            raise exception sqlstate '22003' using message='threads not equals';
        end if;

        new.pathtree := new.pathtree || parent_path || new.id;

    end if;
    UPDATE forum SET posts=posts + 1 WHERE lower(forum.slug) = lower(new.forum);
    UPDATE Stat SET posts = Stat.posts + 1;
    return new;
end
$insertPathTree$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION incThreadCounter() RETURNS trigger as
$incThreadCounter$
begin
    UPDATE forum SET threads = threads + 1 WHERE lower(forum.slug) = lower(new.forum);
    UPDATE Stat SET threads = Stat.threads + 1;
    return new;
end;
$incThreadCounter$ LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION incPostsCounter() RETURNS trigger as
-- $incPostsCounter$
-- begin
--     UPDATE Stat SET posts = Stat.posts + 1;
--     return new;
-- end;
-- $incPostsCounter$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION incForumsCounter() RETURNS trigger as
$incForumsCounter$
begin
    UPDATE Stat SET forums = Stat.forums + 1;
    return new;
end;
$incForumsCounter$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION incUsersCounter() RETURNS trigger as
$incUsersCounter$
begin
    UPDATE Stat SET users = Stat.users + 1;
    return new;
end;
$incUsersCounter$ LANGUAGE plpgsql;


CREATE TRIGGER insertPathTreeTrigger
    BEFORE INSERT
    on Post
    for each row
EXECUTE Function insertPathTree();

CREATE TRIGGER incThreadCounter
    AFTER INSERT
    on Thread
    for each row
EXECUTE FUNCTION incThreadCounter();

CREATE TRIGGER incUsersCounter
    AFTER INSERT
    on Users
    for each row
EXECUTE FUNCTION incUsersCounter();

CREATE TRIGGER incForumCounter
    AFTER INSERT
    on Forum
    for each row
EXECUTE FUNCTION incForumsCounter();

VACUUM ANALYZE;
