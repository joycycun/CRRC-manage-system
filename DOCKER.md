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

默认账号仍使用当前系统里的账号，密码都是之前配置的值，例如：

```text
admin / 123456
卢进 / 123456
郑宇 / 123456
王宇 / 123456
丁sir / 123456
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
```

如果使用 Docker Compose 内部数据库端口映射为 `3307:3306`，则把端口改为：

```bash
mysql -h 127.0.0.1 -P 3307 -ucrrc_user -p123456 crrc_pm < docker/mysql/migrations/20260623_runtime_schema.sql
```
