CREATE TABLE User(
    id INTEGER NOT NULL AUTO_INCREMENT,
    username varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    creationdate varchar(512) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE UserAction(
    id INTEGER NOT NULL AUTO_INCREMENT,
    userID INTEGER NOT NULL,
    activity varchar(255),
    accurance datetime,
    status varchar(255),
    PRIMARY KEY (id),
    FOREIGN KEY (userID) REFERENCES User(id)
);

CREATE TABLE Transaction(


);