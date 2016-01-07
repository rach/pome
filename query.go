package main

import (
	"fmt"
	"github.com/rach/pome/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	_ "github.com/rach/pome/Godeps/_workspace/src/github.com/lib/pq"
	"log"
)

func connectDB(dbURL string) *sqlx.DB {
	db := sqlx.MustOpen("postgres", dbURL)
	return db
}

func connectionString(host string, dbname string, username string, password string, sslmode string) string {
	// TODO: escape single quote
	return fmt.Sprintf("host='%s' dbname='%s' user='%s' password='%s' sslmode='%s'", host, dbname, username, password, sslmode)
}

type IndexBloatDatabaseResult struct {
	Key        string  `db:"key"`
	Schema     string  `db:"schema"`
	Table      string  `db:"table"`
	Index      string  `db:"index"`
	BloatBytes int64   `db:"bloat_bytes"`
	BloatRatio float64 `db:"bloat_ratio"`
}

type NumberOfConnectionResult struct {
	Count int `db:"num_connections"`
}
type DatabaseSizeResult struct {
	TableSize  int     `db:"table_size"`
	IndexSize  int     `db:"index_size"`
	TotalSize  int     `db:"total_size"`
	IndexRatio float64 `db:"index_ratio"`
}

type TableBloatDatabaseResult struct {
	Key        string  `db:"key"`
	Schema     string  `db:"schema"`
	Table      string  `db:"table"`
	BloatBytes int64   `db:"bloat_bytes"`
	BloatRatio float64 `db:"bloat_ratio"`
}

type databaseResultFct func(db *sqlx.DB) interface{}

