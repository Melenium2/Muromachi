create table if not exists app_tracking
(
    id          serial primary key not null,
    bundle      varchar(255) not null,
    category    varchar(128) not null,
    developerId varchar(128),
    developer   varchar(500),
    Geo         varchar(10) not null,
    startAt     timestamp,
    period      int
);
create table if not exists category_tracking
(
    id       serial primary key not null,
    bundleId int references app_tracking,
    type     varchar(128) not null,
    place    int not null,
    date     timestamp not null
);
create table if not exists keyword_tracking
(
    id       serial primary key not null,
    bundleId int references app_tracking,
    type     varchar(128) not null,
    place    int not null,
    date     timestamp not null
);
create
type developerContacts as
(
    email    text,
    contacts text
);
create table if not exists meta_tracking
(
    id               bigserial primary key not null,
    bundleId         int references app_tracking,
    title            varchar(300),
    price            varchar(50),
    picture          text,
    screenshots      text[],
    rating           varchar(50),
    reviewCount      varchar(50),
    ratingHistogram  varchar(50)[],
    description      text,
    shortDescription text,
    recentChanges    text,
    releaseDate      varchar(50),
    lastUpdateDate   varchar(50),
    appSize          varchar(50),
    installs         varchar(50),
    version          varchar(100),
    androidVersion   varchar(100),
    contentRating    varchar(100),
    devContacts      developerContacts,
    privacyPolicy    text,
    date             timestamp
);
create table if not exists users
(
    id           bigserial primary key not null,
    clientId     varchar(50) not null,
    clientSecret varchar(70) not null,
    company      varchar(250) not null,
    addedAt      timestamp with time zone NOT NULL DEFAULT now()
);
create table if not exists refresh_sessions
(
    id           bigserial primary key not null,
    userId       int REFERENCES users (id) ON DELETE CASCADE,
    refreshToken text not null,
    useragent    text,
    ip           character varying(15) NOT NULL,
    expiresIn    timestamp with time zone NOT NULL NOT NULL,
    createdAt    timestamp with time zone NOT NULL DEFAULT now()
);
