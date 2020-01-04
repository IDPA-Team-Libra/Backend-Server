DROP DATABASE IF EXISTS libra;

CREATE DATABASE libra;

USE libra;

CREATE TABLE USER (
       id INTEGER NOT NULL AUTO_INCREMENT,
       username varchar (255) NOT NULL,
       password varchar (255) NOT NULL,
       email varchar (255) NOT NULL,
       creationdate varchar (512) NOT NULL,
       PRIMARY KEY (id)
);

CREATE TABLE TRANSACTION (
       id INTEGER NOT NULL AUTO_INCREMENT,
       userid INTEGER NOT NULL,
       action varchar (255),
       symbol varchar (255),
       amount INTEGER,
       value Decimal (50, 10),
       current_balance Decimal(50,10),
       date date,
       processed boolean,
       PRIMARY KEY (id),
       FOREIGN KEY (userid) REFERENCES USER (id)
);

CREATE TABLE stock (
       id INTEGER AUTO_INCREMENT,
       symbol varchar (255),
       company varchar (255),
       timedata varchar(255),
       DATA longtext,
       price varchar (255),
       last_query datetime,
       PRIMARY KEY (id)
);

CREATE TABLE portfolio (
       id INTEGER AUTO_INCREMENT,
       user_id INTEGER,
       current_value Decimal (50, 10),
       total_stocks INTEGER,
       balance Decimal (50, 10),
       start_capital Decimal (50, 10),
       PRIMARY KEY (id),
       FOREIGN KEY (user_id) REFERENCES USER (id)
);

CREATE TABLE portfolio_item (
       id INTEGER AUTO_INCREMENT,
       stock_id INTEGER,
       buy_price Decimal (50, 10),
       buy_date_time datetime,
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
       FOREIGN KEY (portfolio_id) REFERENCES portfolio (id),
       FOREIGN KEY (portfolio_item_id) REFERENCES portfolio_item (id)
);

