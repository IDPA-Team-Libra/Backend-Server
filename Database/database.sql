DROP DATABASE IF EXISTS libra;
CREATE DATABASE libra;
use libra;


CREATE TABLE User (
        id INTEGER NOT NULL AUTO_INCREMENT,
        username varchar (255) NOT NULL,
        password varchar (255) NOT NULL,
        email varchar (255) NOT NULL,
        creationdate varchar (512) NOT NULL,
        PRIMARY KEY (id)
);

CREATE TABLE UserAction (
        id INTEGER NOT NULL AUTO_INCREMENT,
        userID INTEGER NOT NULL,
        activity varchar (255),
        accurance datetime,
        status varchar (255),
        PRIMARY KEY (id),
        FOREIGN KEY (userID) REFERENCES User (id)
);

CREATE TABLE transaction (
        id INTEGER NOT NULL AUTO_INCREMENT,
        userID INTEGER NOT NULL,
        action varchar (255),
        description varchar (255),
        amount INTEGER,
        value Decimal (50, 10),
        date date,
        processed boolean,
        PRIMARY KEY (id),
        FOREIGN KEY (userID) REFERENCES User (id)
);

CREATE TABLE stock (
        id INTEGER AUTO_INCREMENT,
        symbol varchar (255),
        company varchar (255),
        timeData ENUM ('Daily', '5', '15', '30', '60'),
        data LONGTEXT,
        price varchar (255),
        last_query DATETIME,
        PRIMARY KEY (id)
);

CREATE TABLE Portfolio (
        id INTEGER AUTO_INCREMENT,
        user_id INTEGER,
        current_value Decimal (50, 10),
        total_stocks INTEGER,
        balance Decimal (50, 10),
        start_capital Decimal (50, 10),
        PRIMARY KEY (id),
        FOREIGN KEY (user_id) REFERENCES User (id)
);

CREATE TABLE portfolio_item (
        id INTEGER AUTO_INCREMENT,
        stock_id INTEGER,
        buy_price Decimal (50,10),
        buy_date_time DATETIME,
        quantity INTEGER,
        total_buy_price Decimal (50, 10),
        PRIMARY KEY (id),
        FOREIGN KEY (stock_id) REFERENCES stock (id)
);

CREATE TABLE portfolio_to_item (
        id INTEGER AUTO_INCREMENT,
        portfolio_id INTEGER,
        portfolio_item_id INTEGER,
        PRIMARY KEY (id),
        FOREIGN KEY (portfolio_id) REFERENCES Portfolio (id),
        FOREIGN KEY (portfolio_item_id) REFERENCES portfolio_item (id)
);

CREATE TABLE performance (
        id INTEGER,
        userID INTEGER,
        date date,
        performance DECIMAL(10, 2),
        PRIMARY KEY (id),
        FOREIGN KEY (userID) REFERENCES User (id)
);


INSERT INTO stock
    (symbol,company,timeData,last_query)