func GetIndexBloatResult(db *sqlx.DB) interface{} {
	if *verbose {
		log.Printf("Fetch Indexes Bloat")
	}
	r := []IndexBloatDatabaseResult{}
	err := db.Select(&r, IndexBloatSql)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func GetTableBloatResult(db *sqlx.DB) interface{} {
	if *verbose {
		log.Printf("Fetch Tables Bloat")
	}
	r := []TableBloatDatabaseResult{}
	err := db.Select(&r, TableBloatSql)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func GetNumberOfConnectionResult(db *sqlx.DB) interface{} {
	if *verbose {
		log.Printf("Fetch Number of Connections")
	}
	r := NumberOfConnectionResult{}
	err := db.Get(&r, NumberOfConnectionSql)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func GetDatabeSizeResult(db *sqlx.DB) interface{} {
	if *verbose {
		log.Printf("Fetch Database Size")
	}
	r := DatabaseSizeResult{}
	err := db.Get(&r, DatabaseSizeSql)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

// Modified Query from:
// https://gist.github.com/gullevek/32881d6b4c5b1ed0135c
const (
	IndexBloatSql = `
WITH btree_index_atts AS (
    SELECT nspname, relname, reltuples, relpages, indrelid, relam,
        regexp_split_to_table(indkey::text, ' ')::smallint AS attnum,
        indexrelid as index_oid
    FROM pg_index
    JOIN pg_class ON pg_class.oid=pg_index.indexrelid
    JOIN pg_namespace ON pg_namespace.oid = pg_class.relnamespace
    JOIN pg_am ON pg_class.relam = pg_am.oid
    WHERE pg_am.amname = 'btree'
    ),
index_item_sizes AS (
    SELECT
    i.nspname, i.relname, i.reltuples, i.relpages, i.relam,
    (quote_ident(s.schemaname) || '.' || quote_ident(s.tablename))::regclass AS starelid, a.attrelid AS table_oid, index_oid,
    current_setting('block_size')::numeric AS bs,
    /* MAXALIGN: 4 on 32bits, 8 on 64bits (and mingw32 ?) */
    CASE
        WHEN version() ~ 'mingw32' OR version() ~ '64-bit' THEN 8
        ELSE 4
    END AS maxalign,
    24 AS pagehdr,
    /* per tuple header: add index_attribute_bm if some cols are null-able */
    CASE WHEN max(coalesce(s.null_frac,0)) = 0
        THEN 2
        ELSE 6
    END AS index_tuple_hdr,
    /* data len: we remove null values save space using it fractionnal part from stats */
    sum( (1-coalesce(s.null_frac, 0)) * coalesce(s.avg_width, 2048) ) AS nulldatawidth
    FROM pg_attribute AS a
    JOIN pg_stats AS s ON (quote_ident(s.schemaname) || '.' || quote_ident(s.tablename))::regclass=a.attrelid AND s.attname = a.attname 
    JOIN btree_index_atts AS i ON i.indrelid = a.attrelid AND a.attnum = i.attnum
    WHERE a.attnum > 0
    GROUP BY 1, 2, 3, 4, 5, 6, 7, 8, 9
),
index_aligned AS (
    SELECT maxalign, bs, nspname, relname AS index_name, reltuples,
        relpages, relam, table_oid, index_oid,
      ( 2 +
          maxalign - CASE /* Add padding to the index tuple header to align on MAXALIGN */
            WHEN index_tuple_hdr%maxalign = 0 THEN maxalign
            ELSE index_tuple_hdr%maxalign
          END
        + nulldatawidth + maxalign - CASE /* Add padding to the data to align on MAXALIGN */
            WHEN nulldatawidth::integer%maxalign = 0 THEN maxalign
            ELSE nulldatawidth::integer%maxalign
          END
      )::numeric AS nulldatahdrwidth, pagehdr
    FROM index_item_sizes AS s1
),
otta_calc AS (
  SELECT bs, nspname, table_oid, index_oid, index_name, relpages, coalesce(
    ceil((reltuples*(4+nulldatahdrwidth))/(bs-pagehdr::float)) +
      CASE WHEN am.amname IN ('hash','btree') THEN 1 ELSE 0 END , 0 -- btree and hash have a metadata reserved block
    ) AS otta
  FROM index_aligned AS s2
    LEFT JOIN pg_am am ON s2.relam = am.oid
),
raw_bloat AS (
    SELECT current_database() as dbname, nspname, c.relname AS tablename, index_name,
        bs*(sub.relpages)::bigint AS totalbytes,
        CASE
            WHEN sub.relpages <= otta THEN 0
            ELSE bs*(sub.relpages-otta)::bigint END
            AS wastedbytes,
        CASE
            WHEN sub.relpages <= otta
            THEN 0 ELSE bs*(sub.relpages-otta)::bigint * 100 / (bs*(sub.relpages)::bigint) END
            AS realbloat,
        pg_relation_size(sub.table_oid) as table_bytes
    FROM otta_calc AS sub
    JOIN pg_class AS c ON c.oid=sub.table_oid
)
SELECT  nspname || '.' || tablename || '.' || index_name AS key,
				nspname AS schema,
				tablename as table,
        index_name AS index,
        wastedbytes as bloat_bytes,
        round(realbloat, 1) as bloat_ratio
  --     , totalbytes as index_size,
  --      table_bytes, pg_size_pretty(table_bytes) as table_size
FROM raw_bloat
ORDER BY wastedbytes DESC;
`

	// Modified Query form
	// https://github.com/ioguix/pgsql-bloat-estimation/blob/master/table/table_bloat.sql

	TableBloatSql = `
SELECT schemaname || '.' || tblname AS key, schemaname as schema, tblname as table,
 --  bs*tblpages AS real_size,
 --  (tblpages-est_tblpages)*bs AS extra_size,
 -- CASE WHEN tblpages - est_tblpages > 0
 --   THEN 100 * (tblpages - est_tblpages)/tblpages::float
 --   ELSE 0
 -- END AS extra_ratio, fillfactor,
  CASE WHEN (tblpages-est_tblpages_ff)*bs > 0
    THEN ((tblpages-est_tblpages_ff)*bs)::bigint
    ELSE 0
  END AS bloat_bytes,
  CASE WHEN tblpages - est_tblpages_ff > 0
    THEN round((100 * (tblpages - est_tblpages_ff)/tblpages::float)::numeric, 1)
    ELSE 0
  END AS bloat_ratio
-- , is_na
FROM (
  SELECT ceil( reltuples / ( (bs-page_hdr)/tpl_size ) ) + ceil( toasttuples / 4 ) AS est_tblpages,
    ceil( reltuples / ( (bs-page_hdr)*fillfactor/(tpl_size*100) ) ) + ceil( toasttuples / 4 ) AS est_tblpages_ff,
    tblpages, fillfactor, bs, tblid, schemaname, tblname, heappages, toastpages, is_na
  FROM (
    SELECT
      ( 4 + tpl_hdr_size + tpl_data_size + (2*ma)
        - CASE WHEN tpl_hdr_size%ma = 0 THEN ma ELSE tpl_hdr_size%ma END
        - CASE WHEN ceil(tpl_data_size)::int%ma = 0 THEN ma ELSE ceil(tpl_data_size)::int%ma END
      ) AS tpl_size, bs - page_hdr AS size_per_block, (heappages + toastpages) AS tblpages, heappages,
      toastpages, reltuples, toasttuples, bs, page_hdr, tblid, schemaname, tblname, fillfactor, is_na
    FROM (
      SELECT
        tbl.oid AS tblid, ns.nspname AS schemaname, tbl.relname AS tblname, tbl.reltuples,
        tbl.relpages AS heappages, coalesce(toast.relpages, 0) AS toastpages,
        coalesce(toast.reltuples, 0) AS toasttuples,
        coalesce(substring(
          array_to_string(tbl.reloptions, ' ')
          FROM '%fillfactor=#"__#"%' FOR '#')::smallint, 100) AS fillfactor,
        current_setting('block_size')::numeric AS bs,
        CASE WHEN version()~'mingw32' OR version()~'64-bit|x86_64|ppc64|ia64|amd64' THEN 8 ELSE 4 END AS ma,
        24 AS page_hdr,
        23 + CASE WHEN MAX(coalesce(null_frac,0)) > 0 THEN ( 7 + count(*) ) / 8 ELSE 0::int END
          + CASE WHEN tbl.relhasoids THEN 4 ELSE 0 END AS tpl_hdr_size,
        sum( (1-coalesce(s.null_frac, 0)) * coalesce(s.avg_width, 1024) ) AS tpl_data_size,
        bool_or(att.atttypid = 'pg_catalog.name'::regtype) AS is_na
      FROM pg_attribute AS att
        JOIN pg_class AS tbl ON att.attrelid = tbl.oid
        JOIN pg_namespace AS ns ON ns.oid = tbl.relnamespace
        JOIN pg_stats AS s ON s.schemaname=ns.nspname
          AND s.tablename = tbl.relname AND s.inherited=false AND s.attname=att.attname
        LEFT JOIN pg_class AS toast ON tbl.reltoastrelid = toast.oid
      WHERE att.attnum > 0 AND NOT att.attisdropped
        AND tbl.relkind = 'r'
      GROUP BY 1,2,3,4,5,6,7,8,9,10, tbl.relhasoids
      ORDER BY 2,3
    ) AS s
  ) AS s2
) AS s3
`

	NumberOfConnectionSql = `
SELECT numbackends as num_connections FROM pg_stat_database WHERE datname = current_database()
`

	DatabaseSizeSql = `
 SELECT
     SUM(table_size)::bigint AS table_size,
     SUM(indexes_size)::bigint AS index_size,
     SUM(total_size)::bigint AS total_size,
     round(100 * SUM(indexes_size)/SUM(total_size), 1) as index_ratio
 FROM (
     SELECT
         table_name,
         pg_table_size(table_name) AS table_size,
         pg_indexes_size(table_name) AS indexes_size,
         pg_total_relation_size(table_name) AS total_size
     FROM (
         SELECT ('"' || table_schema || '"."' || table_name || '"') AS table_name
         FROM information_schema.tables
     ) AS all_tables
     ORDER BY total_size DESC
 ) AS pretty_sizes;
 `
)
