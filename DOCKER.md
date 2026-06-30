# Docker 部署

已内置一个干净的 MySQL 初始化脚本：

- 业务表数据已清空。
- `users`、`roles`、`permissions`、`user_roles`、`role_permissions` 会保留并初始化。

启动：

```bash
docker compose up -d --build
```

访问：

```text
http://localhost
```

上传文件存储：

- 新上传的业务文件不再写入 MySQL 的 `file_data`，会保存到后端容器固定目录 `/app/uploads`。
- Docker Compose 已将 `/app/uploads` 挂载到 `uploads_data` 卷，重建容器不会丢失上传文件。
- 文件按类型分目录保存：`documents`、`texts`、`pdf`、`spreadsheets`、`archives`、`images`、`others`。
- 数据库只保存文件名、类型、大小、分类和磁盘路径；旧数据库里的 BLOB 文件仍可兼容下载。
- 如需把旧数据库中已有的 BLOB 文件迁移到 `/app/uploads` 并清空数据库文件内容，运行：

```bash
cd backend
UPLOAD_ROOT=/app/uploads DB_HOST=127.0.0.1 DB_PORT=3306 DB_USER=crrc_user DB_PASSWORD=123456 DB_NAME=crrc_pm go run ./cmd/migrate_file_blobs
```

Docker Compose 端口是 `3307:3306` 时，把 `DB_PORT` 改为 `3307`。

默认账号仍使用当前系统里的账号，密码都是之前配置的值，例如：

```text
admin / 123456
卢进 / 123456
郑宇 / 123456
王宇 / 123456
丁宇 / 123456
寸诗睿 / 123456
刘克英 / 123456
王洪玮 / 123456
袁晓兰 / 123456
未知 / 123456
傅建豪 / 123456
彭泉鑫 / 123456
```

如果需要重新使用初始化 SQL，需要先删除 Docker 卷：

```bash
docker compose down -v
docker compose up -d --build
```

已有数据库升级时，不要删除 Docker 卷，直接执行迁移脚本：

```bash
mysql -h 127.0.0.1 -P 3306 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260623_runtime_schema.sql
mysql -h 127.0.0.1 -P 3306 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260630_file_storage_path.sql
mysql -h 127.0.0.1 -P 3306 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260630_auth_device_updates.sql
```

如果使用 Docker Compose 内部数据库端口映射为 `3307:3306`，则把端口改为：

```bash
mysql -h 127.0.0.1 -P 3307 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260623_runtime_schema.sql
mysql -h 127.0.0.1 -P 3307 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260630_file_storage_path.sql
mysql -h 127.0.0.1 -P 3307 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260630_auth_device_updates.sql
```
