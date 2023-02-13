-- Creating a new User.

CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT ON snippetbox.* TO 'web'@'localhost';

-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';