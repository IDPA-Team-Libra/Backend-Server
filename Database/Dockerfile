FROM mysql:5.6
COPY mysql_install.sh /mysql/setup.sh
COPY database.sql /mysql/setup.sql
RUN ["chmod", "+x", "/mysql/setup.sh"]
RUN /mysql/setup.sh
