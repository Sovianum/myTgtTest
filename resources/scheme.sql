DROP INDEX IF EXISTS stats_ts_idx ;
DROP INDEX IF EXISTS stats_userid_idx ;
DROP TABLE IF EXISTS Client CASCADE ;
DROP TABLE IF EXISTS Stats CASCADE;
DROP TYPE IF EXISTS Esex CASCADE ;
DROP TYPE IF EXISTS EAction CASCADE ;

CREATE TYPE ESex AS ENUM ('M', 'F');
CREATE TYPE EAction AS ENUM ('login', 'like', 'comments', 'exit');

CREATE TABLE Client (
  id INTEGER PRIMARY KEY ,
  age INTEGER,
  sex ESex
);

CREATE TABLE Stats (
  id SERIAL PRIMARY KEY ,
  userId INTEGER REFERENCES Client(id),
  ts TIMESTAMP,
  action EAction,
  counter INT DEFAULT 1,
  UNIQUE (userId, ts, action)
);

CREATE INDEX stats_userid_idx ON Stats (userId);
CREATE INDEX stats_ts_userid_idx ON Stats (userId, ts);