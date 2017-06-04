CREATE TABLE items (
  id integer PRIMARY KEY AUTOINCREMENT,
  itemID integer,
  url text,
  name text,
  image text,
  price text,
  description text,
  createdAt integer
);


CREATE INDEX index_itemID ON items (itemID);
