在 Linux 上，你可以使用 **systemd** 创建一个服务，以便让你的 `crontab_manager` 以守护进程方式运行，并在系统启动时自动启动。

---

## **步骤 1：创建 systemd 服务文件**
使用 `vim` 或 `nano` 创建一个新的服务文件：
```sh
sudo nano /etc/systemd/system/crontab_manager.service
```

然后，添加以下内容：
```ini
[Unit]
Description=Crontab Manager Web Service
After=network.target

[Service]
ExecStart=/usr/local/bin/crontab_manager -p 10010 -u admin -pwd password -d
Restart=always
RestartSec=5
KillMode=process
WorkingDirectory=/usr/local/bin
StandardOutput=append:/var/log/crontab_manager.log
StandardError=append:/var/log/crontab_manager.log

[Install]
WantedBy=multi-user.target
```

---

## **步骤 2：设置权限**
保存文件后，执行以下命令以确保 systemd 识别它：
```sh
sudo chmod 644 /etc/systemd/system/crontab_manager.service
```

---

## **步骤 3：启动并启用服务**
**1. 重新加载 systemd 配置**
```sh
sudo systemctl daemon-reload
```

**2. 启动服务**
```sh
sudo systemctl start crontab_manager
```

**3. 设置开机自启**
```sh
sudo systemctl enable crontab_manager
```

---

## **步骤 4：检查服务状态**
你可以使用以下命令检查服务是否正常运行：
```sh
sudo systemctl status crontab_manager
```
如果成功，你应该会看到类似下面的输出：
```
● crontab_manager.service - Crontab Manager Web Service
   Loaded: loaded (/etc/systemd/system/crontab_manager.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2025-02-08 10:00:00 UTC; 2min ago
 Main PID: 12345 (crontab_manager)
   Tasks: 5
   Memory: 10M
   CGroup: /system.slice/crontab_manager.service
           └─12345 /usr/local/bin/crontab_manager -port 8080 -user admin -pass password -daemon
```

---

## **步骤 5：管理服务**
| 操作        | 命令 |
|-------------|--------------------------------|
| **启动服务**  | `sudo systemctl start crontab_manager` |
| **停止服务**  | `sudo systemctl stop crontab_manager` |
| **重启服务**  | `sudo systemctl restart crontab_manager` |
| **查看日志**  | `journalctl -u crontab_manager -f` |

---

### **至此，你的 `crontab_manager` 已成功作为服务运行，并支持守护进程模式！🚀**
这样，每次服务器重启后，`crontab_manager` 都会自动启动，不需要手动运行。