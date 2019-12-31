
create type GenDateVal AS (
    modifier int,
    calendar int,
    year int,
    month int,
    day int
);

create type GenDate AS (
    type text,
    start_date GenDateVal,
    end_date GenDateVal,
    sort_key bigint
);

create table note
(
    id text primary key,
    detail text,
    date GenDate
);

create table holder
(
    id bigint primary key,
    name text,
    address1 text,
    address2 text,
    address3 text,
    address4 text,
    address5 text,
    phone text,
    email text,
    url text,
    abbreviation text,
    detail text
);

create table source
(
    id text primary key,
    detail text,
    author text,
    title text,
    location text,
    holder bigint references holder(id),
    date GenDate
);



create table person
(
    id text primary key,
    userid text,
    birthsex text,
    isPrivate boolean,
    noteid text references note(id),
    sourceid text references source(id)
);

create table family
(
    id text primary key,
    userid text,
    primeid text references person(id),
    partnerid text references person(id),
    noteid text references note(id),
    sourceid text references source(id)
);

create table child
(
    id text primary key,
    personid text references person(id),
    familyid text references family(id),
    parent1id text references person(id),
    parent1rel text,
    parent2id text references person(id),
    parent2rel text,
    ordinal int,
    noteid text references note(id),
    sourceid text references source(id)
);

create table fact
(
    id text primary key,
    type text,
    referenceid text,
    detail text,
    surety int,
    date GenDate,
    place text,
    noteid text references note(id),
    sourceid text references source(id)
);

create table attach
(
    id text primary key,
    referenceid text,
    filename text,
    fileinfo text,
    detail text,
    sourceid text references source(id)
);

create table name
(
    id text primary key,
    type text,
    isPrefered boolean,
    personid text references person(id),
    given text,
    surname text,
    familiar text,
    title text,
    suffix text,
    surety int,
    date GenDate,
    displayAs text,
    noteid text references note(id),
    sourceid text references source(id)
);

create table header
(
    version text,
    dbinfo text,
    copyright text,
    comment text,
    peoplecount bigint,
    familycount bigint,
    name text,
    address1 text,
    address2 text,
    address3 text,
    address4 text,
    address5 text,
    phone text,
    email text,
    url text,
    date GenDate
);
