# 初始化

```
goo_clickhouse.Init(clickhouse.Config{
    Driver:   "clickhouse",
    Addr:     "192.168.1.100:9000",
    User:     "root",
    Password: "123456",
    Database: "test",
})
```

# 创建数据库

```
CREATE DATABASE IF NOT EXISTS test;
```

# 创建表

```
sqlstr := `
    CREATE TABLE IF NOT EXISTS user
    (
        name String,
        gender String,
        birthday Date
    )
    ENGINE = MergeTree()
    ORDER BY (name, gender)
    PARTITION BY toYYYYMM(birthday)
`
if _, err := goo_clickhouse.DB().Exec(sqlstr); err != nil {
    log.Fatal(err)
}
```

# 添加数据

```
func insert() {
	var (
		tx, _   = goo_clickhouse.DB().Begin()
		stmt, _ = tx.Prepare(`INSERT INTO user(name, gender) VALUES(?, ?)`)
	)

    data := []interface{}{"", ""}
    if _, err := stmt.Exec(data...); err != nil {
        goo_log.Error(err)
        return
    }

    if err := tx.Commit(); err != nil {
        goo_log.Error(err)
	}
}
```

# 查询数据

```
rows, err := DB().Query("SELECT name, gender FROM user")
if err != nil {
    goo_log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var (
        name string
        gender string
    )
    if err := rows.Scan(&name, &gender); err != nil {
        log.Fatal(err)
    }
    fmt.Println(name, gender)
}

if err := rows.Err(); err != nil {
    goo_log.Fatal(err)
}
```

# 删除表

```
if _, err := DB().Exec("DROP TABLE user"); err != nil {
	goo_log.Fatal(err)
}
```