CREATE TABLE performance (
       id INTEGER,
       userid INTEGER,
       date date,
       performance decimal(10, 2),
       PRIMARY KEY (id),
       FOREIGN KEY (userid) REFERENCES USER (id)
);
INSERT INTO stock (symbol, timedata, last_query)
VALUES ('AAPL','1', NOW()),
('AAPL','5', NOW()),
('AAPL','15', NOW()),
('AAPL','30', NOW()),
('AAPL','60', NOW()),
('XOM','1', NOW()),
('XOM','5', NOW()),
('XOM','15', NOW()),
('XOM','30', NOW()),
('XOM','60', NOW()),
('MSFT','1', NOW()),
('MSFT','5', NOW()),
('MSFT','15', NOW()),
('MSFT','30', NOW()),
('MSFT','60', NOW()),
('BAC^I','1', NOW()),
('BAC^I','5', NOW()),
('BAC^I','15', NOW()),
('BAC^I','30', NOW()),
('BAC^I','60', NOW()),
('IBM','1', NOW()),
('IBM','5', NOW()),
('IBM','15', NOW()),
('IBM','30', NOW()),
('IBM','60', NOW()),
('CVX','1', NOW()),
('CVX','5', NOW()),
('CVX','15', NOW()),
('CVX','30', NOW()),
('CVX','60', NOW()),
('GE','1', NOW()),
('GE','5', NOW()),
('GE','15', NOW()),
('GE','30', NOW()),
('GE','60', NOW()),
('WMT','1', NOW()),
('WMT','5', NOW()),
('WMT','15', NOW()),
('WMT','30', NOW()),
('WMT','60', NOW()),
('T','1', NOW()),
('T','5', NOW()),
('T','15', NOW()),
('T','30', NOW()),
('T','60', NOW()),
('JNJ','1', NOW()),
('JNJ','5', NOW()),
('JNJ','15', NOW()),
('JNJ','30', NOW()),
('JNJ','60', NOW()),
('PG','1', NOW()),
('PG','5', NOW()),
('PG','15', NOW()),
('PG','30', NOW()),
('PG','60', NOW()),
('WFC','1', NOW()),
('WFC','5', NOW()),
('WFC','15', NOW()),
('WFC','30', NOW()),
('WFC','60', NOW()),
('KO','1', NOW()),
('KO','5', NOW()),
('KO','15', NOW()),
('KO','30', NOW()),
('KO','60', NOW()),
('PFE','1', NOW()),
('PFE','5', NOW()),
('PFE','15', NOW()),
('PFE','30', NOW()),
('PFE','60', NOW()),
('JPM','1', NOW()),
('JPM','5', NOW()),
('JPM','15', NOW()),
('JPM','30', NOW()),
('JPM','60', NOW()),
('GOOG','1', NOW()),
('GOOG','5', NOW()),
('GOOG','15', NOW()),
('GOOG','30', NOW()),
('GOOG','60', NOW()),
('PM','1', NOW()),
('PM','5', NOW()),
('PM','15', NOW()),
('PM','30', NOW()),
('PM','60', NOW()),
('VOD','1', NOW()),
('VOD','5', NOW()),
('VOD','15', NOW()),
('VOD','30', NOW()),
('VOD','60', NOW()),
('ORCL','1', NOW()),
('ORCL','5', NOW()),
('ORCL','15', NOW()),
('ORCL','30', NOW()),
('ORCL','60', NOW()),
('INTC','1', NOW()),
('INTC','5', NOW()),
('INTC','15', NOW()),
('INTC','30', NOW()),
('INTC','60', NOW()),
('MRK','1', NOW()),
('MRK','5', NOW()),
('MRK','15', NOW()),
('MRK','30', NOW()),
('MRK','60', NOW()),
('VZ','1', NOW()),
('VZ','5', NOW()),
('VZ','15', NOW()),
('VZ','30', NOW()),
('VZ','60', NOW()),
('BRK/A','1', NOW()),
('BRK/A','5', NOW()),
('BRK/A','15', NOW()),
('BRK/A','30', NOW()),
('BRK/A','60', NOW()),
('QCOM','1', NOW()),
('QCOM','5', NOW()),
('QCOM','15', NOW()),
('QCOM','30', NOW()),
('QCOM','60', NOW()),
('PEP','1', NOW()),
('PEP','5', NOW()),
('PEP','15', NOW()),
('PEP','30', NOW()),
('PEP','60', NOW()),
('CSCO','1', NOW()),
('CSCO','5', NOW()),
('CSCO','15', NOW()),
('CSCO','30', NOW()),
('CSCO','60', NOW()),
('AMZN','1', NOW()),
('AMZN','5', NOW()),
('AMZN','15', NOW()),
('AMZN','30', NOW()),
('AMZN','60', NOW()),
('ABT','1', NOW()),
('ABT','5', NOW()),
('ABT','15', NOW()),
('ABT','30', NOW()),
('ABT','60', NOW()),
('MCD','1', NOW()),
('MCD','5', NOW()),
('MCD','15', NOW()),
('MCD','30', NOW()),
('MCD','60', NOW()),
('SLB','1', NOW()),
('SLB','5', NOW()),
('SLB','15', NOW()),
('SLB','30', NOW()),
('SLB','60', NOW()),
('C','1', NOW()),
('C','5', NOW()),
('C','15', NOW()),
('C','30', NOW()),
('C','60', NOW()),
('BRK/B','1', NOW()),
('BRK/B','5', NOW()),
('BRK/B','15', NOW()),
('BRK/B','30', NOW()),
('BRK/B','60', NOW()),
('BAC','1', NOW()),
('BAC','5', NOW()),
('BAC','15', NOW()),
('BAC','30', NOW()),
('BAC','60', NOW()),
('RY','1', NOW()),
('RY','5', NOW()),
('RY','15', NOW()),
('RY','30', NOW()),
('RY','60', NOW()),
('HD','1', NOW()),
('HD','5', NOW()),
('HD','15', NOW()),
('HD','30', NOW()),
('HD','60', NOW()),
('DIS','1', NOW()),
('DIS','5', NOW()),
('DIS','15', NOW()),
('DIS','30', NOW()),
('DIS','60', NOW()),
('TD','1', NOW()),
('TD','5', NOW()),
('TD','15', NOW()),
('TD','30', NOW()),
('TD','60', NOW()),
('UTX','1', NOW()),
('UTX','5', NOW()),
('UTX','15', NOW()),
('UTX','30', NOW()),
('UTX','60', NOW()),
('OXY','1', NOW()),
('OXY','5', NOW()),
('OXY','15', NOW()),
('OXY','30', NOW()),
('OXY','60', NOW()),
('AXP','1', NOW()),
('AXP','5', NOW()),
('AXP','15', NOW()),
('AXP','30', NOW()),
('AXP','60', NOW()),
('KFT','1', NOW()),
('KFT','5', NOW()),
('KFT','15', NOW()),
('KFT','30', NOW()),
('KFT','60', NOW()),
('COP','1', NOW()),
('COP','5', NOW()),
('COP','15', NOW()),
('COP','30', NOW()),
('COP','60', NOW()),
('MO','1', NOW()),
('MO','5', NOW()),
('MO','15', NOW()),
('MO','30', NOW()),
('MO','60', NOW()),
('CAT','1', NOW()),
('CAT','5', NOW()),
('CAT','15', NOW()),
('CAT','30', NOW()),
('CAT','60', NOW()),
('V','1', NOW()),
('V','5', NOW()),
('V','15', NOW()),
('V','30', NOW()),
('V','60', NOW()),
('CMCSA','1', NOW()),
('CMCSA','5', NOW()),
('CMCSA','15', NOW()),
('CMCSA','30', NOW()),
('CMCSA','60', NOW()),
('MMM','1', NOW()),
('MMM','5', NOW()),
('MMM','15', NOW()),
('MMM','30', NOW()),
('MMM','60', NOW()),
('USB','1', NOW()),
('USB','5', NOW()),
('USB','15', NOW()),
('USB','30', NOW()),
('USB','60', NOW()),
('AIG','1', NOW()),
('AIG','5', NOW()),
('AIG','15', NOW()),
('AIG','30', NOW()),
('AIG','60', NOW()),
('EMC','1', NOW()),
('EMC','5', NOW()),
('EMC','15', NOW()),
('EMC','30', NOW()),
('EMC','60', NOW()),
('CVS','1', NOW()),
('CVS','5', NOW()),
('CVS','15', NOW()),
('CVS','30', NOW()),
('CVS','60', NOW()),
('BNS','1', NOW()),
('BNS','5', NOW()),
('BNS','15', NOW()),
('BNS','30', NOW()),
('BNS','60', NOW()),
('UNH','1', NOW()),
('UNH','5', NOW()),
('UNH','15', NOW()),
('UNH','30', NOW()),
('UNH','60', NOW()),
('BA','1', NOW()),
('BA','5', NOW()),
('BA','15', NOW()),
('BA','30', NOW()),
('BA','60', NOW()),
('UPS','1', NOW()),
('UPS','5', NOW()),
('UPS','15', NOW()),
('UPS','30', NOW()),
('UPS','60', NOW()),
('BMY','1', NOW()),
('BMY','5', NOW()),
('BMY','15', NOW()),
('BMY','30', NOW()),
('BMY','60', NOW()),
('AMGN','1', NOW()),
('AMGN','5', NOW()),
('AMGN','15', NOW()),
('AMGN','30', NOW()),
('AMGN','60', NOW()),
('UNP','1', NOW()),
('UNP','5', NOW()),
('UNP','15', NOW()),
('UNP','30', NOW()),
('UNP','60', NOW()),
('GS','1', NOW()),
('GS','5', NOW()),
('GS','15', NOW()),
('GS','30', NOW()),
('GS','60', NOW()),
('MA','1', NOW()),
('MA','5', NOW()),
('MA','15', NOW()),
('MA','30', NOW()),
('MA','60', NOW()),
('EBAY','1', NOW()),
('EBAY','5', NOW()),
('EBAY','15', NOW()),
('EBAY','30', NOW()),
('EBAY','60', NOW()),
('DD','1', NOW()),
('DD','5', NOW()),
('DD','15', NOW()),
('DD','30', NOW()),
('DD','60', NOW()),
('HPQ','1', NOW()),
('HPQ','5', NOW()),
('HPQ','15', NOW()),
('HPQ','30', NOW()),
('HPQ','60', NOW()),
('LLY','1', NOW()),
('LLY','5', NOW()),
('LLY','15', NOW()),
('LLY','30', NOW()),
('LLY','60', NOW()),
('CL','1', NOW()),
('CL','5', NOW()),
('CL','15', NOW()),
('CL','30', NOW()),
('CL','60', NOW()),
('SU','1', NOW()),
('SU','5', NOW()),
('SU','15', NOW()),
('SU','30', NOW()),
('SU','60', NOW()),
('SPG','1', NOW()),
('SPG','5', NOW()),
('SPG','15', NOW()),
('SPG','30', NOW()),
('SPG','60', NOW()),
('EPD','1', NOW()),
('EPD','5', NOW()),
('EPD','15', NOW()),
('EPD','30', NOW()),
('EPD','60', NOW()),
('HON','1', NOW()),
('HON','5', NOW()),
('HON','15', NOW()),
('HON','30', NOW()),
('HON','60', NOW()),
('LVS','1', NOW()),
('LVS','5', NOW()),
('LVS','15', NOW()),
('LVS','30', NOW()),
('LVS','60', NOW()),
('ESRX','1', NOW()),
('ESRX','5', NOW()),
('ESRX','15', NOW()),
('ESRX','30', NOW()),
('ESRX','60', NOW()),
('SBUX','1', NOW()),
('SBUX','5', NOW()),
('SBUX','15', NOW()),
('SBUX','30', NOW()),
('SBUX','60', NOW()),
('NKE','1', NOW()),
('NKE','5', NOW()),
('NKE','15', NOW()),
('NKE','30', NOW()),
('NKE','60', NOW()),
('ACN','1', NOW()),
('ACN','5', NOW()),
('ACN','15', NOW()),
('ACN','30', NOW()),
('ACN','60', NOW()),
('F','1', NOW()),
('F','5', NOW()),
('F','15', NOW()),
('F','30', NOW()),
('F','60', NOW()),
('MDT','1', NOW()),
('MDT','5', NOW()),
('MDT','15', NOW()),
('MDT','30', NOW()),
('MDT','60', NOW()),
('SO','1', NOW()),
('SO','5', NOW()),
('SO','15', NOW()),
('SO','30', NOW()),
('SO','60', NOW()),
('MON','1', NOW()),
('MON','5', NOW()),
('MON','15', NOW()),
('MON','30', NOW()),
('MON','60', NOW()),
('TEVA','1', NOW()),
('TEVA','5', NOW()),
('TEVA','15', NOW()),
('TEVA','30', NOW()),
('TEVA','60', NOW()),
('DOW','1', NOW()),
('DOW','5', NOW()),
('DOW','15', NOW()),
('DOW','30', NOW()),
('DOW','60', NOW()),
('GILD','1', NOW()),
('GILD','5', NOW()),
('GILD','15', NOW()),
('GILD','30', NOW()),
('GILD','60', NOW()),
('ABX','1', NOW()),
('ABX','5', NOW()),
('ABX','15', NOW()),
('ABX','30', NOW()),
('ABX','60', NOW()),
('IMO','1', NOW()),
('IMO','5', NOW()),
('IMO','15', NOW()),
('IMO','30', NOW()),
('IMO','60', NOW()),
('LOW','1', NOW()),
('LOW','5', NOW()),
('LOW','15', NOW()),
('LOW','30', NOW()),
('LOW','60', NOW()),
('DHR','1', NOW()),
('DHR','5', NOW()),
('DHR','15', NOW()),
('DHR','30', NOW()),
('DHR','60', NOW()),
('CNI','1', NOW()),
('CNI','5', NOW()),
('CNI','15', NOW()),
('CNI','30', NOW()),
('CNI','60', NOW()),
('DB','1', NOW()),
('DB','5', NOW()),
('DB','15', NOW()),
('DB','30', NOW()),
('DB','60', NOW()),
('TGT','1', NOW()),
('TGT','5', NOW()),
('TGT','15', NOW()),
('TGT','30', NOW()),
('TGT','60', NOW()),
('PCLN','1', NOW()),
('PCLN','5', NOW()),
('PCLN','15', NOW()),
('PCLN','30', NOW()),
('PCLN','60', NOW()),
('BMO','1', NOW()),
('BMO','5', NOW()),
('BMO','15', NOW()),
('BMO','30', NOW()),
('BMO','60', NOW()),
('POT','1', NOW()),
('POT','5', NOW()),
('POT','15', NOW()),
('POT','30', NOW()),
('POT','60', NOW()),
('MET','1', NOW()),
('MET','5', NOW()),
('MET','15', NOW()),
('MET','30', NOW()),
('MET','60', NOW()),
('COST','1', NOW()),
('COST','5', NOW()),
('COST','15', NOW()),
('COST','30', NOW()),
('COST','60', NOW()),
('EMR','1', NOW()),
('EMR','5', NOW()),
('EMR','15', NOW()),
('EMR','30', NOW()),
('EMR','60', NOW()),
('BIDU','1', NOW()),
('BIDU','5', NOW()),
('BIDU','15', NOW()),
('BIDU','30', NOW()),
('BIDU','60', NOW()),
('TXN','1', NOW()),
('TXN','5', NOW()),
('TXN','15', NOW()),
('TXN','30', NOW()),
('TXN','60', NOW()),
('GM','1', NOW()),
('GM','5', NOW()),
('GM','15', NOW()),
('GM','30', NOW()),
('GM','60', NOW()),
('TWX','1', NOW()),
('TWX','5', NOW()),
('TWX','15', NOW()),
('TWX','30', NOW()),
('TWX','60', NOW()),
('CNQ','1', NOW()),
('CNQ','5', NOW()),
('CNQ','15', NOW()),
('CNQ','30', NOW()),
('CNQ','60', NOW()),
('FCX','1', NOW()),
('FCX','5', NOW()),
('FCX','15', NOW()),
('FCX','30', NOW()),
('FCX','60', NOW()),
('APC','1', NOW()),
('APC','5', NOW()),
('APC','15', NOW()),
('APC','30', NOW()),
('APC','60', NOW()),
('PNC','1', NOW()),
('PNC','5', NOW()),
('PNC','15', NOW()),
('PNC','30', NOW()),
('PNC','60', NOW()),
('PX','1', NOW()),
('PX','5', NOW()),
('PX','15', NOW()),
('PX','30', NOW()),
('PX','60', NOW()),
('APA','1', NOW()),
('APA','5', NOW()),
('APA','15', NOW()),
('APA','30', NOW()),
('APA','60', NOW()),
('BP','1', NOW()),
('BP','5', NOW()),
('BP','15', NOW()),
('BP','30', NOW()),
('BP','60', NOW()),
('YUM','1', NOW()),
('YUM','5', NOW()),
('YUM','15', NOW()),
('YUM','30', NOW()),
('YUM','60', NOW()),
('DTV','1', NOW()),
('DTV','5', NOW()),
('DTV','15', NOW()),
('DTV','30', NOW()),
('DTV','60', NOW()),
('DE','1', NOW()),
('DE','5', NOW()),
('DE','15', NOW()),
('DE','30', NOW()),
('DE','60', NOW()),
('NWSA','1', NOW()),
('NWSA','5', NOW()),
('NWSA','15', NOW()),
('NWSA','30', NOW()),
('NWSA','60', NOW()),
('QQQ','1', NOW()),
('QQQ','5', NOW()),
('QQQ','15', NOW()),
('QQQ','30', NOW()),
('QQQ','60', NOW()),
('MS','1', NOW()),
('MS','5', NOW()),
('MS','15', NOW()),
('MS','30', NOW()),
('MS','60', NOW()),
('BCE','1', NOW()),
('BCE','5', NOW()),
('BCE','15', NOW()),
('BCE','30', NOW()),
('BCE','60', NOW()),
('ENB','1', NOW()),
('ENB','5', NOW()),
('ENB','15', NOW()),
('ENB','30', NOW()),
('ENB','60', NOW()),
('BIIB','1', NOW()),
('BIIB','5', NOW()),
('BIIB','15', NOW()),
('BIIB','30', NOW()),
('BIIB','60', NOW()),
('TJX','1', NOW()),
('TJX','5', NOW()),
('TJX','15', NOW()),
('TJX','30', NOW()),
('TJX','60', NOW()),
('KMB','1', NOW()),
('KMB','5', NOW()),
('KMB','15', NOW()),
('KMB','30', NOW()),
('KMB','60', NOW()),
('CELG','1', NOW()),
('CELG','5', NOW()),
('CELG','15', NOW()),
('CELG','30', NOW()),
('CELG','60', NOW()),
('BAX','1', NOW()),
('BAX','5', NOW()),
('BAX','15', NOW()),
('BAX','30', NOW()),
('BAX','60', NOW()),
('TRP','1', NOW()),
('TRP','5', NOW()),
('TRP','15', NOW()),
('TRP','30', NOW()),
('TRP','60', NOW()),
('NOV','1', NOW()),
('NOV','5', NOW()),
('NOV','15', NOW()),
('NOV','30', NOW()),
('NOV','60', NOW()),
('HAL','1', NOW()),
('HAL','5', NOW()),
('HAL','15', NOW()),
('HAL','30', NOW()),
('HAL','60', NOW()),
('D','1', NOW()),
('D','5', NOW()),
('D','15', NOW()),
('D','30', NOW()),
('D','60', NOW()),
('GG','1', NOW()),
('GG','5', NOW()),
('GG','15', NOW()),
('GG','30', NOW()),
('GG','60', NOW()),
('CM','1', NOW()),
('CM','5', NOW()),
('CM','15', NOW()),
('CM','30', NOW()),
('CM','60', NOW()),
('MCI','1', NOW()),
('MCI','5', NOW()),
('MCI','15', NOW()),
('MCI','30', NOW()),
('MCI','60', NOW()),
('WAG','1', NOW()),
('WAG','5', NOW()),
('WAG','15', NOW()),
('WAG','30', NOW()),
('WAG','60', NOW()),
('RDS/A','1', NOW()),
('RDS/A','5', NOW()),
('RDS/A','15', NOW()),
('RDS/A','30', NOW()),
('RDS/A','60', NOW()),
('DUK','1', NOW()),
('DUK','5', NOW()),
('DUK','15', NOW()),
('DUK','30', NOW()),
('DUK','60', NOW()),
('LMT','1', NOW()),
('LMT','5', NOW()),
('LMT','15', NOW()),
('LMT','30', NOW()),
('LMT','60', NOW()),
('DELL','1', NOW()),
('DELL','5', NOW()),
('DELL','15', NOW()),
('DELL','30', NOW()),
('DELL','60', NOW()),
('EOG','1', NOW()),
('EOG','5', NOW()),
('EOG','15', NOW()),
('EOG','30', NOW()),
('EOG','60', NOW()),
('AGN','1', NOW()),
('AGN','5', NOW()),
('AGN','15', NOW()),
('AGN','30', NOW()),
('AGN','60', NOW()),
('FDX','1', NOW()),
('FDX','5', NOW()),
('FDX','15', NOW()),
('FDX','30', NOW()),
('FDX','60', NOW()),
('ERIC','1', NOW()),
('ERIC','5', NOW()),
('ERIC','15', NOW()),
('ERIC','30', NOW()),
('ERIC','60', NOW()),
('BK','1', NOW()),
('BK','5', NOW()),
('BK','15', NOW()),
('BK','30', NOW()),
('BK','60', NOW()),
('SCCO','1', NOW()),
('SCCO','5', NOW()),
('SCCO','15', NOW()),
('SCCO','30', NOW()),
('SCCO','60', NOW()),
('ITW','1', NOW()),
('ITW','5', NOW()),
('ITW','15', NOW()),
('ITW','30', NOW()),
('ITW','60', NOW()),
('NEE','1', NOW()),
('NEE','5', NOW()),
('NEE','15', NOW()),
('NEE','30', NOW()),
('NEE','60', NOW()),
('DVN','1', NOW()),
('DVN','5', NOW()),
('DVN','15', NOW()),
('DVN','30', NOW()),
('DVN','60', NOW()),
('ADP','1', NOW()),
('ADP','5', NOW()),
('ADP','15', NOW()),
('ADP','30', NOW()),
('ADP','60', NOW()),
('AMT','1', NOW()),
('AMT','5', NOW()),
('AMT','15', NOW()),
('AMT','30', NOW()),
('AMT','60', NOW()),
('COV','1', NOW()),
('COV','5', NOW()),
('COV','15', NOW()),
('COV','30', NOW()),
('COV','60', NOW()),
('INFY','1', NOW()),
('INFY','5', NOW()),
('INFY','15', NOW()),
('INFY','30', NOW()),
('INFY','60', NOW()),
('ACE','1', NOW()),
('ACE','5', NOW()),
('ACE','15', NOW()),
('ACE','30', NOW()),
('ACE','60', NOW()),
('BLK','1', NOW()),
('BLK','5', NOW()),
('BLK','15', NOW()),
('BLK','30', NOW()),
('BLK','60', NOW()),
('EXC','1', NOW()),
('EXC','5', NOW()),
('EXC','15', NOW()),
('EXC','30', NOW()),
('EXC','60', NOW()),
('TYC','1', NOW()),
('TYC','5', NOW()),
('TYC','15', NOW()),
('TYC','30', NOW()),
('TYC','60', NOW()),
('BEN','1', NOW()),
('BEN','5', NOW()),
('BEN','15', NOW()),
('BEN','30', NOW()),
('BEN','60', NOW()),
('PCP','1', NOW()),
('PCP','5', NOW()),
('PCP','15', NOW()),
('PCP','30', NOW()),
('PCP','60', NOW()),
('TRV','1', NOW()),
('TRV','5', NOW()),
('TRV','15', NOW()),
('TRV','30', NOW()),
('TRV','60', NOW()),
('COF','1', NOW()),
('COF','5', NOW()),
('COF','15', NOW()),
('COF','30', NOW()),
('COF','60', NOW()),
('KMI','1', NOW()),
('KMI','5', NOW()),
('KMI','15', NOW()),
('KMI','30', NOW()),
('KMI','60', NOW()),
('PRU','1', NOW()),
('PRU','5', NOW()),
('PRU','15', NOW()),
('PRU','30', NOW()),
('PRU','60', NOW()),
('GIS','1', NOW()),
('GIS','5', NOW()),
('GIS','15', NOW()),
('GIS','30', NOW()),
('GIS','60', NOW()),
('TWC','1', NOW()),
('TWC','5', NOW()),
('TWC','15', NOW()),
('TWC','30', NOW()),
('TWC','60', NOW()),
('TRI','1', NOW()),
('TRI','5', NOW()),
('TRI','15', NOW()),
('TRI','30', NOW()),
('TRI','60', NOW()),
('CVE','1', NOW()),
('CVE','5', NOW()),
('CVE','15', NOW()),
('CVE','30', NOW()),
('CVE','60', NOW()),
('GD','1', NOW()),
('GD','5', NOW()),
('GD','15', NOW()),
('GD','30', NOW()),
('GD','60', NOW()),
('PSA','1', NOW()),
('PSA','5', NOW()),
('PSA','15', NOW()),
('PSA','30', NOW()),
('PSA','60', NOW()),
('AMX','1', NOW()),
('AMX','5', NOW()),
('AMX','15', NOW()),
('AMX','30', NOW()),
('AMX','60', NOW()),
('CTL','1', NOW()),
('CTL','5', NOW()),
('CTL','15', NOW()),
('CTL','30', NOW()),
('CTL','60', NOW()),
('VIAB','1', NOW()),
('VIAB','5', NOW()),
('VIAB','15', NOW()),
('VIAB','30', NOW()),
('VIAB','60', NOW()),
('NSC','1', NOW()),
('NSC','5', NOW()),
('NSC','15', NOW()),
('NSC','30', NOW()),
('NSC','60', NOW()),
('RAI','1', NOW()),
('RAI','5', NOW()),
('RAI','15', NOW()),
('RAI','30', NOW()),
('RAI','60', NOW()),
('CSX','1', NOW()),
('CSX','5', NOW()),
('CSX','15', NOW()),
('CSX','30', NOW()),
('CSX','60', NOW()),
('MFC','1', NOW()),
('MFC','5', NOW()),
('MFC','15', NOW()),
('MFC','30', NOW()),
('MFC','60', NOW()),
('LYB','1', NOW()),
('LYB','5', NOW()),
('LYB','15', NOW()),
('LYB','30', NOW()),
('LYB','60', NOW()),
('NEM','1', NOW()),
('NEM','5', NOW()),
('NEM','15', NOW()),
('NEM','30', NOW()),
('NEM','60', NOW()),
('ISRG','1', NOW()),
('ISRG','5', NOW()),
('ISRG','15', NOW()),
('ISRG','30', NOW()),
('ISRG','60', NOW()),
('STT','1', NOW()),
('STT','5', NOW()),
('STT','15', NOW()),
('STT','30', NOW()),
('STT','60', NOW()),
('BBT','1', NOW()),
('BBT','5', NOW()),
('BBT','15', NOW()),
('BBT','30', NOW()),
('BBT','60', NOW()),
('WLP','1', NOW()),
('WLP','5', NOW()),
('WLP','15', NOW()),
('WLP','30', NOW()),
('WLP','60', NOW()),
('JCI','1', NOW()),
('JCI','5', NOW()),
('JCI','15', NOW()),
('JCI','30', NOW()),
('JCI','60', NOW()),
('ADM','1', NOW()),
('ADM','5', NOW()),
('ADM','15', NOW()),
('ADM','30', NOW()),
('ADM','60', NOW()),
('MPV','1', NOW()),
('MPV','5', NOW()),
('MPV','15', NOW()),
('MPV','30', NOW()),
('MPV','60', NOW()),
('CTSH','1', NOW()),
('CTSH','5', NOW()),
('CTSH','15', NOW()),
('CTSH','30', NOW()),
('CTSH','60', NOW()),
('MCK','1', NOW()),
('MCK','5', NOW()),
('MCK','15', NOW()),
('MCK','30', NOW()),
('MCK','60', NOW()),
('CRM','1', NOW()),
('CRM','5', NOW()),
('CRM','15', NOW()),
('CRM','30', NOW()),
('CRM','60', NOW()),
('COH','1', NOW()),
('COH','5', NOW()),
('COH','15', NOW()),
('COH','30', NOW()),
('COH','60', NOW()),
('EP','1', NOW()),
('EP','5', NOW()),
('EP','15', NOW()),
('EP','30', NOW()),
('EP','60', NOW()),
('GLW','1', NOW()),
('GLW','5', NOW()),
('GLW','15', NOW()),
('GLW','30', NOW()),
('GLW','60', NOW()),
('CMI','1', NOW()),
('CMI','5', NOW()),
('CMI','15', NOW()),
('CMI','30', NOW()),
('CMI','60', NOW()),
('SYK','1', NOW()),
('SYK','5', NOW()),
('SYK','15', NOW()),
('SYK','30', NOW()),
('SYK','60', NOW()),
('AFL','1', NOW()),
('AFL','5', NOW()),
('AFL','15', NOW()),
('AFL','30', NOW()),
('AFL','60', NOW()),
('CBS','1', NOW()),
('CBS','5', NOW()),
('CBS','15', NOW()),
('CBS','30', NOW()),
('CBS','60', NOW()),
('BAM','1', NOW()),
('BAM','5', NOW()),
('BAM','15', NOW()),
('BAM','30', NOW()),
('BAM','60', NOW()),
('TCK','1', NOW()),
('TCK','5', NOW()),
('TCK','15', NOW()),
('TCK','30', NOW()),
('TCK','60', NOW()),
('ASML','1', NOW()),
('ASML','5', NOW()),
('ASML','15', NOW()),
('ASML','30', NOW()),
('ASML','60', NOW()),
('CB','1', NOW()),
('CB','5', NOW()),
('CB','15', NOW()),
('CB','30', NOW()),
('CB','60', NOW()),
('TMO','1', NOW()),
('TMO','5', NOW()),
('TMO','15', NOW()),
('TMO','30', NOW()),
('TMO','60', NOW()),
('SE','1', NOW()),
('SE','5', NOW()),
('SE','15', NOW()),
('SE','30', NOW()),
('SE','60', NOW()),
('FE','1', NOW()),
('FE','5', NOW()),
('FE','15', NOW()),
('FE','30', NOW()),
('FE','60', NOW()),
('WMB','1', NOW()),
('WMB','5', NOW()),
('WMB','15', NOW()),
('WMB','30', NOW()),
('WMB','60', NOW()),
('CCL','1', NOW()),
('CCL','5', NOW()),
('CCL','15', NOW()),
('CCL','30', NOW()),
('CCL','60', NOW()),
('KMP','1', NOW()),
('KMP','5', NOW()),
('KMP','15', NOW()),
('KMP','30', NOW()),
('KMP','60', NOW()),
('PSX','1', NOW()),
('PSX','5', NOW()),
('PSX','15', NOW()),
('PSX','30', NOW()),
('PSX','60', NOW()),
('MRO','1', NOW()),
('MRO','5', NOW()),
('MRO','15', NOW()),
('MRO','30', NOW()),
('MRO','60', NOW()),
('EQR','1', NOW()),
('EQR','5', NOW()),
('EQR','15', NOW()),
('EQR','30', NOW()),
('EQR','60', NOW()),
('PCG','1', NOW()),
('PCG','5', NOW()),
('PCG','15', NOW()),
('PCG','30', NOW()),
('PCG','60', NOW()),
('AEP','1', NOW()),
('AEP','5', NOW()),
('AEP','15', NOW()),
('AEP','30', NOW()),
('AEP','60', NOW()),
('MMC','1', NOW()),
('MMC','5', NOW()),
('MMC','15', NOW()),
('MMC','30', NOW()),
('MMC','60', NOW()),
('ECL','1', NOW()),
('ECL','5', NOW()),
('ECL','15', NOW()),
('ECL','30', NOW()),
('ECL','60', NOW()),
('YHOO','1', NOW()),
('YHOO','5', NOW()),
('YHOO','15', NOW()),
('YHOO','30', NOW()),
('YHOO','60', NOW()),
('SHPGY','1', NOW()),
('SHPGY','5', NOW()),
('SHPGY','15', NOW()),
('SHPGY','30', NOW()),
('SHPGY','60', NOW()),
('CMCSK','1', NOW()),
('CMCSK','5', NOW()),
('CMCSK','15', NOW()),
('CMCSK','30', NOW()),
('CMCSK','60', NOW()),
('APD','1', NOW()),
('APD','5', NOW()),
('APD','15', NOW()),
('APD','30', NOW()),
('APD','60', NOW()),
('BHI','1', NOW()),
('BHI','5', NOW()),
('BHI','15', NOW()),
('BHI','30', NOW()),
('BHI','60', NOW()),
('K','1', NOW()),
('K','5', NOW()),
('K','15', NOW()),
('K','30', NOW()),
('K','60', NOW()),
('RTN','1', NOW()),
('RTN','5', NOW()),
('RTN','15', NOW()),
('RTN','30', NOW()),
('RTN','60', NOW()),
('GSK','1', NOW()),
('GSK','5', NOW()),
('GSK','15', NOW()),
('GSK','30', NOW()),
('GSK','60', NOW()),
('DFS','1', NOW()),
('DFS','5', NOW()),
('DFS','15', NOW()),
('DFS','30', NOW()),
('DFS','60', NOW()),
('WPZ','1', NOW()),
('WPZ','5', NOW()),
('WPZ','15', NOW()),
('WPZ','30', NOW()),
('WPZ','60', NOW()),
('SDRL','1', NOW()),
('SDRL','5', NOW()),
('SDRL','15', NOW()),
('SDRL','30', NOW()),
('SDRL','60', NOW()),
('HES','1', NOW()),
('HES','5', NOW()),
('HES','15', NOW()),
('HES','30', NOW()),
('HES','60', NOW());



CREATE USER 'administrator' @'localhost' IDENTIFIED BY 'LOCAL1234';

GRANT ALL PRIVILEGES ON libra.* TO 'administrator' @'localhost';
