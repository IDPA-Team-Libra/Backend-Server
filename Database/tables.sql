CREATE TABLE User
(
        id INTEGER NOT NULL
        AUTO_INCREMENT,
    username varchar
        (255) NOT NULL,
    password varchar
        (255) NOT NULL,
    email varchar
        (255) NOT NULL,
    creationdate varchar
        (512) NOT NULL,
    PRIMARY KEY
        (id)
);

        CREATE TABLE UserAction
        (
                id INTEGER NOT NULL
                AUTO_INCREMENT,
    userID INTEGER NOT NULL,
    activity varchar
                (255),
    accurance datetime,
    status varchar
                (255),
    PRIMARY KEY
                (id),
    FOREIGN KEY
                (userID) REFERENCES User
                (id)
);


                CREATE TABLE Transaction
                (
                        id INTEGER NOT NULL
                        AUTO_INCREMENT,
    userID INTEGER NOT NULL,
    action varchar
                        (255),
    description varchar
                        (255),
    amount INTEGER,
    value Decimal
                        (10,4),
    date date,
    PRIMARY KEY
                        (id),
    FOREIGN KEY
                        (userID) REFERENCES User
                        (id)
);




                        CREATE TABLE Stock
                        (
                                id INTEGER
                                AUTO_INCREMENT,
    symbol varchar
                                (255),
    company varchar
                                (255),
    timeData ENUM
                                ('Daily','5','15','30','60'),
    data LONGTEXT,
    price varchar
                                (255),
    last_query DATETIME,
    PRIMARY KEY
                                (id)
);
                                CREATE TABLE Portfolio
                                (
                                        id INTEGER
                                        AUTO_INCREMENT,
                                user_id INTEGER,
                                current_value Decimal
                                        (15,4),
                                total_stocks INTEGER,
                                start_capital Decimal
                                        (15,4),
                                PRIMARY KEY
                                        (id),
                                FOREIGN KEY
                                        (user_id) REFERENCES User
                                        (id)
                        );

                                        CREATE TABLE portfolio_item
                                        (
                                                id INTEGER
                                                AUTO_INCREMENT,
                                stock_id INTEGER,
                                buy_price Decimal
                                                (10,4),
                                buy_date_time DATETIME,
                                quantity INTEGER,
                                total_buy_price Decimal
                                                (10,4),
                                );