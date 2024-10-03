# eventrunner

## ClickHouse

```sql
CREATE TABLE IF NOT EXISTS events (
    id String,
    source String,
    type String,
    subject Nullable(String),
    time DateTime,
    data_contenttype String,
    data String,
    specversion String
) ENGINE = MergeTree()
ORDER BY (time, id);
```