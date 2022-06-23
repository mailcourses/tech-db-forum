CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE Users
(
    id       serial             NOT NULL,
    nickname citext COLLATE "C" NOT NULL UNIQUE PRIMARY KEY,
    fullname varchar(255)       NOT NULL,
    about    text               NOT NULL,
    email    citext             NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE Forum
(
    id            serial             NOT NULL,
    title         varchar(255)       NOT NULL,
    user_nickname citext COLLATE "C" NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    slug          citext             NOT NULL UNIQUE,
    posts         int                NOT NULL default 0,
    threads       int                NOT NULL default 0,
    CONSTRAINT ForumPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE ForumUsers
(
    nickname citext COLLATE "C" NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    fullname varchar(255)       NOT NULL,
    about    text               NOT NULL,
    email    citext             NOT NULL REFERENCES Users (email) ON DELETE CASCADE,
    forum    citext             NOT NULL REFERENCES Forum (slug) ON DELETE CASCADE
);

CREATE UNLOGGED TABLE Thread
(
    id            serial             NOT NULL,
    title         varchar(255)       NOT NULL,
    user_nickname citext COLLATE "C" NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    forum         citext             NOT NULL REFERENCES Forum (slug) ON DELETE CASCADE,
    message       text               NOT NULL,
    votes         int         default 0,
    slug          citext             NOT NULL,
    created       timestamptz DEFAULT now(),
    CONSTRAINT ThreadPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE Post
(
    id        serial             NOT NULL,
    parent    int                NOT NULL default 0,
    author    citext COLLATE "C" NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    message   text               NOT NULL,
    is_edited boolean            NOT NULL,
    forum     citext             NOT NULL REFERENCES Forum (slug) ON DELETE CASCADE,
    thread    int                NOT NULL,
    created   timestamptz                 DEFAULT now(),
    pathTree  bigint[]                    default array []::bigint[],
    CONSTRAINT PostPk PRIMARY KEY (id)
);

CREATE UNLOGGED TABLE Vote
(
    threadId int                NOT NULL REFERENCES Thread (id) ON DELETE CASCADE,
    nickname citext COLLATE "C" NOT NULL REFERENCES Users (nickname) ON DELETE CASCADE,
    voice    int                NOT NULL default 0
);

-- Users
CREATE INDEX IF NOT EXISTS index_emails_users ON Users USING hash (email);
CREATE INDEX IF NOT EXISTS index_nicknames_hash_users ON Users USING hash (nickname);
CREATE INDEX IF NOT EXISTS index_nicknames_btree_users ON Users USING btree (nickname);
CREATE UNIQUE INDEX IF NOT EXISTS index_lowercase_emails on Users (lower(email));
CREATE UNIQUE INDEX IF NOT EXISTS index_lowercase_nicknames on Users (lower(nickname));
CREATE UNIQUE INDEX IF NOT EXISTS index_unique_users on Users (email, nickname);

-- Forum
CREATE INDEX IF NOT EXISTS index_user_forum on Forum using hash (user_nickname);
CREATE INDEX IF NOT EXISTS index_slug_forum on Forum using hash (slug);
CREATE UNIQUE INDEX IF NOT EXISTS index_lowercase_slug_forum on Forum (lower(slug));

-- ForumUsers
CREATE UNIQUE INDEX IF NOT EXISTS index_uniq_user on ForumUsers (lower(email), lower(nickname), about, lower(forum), fullname);


-- Threads
CREATE INDEX IF NOT EXISTS index_forum_thread on Thread using hash (forum);
CREATE INDEX IF NOT EXISTS index_slug_thread on Thread using hash (slug);
CREATE INDEX IF NOT EXISTS index_user_thread on Thread using hash (user_nickname);

-- Post
CREATE INDEX IF NOT EXISTS index_author_post on Post using hash (author);
CREATE INDEX IF NOT EXISTS index_forum_post on Post using hash (forum);
CREATE INDEX IF NOT EXISTS index_parent_post on Post using btree (parent);
CREATE INDEX IF NOT EXISTS thread_parentTree_post on post (thread, pathtree);
CREATE INDEX IF NOT EXISTS first_parent_post on post ((pathtree[1]), pathtree);

-- Vote
CREATE INDEX IF NOT EXISTS index_nickname_vote on Vote using hash (nickname);
CREATE INDEX IF NOT EXISTS index_threadId_vote on Vote using hash (threadid);
CREATE UNIQUE INDEX IF NOT EXISTS index_unique_vote on Vote (threadId, nickname);


CREATE OR REPLACE FUNCTION insertPathTree() RETURNS trigger as
$insertPathTree$
Declare
    parent_path BIGINT[];
begin
    if (new.parent = 0) then
        new.pathtree := array_append(new.pathtree, new.id);
    else
        select pathtree from post where id = new.parent into parent_path;
        new.pathtree := new.pathtree || parent_path || new.id;
    end if;
    return new;
end
$insertPathTree$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION insertThreadsVotes() RETURNS trigger as
$insertThreadsVotes$
begin
    update thread set votes = thread.votes + new.voice where id = new.threadid;
    return new;
end;
$insertThreadsVotes$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION updateThreadsVotes() RETURNS trigger as
$updateThreadsVotes$
begin
    update thread set votes = thread.votes + new.voice - old.voice where id = new.threadid;
    return new;
end;
$updateThreadsVotes$ LANGUAGE plpgsql;

CREATE TRIGGER insertThreadsVotesTrigger
    AFTER INSERT
    on vote
    for each row
EXECUTE Function insertThreadsVotes();

CREATE TRIGGER updateThreadsVotesTrigger
    AFTER UPDATE
    on vote
    for each row
EXECUTE Function updateThreadsVotes();

CREATE TRIGGER insertPathTreeTrigger
    BEFORE INSERT
    on Post
    for each row
EXECUTE Function insertPathTree();

VACUUM ANALYZE;


