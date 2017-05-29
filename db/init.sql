CREATE TABLE items (
  id integer PRIMARY KEY AUTOINCREMENT,
 itemID integer ,
 url text NOT NULL,
 name text NOT NULL,
 image text NOT NULL UNIQUE,
 price text NOT NULL UNIQUE,
 description text NOT NULL UNIQUE
);
