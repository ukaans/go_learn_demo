create table access
(
    id          int auto_increment
        primary key,
    module_name varchar(255) default '0' null,
    type        tinyint(1)               null,
    action_name varchar(255)             null,
    url         varchar(255)             null,
    module_id   int                      null,
    sort        int                      null,
    description varchar(255)             null,
    add_time    int                      null,
    status      int                      null
)
    charset = utf8mb3
    row_format = COMPACT;

create table address
(
    id              int auto_increment
        primary key,
    uid             int          null,
    name            varchar(255) null,
    phone           varchar(11)  null,
    address         varchar(255) null,
    default_address tinyint(1)   null,
    add_time        int          null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table focus
(
    id         int auto_increment
        primary key,
    title      varchar(255) null,
    focus_type tinyint(1)   null,
    focus_img  varchar(255) null,
    link       varchar(255) null,
    sort       int          null,
    status     tinyint(1)   null,
    add_time   int          null
)
    charset = utf8mb3
    row_format = COMPACT;

create table goods
(
    id             int auto_increment
        primary key,
    title          varchar(255)   null,
    sub_title      varchar(255)   null,
    goods_sn       varchar(255)   null,
    cate_id        int            null,
    click_count    int            null,
    goods_number   int            null,
    price          decimal(10, 2) null,
    market_price   decimal(10, 2) null,
    relation_goods varchar(255)   null,
    goods_attr     varchar(1024)  null,
    goods_color    varchar(255)   null,
    goods_version  varchar(255)   null,
    goods_img      varchar(255)   null,
    goods_gift     varchar(255)   null,
    goods_fitting  varchar(255)   null,
    goods_keywords varchar(255)   null,
    goods_desc     varchar(255)   null,
    goods_content  text           null,
    is_delete      tinyint        null,
    is_hot         tinyint        null,
    is_best        tinyint        null,
    is_new         tinyint        null,
    goods_type_id  int            null,
    sort           int            null,
    status         tinyint        null,
    add_time       int            null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table goods_attr
(
    id                int auto_increment
        primary key,
    goods_id          int          null,
    attribute_cate_id int          null,
    attribute_id      int          null,
    attribute_title   varchar(255) null,
    attribute_type    tinyint(1)   null,
    attribute_value   varchar(255) null,
    sort              int          null,
    add_time          int          null,
    status            tinyint(1)   null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table goods_cate
(
    id          int auto_increment
        primary key,
    title       varchar(255)  null,
    cate_img    varchar(255)  null,
    link        varchar(255)  null,
    template    varchar(255)  null,
    pid         int           null,
    sub_title   varchar(255)  null,
    keywords    varchar(255)  null,
    description varchar(1024) null,
    status      tinyint(1)    null,
    sort        varchar(255)  null,
    add_time    int           null
)
    charset = utf8mb3
    row_format = DYNAMIC;

create table goods_color
(
    id          int auto_increment
        primary key,
    color_name  varchar(255) null,
    color_value varchar(255) null,
    status      int          null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table goods_image
(
    id       int auto_increment
        primary key,
    goods_id int          null,
    img_url  varchar(255) null,
    color_id int          null,
    sort     int          null,
    add_time int          null,
    status   tinyint(1)   null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table goods_type
(
    id          int auto_increment
        primary key,
    title       varchar(255) null,
    description varchar(255) null,
    status      int          null,
    add_time    int          null
)
    charset = utf8mb3
    row_format = DYNAMIC;

create table goods_type_attribute
(
    id         int auto_increment
        primary key,
    cate_id    int           null,
    title      varchar(255)  null,
    attr_type  varchar(255)  null,
    attr_value varchar(1024) null,
    status     tinyint(1)    null,
    sort       int           null,
    add_time   int           null
)
    charset = utf8mb3
    row_format = DYNAMIC;

create table manager
(
    id       int auto_increment
        primary key,
    username varchar(255)         null,
    password varchar(32)          null,
    mobile   varchar(11)          null,
    email    varchar(255)         null,
    status   tinyint(1)           null,
    role_id  int                  null,
    add_time int                  null,
    is_super tinyint(1) default 0 null
)
    charset = utf8mb3
    row_format = COMPACT;

create index role_id
    on manager (role_id);

create table nav
(
    id         int auto_increment
        primary key,
    title      varchar(255) null,
    link       varchar(255) null,
    position   tinyint(1)   null,
    is_opennew tinyint(1)   null,
    relation   varchar(255) null,
    sort       int          null,
    status     tinyint(1)   null,
    add_time   int          null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table `order`
(
    id                int auto_increment
        primary key,
    uid               int            null,
    order_id          varchar(255)   null,
    all_price         decimal(10, 2) null,
    name              varchar(255)   null,
    phone             varchar(11)    null,
    address           varchar(255)   null,
    zipcode           varchar(255)   null,
    pay_status        tinyint(1)     null,
    pay_type          tinyint(1)     null,
    order_status      tinyint(1)     null,
    add_time          int            null,
    pay_time          int            null,
    distribution_time int default 0  not null,
    cancel_time       int            null,
    logistics_company varchar(255)   null,
    waybill_no        varchar(255)   null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table order_item
(
    id            int auto_increment
        primary key,
    order_id      int            null,
    uid           int            null,
    product_title varchar(255)   null,
    product_id    int            null,
    product_img   varchar(255)   null,
    product_price decimal(10, 2) null,
    product_num   int            null,
    goods_version varchar(255)   null,
    goods_color   varchar(255)   null,
    add_time      int            null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table role
(
    id          int auto_increment
        primary key,
    title       varchar(255) null,
    description varchar(255) null,
    status      tinyint(1)   null,
    add_time    int          null
)
    row_format = COMPACT;

create table role_access
(
    role_id   int not null,
    access_id int not null
)
    charset = utf8mb3
    row_format = COMPACT;

create index access_id
    on role_access (access_id);

create index role_id
    on role_access (role_id);

create table setting
(
    id               int auto_increment
        primary key,
    site_title       varchar(255) null,
    site_logo        varchar(255) null,
    site_keywords    varchar(255) null,
    site_description varchar(255) null,
    no_picture       varchar(255) null,
    site_icp         varchar(255) null,
    site_tel         varchar(255) null,
    search_keywords  varchar(255) null,
    tongji_code      varchar(255) null,
    appid            varchar(255) null,
    app_secret       varchar(255) null,
    end_point        varchar(255) null,
    bucket_name      varchar(255) null,
    oss_status       tinyint(1)   null,
    oss_domain       varchar(255) null,
    thumbnail_size   varchar(255) null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table user
(
    id       int auto_increment
        primary key,
    password varchar(32)  null,
    phone    varchar(11)  null,
    last_ip  varchar(255) null,
    email    varchar(255) null,
    add_time int          null,
    status   tinyint      null
)
    charset = utf8mb4
    row_format = DYNAMIC;

create table user_temp
(
    id         int auto_increment
        primary key,
    ip         varchar(255) null,
    phone      varchar(11)  null,
    send_count int          null,
    add_day    int          null,
    add_time   int          null,
    sign       varchar(32)  null
)
    charset = utf8mb4
    row_format = DYNAMIC;


