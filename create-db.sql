-- Table: ticks

DROP TABLE ticks;

CREATE TABLE ticks
(
  stock  character varying(10) NOT NULL,
  date   date                  NOT NULL,
  open   numeric(8,3)          NOT NULL,
  low    numeric(8,3)          NOT NULL,
  high   numeric(8,3)          NOT NULL,
  close  numeric(8,3)          NOT NULL,
  volume bigint                NOT NULL,
  CONSTRAINT ticks_pk PRIMARY KEY (stock, date)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE ticks
  OWNER TO ubuntu;

-- Index: ticks_date

DROP INDEX ticks_date;

CREATE INDEX ticks_date
  ON ticks
  USING btree
  (date);

-- Index: ticks_stock

DROP INDEX ticks_stock;

CREATE INDEX ticks_stock
  ON ticks
  USING btree
  (stock COLLATE pg_catalog."default");

