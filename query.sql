Copy (
  WITH models AS (
    WITH data AS (
      SELECT
          replace(initcap(table_name::text), '_', '') table_name,
          CASE data_type
          WHEN 'timestamp without time zone' THEN 'sql.NullTime'
          WHEN 'timestamp with time zone' THEN 'sql.NullTime'
          WHEN 'boolean' THEN 'sql.NullBool'
          WHEN 'integer' THEN 'sql.NullInt64'
          -- add your own type converters as needed or it will default to 'string'
          ELSE 'sql.NullString'
          END AS type_info,
          CASE
          WHEN column_name ~ '^[0-9]' THEN 'X' || column_name
          ELSE column_name
          END AS col
      FROM information_schema.columns
      WHERE table_schema IN ('public')
      ORDER BY table_schema, table_name, ordinal_position
    )
      SELECT table_name, STRING_AGG(E'\t' || replace(initcap(col::text), '_', '')  || E'\t' || type_info, E'\n') fields
      FROM data
      GROUP BY table_name
  )
  SELECT 'type ' || table_name || E' struct {\n' || fields || E'\n}' models
  FROM models ORDER BY 1
) TO STDOUT;