CREATE TABLE Users
(
    id       serial       NOT NULL,
    nickname varchar(255) NOT NULL UNIQUE,
    fullname varchar(255) NOT NULL UNIQUE,
    about    varchar(255) NOT NULL,
    email    varchar(255) NOT NULL UNIQUE,
    CONSTRAINT UserPk PRIMARY KEY (id)
);

CREATE TABLE Forum
(
    id      serial       NOT NULL,
    title   varchar(255) NOT NULL UNIQUE,
    user_nickname varchar(255)          NOT NULL,
    slug    varchar(255) NOT NULL,
    posts   int          NOT NULL default 0,
    threads int          NOT NULL default 0,
    CONSTRAINT ForumPk PRIMARY KEY (id)
);

CREATE TABLE Thread
(
    id serial NOT NULL,
    title varchar(255) NOT NULL UNIQUE,
    user_nickname varchar(255) NOT NULL,
    forum_title varchar(255) NOT NULL,
    message varchar(255) NOT NULL,
    votes int default 0,
    slug varchar(255) NOT NULL,
    created date default 'now',
    CONSTRAINT ThreadPk PRIMARY KEY (id)
);

ALTER TABLE Forum
    ADD CONSTRAINT realUserAtForum FOREIGN KEY (user_nickname) REFERENCES Users (nickname) ON DELETE CASCADE;

ALTER TABLE Thread
    ADD CONSTRAINT realForumInThread FOREIGN KEY (forum_title) REFERENCES Forum (title) ON DELETE CASCADE;

ALTER TABLE Thread
    ADD CONSTRAINT realSlugInThread FOREIGN KEY (slug) REFERENCES Forum (slug) ON DELETE CASCADE;

ALTER TABLE Thread
    ADD CONSTRAINT realUserInThread FOREIGN KEY (user_nickname) REFERENCES Users (nickname) ON DELETE CASCADE;


