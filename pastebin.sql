use pastebin;
create table log
(
    id  int auto_increment
        primary key,
    log json null
);


create table paste
(
    id       int auto_increment
        primary key,
    pwd      varchar(300) null,
    times    int          null,
    deadline datetime     null,
    paste    json         null
);