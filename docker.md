要在 **Docker** 中运行你的 `crontab_manager`，你需要：

1. **创建 Dockerfile**
2. **构建 Docker 镜像**
3. **运行 Docker 容器**
4. **持久化 crontab（可选）**
5. **使用 `docker-compose`（可选）**

---

## **1. 创建 `Dockerfile`**
在你的项目根目录下创建一个 `Dockerfile`：

```dockerfile
# 使用更小的基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制所有代码
COPY . .

# 构建 Go 代码
RUN go mod tidy && go build -o crontab_manager

# 生产环境镜像
FROM alpine:latest

WORKDIR /app

# 添加必要的依赖
RUN apk --no-cache add bash ca-certificates tzdata cronie

# 复制二进制文件
COPY --from=builder /app/crontab_manager /usr/local/bin/crontab_manager

# 允许 crontab 命令执行
RUN chmod +x /usr/local/bin/crontab_manager

# 运行 Web 服务器
CMD ["/usr/local/bin/crontab_manager", "-port", "8080", "-user", "admin", "-pass", "password"]
```

---

## **2. 构建 Docker 镜像**
```sh
docker build -t crontab_manager .
```

---

## **3. 运行 Docker 容器**
```sh
docker run -d --name crontab_manager \
  -p 8080:8080 \
  --restart unless-stopped \
  crontab_manager
```
### **说明**
- `-d`：后台运行
- `--name crontab_manager`：容器名称
- `-p 8080:8080`：映射 Web 端口
- `--restart unless-stopped`：支持自动重启

---

## **4. 持久化 crontab（可选）**
默认情况下，Docker 容器重启后 `crontab` 任务会丢失，你可以使用 `volume` 共享 `crontab` 数据：
```sh
docker run -d --name crontab_manager \
  -p 8080:8080 \
  -v /var/spool/cron/crontabs:/var/spool/cron/crontabs \
  --restart unless-stopped \
  crontab_manager
```
**这样即使容器重启，crontab 任务仍然保留！**

---

## **5. 使用 `docker-compose`（推荐）**
创建 `docker-compose.yml`：
```yaml
version: "3"
services:
  crontab_manager:
    image: crontab_manager
    container_name: crontab_manager
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - /var/spool/cron/crontabs:/var/spool/cron/crontabs
    command: ["/usr/local/bin/crontab_manager", "-port", "8080", "-user", "admin", "-pass", "password"]
```

### **启动**
```sh
docker-compose up -d
```

---

### **6. 管理容器**
| 操作        | 命令 |
|-------------|------------------------------|
| **查看运行状态** | `docker ps` |
| **进入容器** | `docker exec -it crontab_manager sh` |
| **查看日志** | `docker logs -f crontab_manager` |
| **重启容器** | `docker restart crontab_manager` |
| **停止容器** | `docker stop crontab_manager` |

---

## **总结**
✅ **Docker 化你的 Go 程序**  
✅ **支持 Web 端口映射 & BA 认证**  
✅ **支持 crontab 任务持久化**  
✅ **可选 `docker-compose` 方式管理**  

现在你可以在任何服务器上运行 `crontab_manager`，而无需额外安装 Go 或其他依赖！ 🚀