VALUES
    ('AMZN', 'Amazone', '5', NOW()),
    ('TSLA', 'Tesla', '5', NOW()),
    ('YHOO', 'Yahoo', '5', NOW()),
    ('CELG', 'Celgene', '5', NOW()),
    ('FCEL', 'Fuelcell Energy', '5', NOW()),
    ('AAPL', 'Apple', '5', NOW()),
    ('INTC', 'Intel', '5', NOW()),
    ('MRVL', 'Marvel', '5', NOW()),
    ('ICFI', 'Ico', '5', NOW()),
    ('IGOI', 'IGO', '5', NOW()),
    ('IMAX', 'Imax', '5', NOW()),
    ('IRIS', 'Iris', '5', NOW()),
    ('ITRI', 'Itron', '5', NOW()),
    ('IXYS', 'IXYS', '5', NOW()),
    ('IMAX', 'Imax', '5', NOW()),
    ('AAPL', 'Apple Inc.', '5', NOW()),
    ('XOM', 'Exxon Mobil Corporation', '5', NOW()),
    ('MSFT', 'Microsoft Corporation', '5', NOW()),
    ('BAC^I', 'Bank of America Corporation', '5', NOW()),
    ('IBM', 'International Business Machines Corporation', '5', NOW()),
    ('CVX', 'Chevron Corporation', '5', NOW()),
    ('GE', 'General Electric Company', '5', NOW()),
    ('WMT', 'Wal-Mart Stores', '5', NOW()),
    ('T', ' Inc.', '5', NOW()),
    ('JNJ', 'AT&T Inc.', '5', NOW()),
    ('PG', 'Johnson & Johnson', '5', NOW()),
    ('WFC', 'Procter & Gamble Company (The)', '5', NOW()),
    ('KO', 'Wells Fargo & Company', '5', NOW()),
    ('PFE', 'Coca-Cola Company (The)', '5', NOW()),
    ('JPM', 'Pfizer', '5', NOW()),
    ('GOOG', ' Inc.', '5', NOW()),
    ('PM', 'J P Morgan Chase & Co', '5', NOW()),
    ('VOD', 'Google Inc.', '5', NOW()),
    ('ORCL', 'Philip Morris International Inc', '5', NOW()),
    ('INTC', 'Vodafone Group Plc', '5', NOW()),
    ('MRK', 'Oracle Corporation', '5', NOW()),
    ('VZ', 'Intel Corporation', '5', NOW()),
    ('BRK/A', 'Merck & Company', '5', NOW()),
    ('QCOM', ' Inc.', '5', NOW()),
    ('PEP', 'Verizon Communications Inc.', '5', NOW()),
    ('CSCO', 'Berkshire Hathaway Inc.', '5', NOW()),
    ('AMZN', 'QUALCOMM Incorporated', '5', NOW()),
    ('ABT', 'Pepsico', '5', NOW()),
    ('MCD', ' Inc.', '5', NOW()),
    ('SLB', 'Cisco Systems', '5', NOW()),
    ('C', ' Inc.', '5', NOW()),
    ('BRK/B', 'Amazon.com', '5', NOW()),
    ('BAC', ' Inc.', '5', NOW()),
    ('RY', 'Abbott Laboratories', '5', NOW()),
    ('HD', 'McDonald&#39;s Corporation', '5', NOW()),
    ('DIS', 'Schlumberger N.V.', '5', NOW()),
    ('TD', 'Citigroup Inc.', '5', NOW()),
    ('UTX', 'Berkshire Hathaway Inc.', '5', NOW()),
    ('OXY', 'Bank of America Corporation', '5', NOW()),
    ('AXP', 'Royal Bank Of Canada', '5', NOW()),
    ('KFT', 'Home Depot', '5', NOW()),
    ('COP', ' Inc. (The)', '5', NOW()),
    ('MO', 'Walt Disney Company (The)', '5', NOW()),
    ('CAT', 'Toronto Dominion Bank (The)', '5', NOW()),
    ('V', 'United Technologies Corporation', '5', NOW()),
    ('CMCSA', 'Occidental Petroleum Corporation', '5', NOW()),
    ('MMM', 'American Express Company', '5', NOW()),
    ('USB', 'Kraft Foods Inc.', '5', NOW()),
    ('AIG', 'ConocoPhillips', '5', NOW()),
    ('EMC', 'Altria Group', '5', NOW()),
    ('CVS', 'Caterpillar', '5', NOW()),
    ('BNS', ' Inc.', '5', NOW()),
    ('UNH', 'Visa Inc.', '5', NOW()),
    ('BA', 'Comcast Corporation', '5', NOW()),
    ('UPS', '3M Company', '5', NOW()),
    ('BMY', 'U.S. Bancorp', '5', NOW()),
    ('AMGN', 'American International Group', '5', NOW()),
    ('UNP', ' Inc.', '5', NOW()),
    ('GS', 'EMC Corporation', '5', NOW()),
    ('MA', 'CVS Corporation', '5', NOW()),
    ('EBAY', 'Bank of Nova Scotia (The)', '5', NOW()),
    ('DD', 'UnitedHealth Group Incorporated', '5', NOW()),
    ('HPQ', 'Boeing Company (The)', '5', NOW()),
    ('LLY', 'United Parcel Service', '5', NOW()),
    ('CL', ' Inc.', '5', NOW()),
    ('SU', 'Bristol-Myers Squibb Company', '5', NOW()),
    ('SPG', 'Amgen Inc.', '5', NOW());

CREATE USER 'administrator'@'localhost' IDENTIFIED BY 'LOCAL1234';
GRANT ALL PRIVILEGES ON libra.* TO 'administrator'@'localhost';
