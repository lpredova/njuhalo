CREATE TABLE queries (
  id integer PRIMARY KEY AUTOINCREMENT,
  name text,
  isActive integer,
  url text,
  filters text,
  createdAt integer
);

CREATE TABLE items (
  id integer PRIMARY KEY AUTOINCREMENT,
  queryID integer,
  itemID integer,
  isNew integer,
  url text,
  name text,
  image text,
  price text,
  description text,
  location text,
  year text,
  mileage text,
  published text,
  createdAt integer,
  FOREIGN KEY (queryID) REFERENCES queries(id) ON DELETE CASCADE
);

CREATE INDEX index_itemID ON items (itemID);
