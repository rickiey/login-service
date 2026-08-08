## login-service


配置文件再 conf目录下，配置好数据库和 redis

导入数据库表配置 sql文件夹下

运行代码：
> go run cmd/login-service

编译可执行文件 ./bin 下, 需要安装 make
> make  

编译成docker镜像
> make docker

运行生成的 docker 镜像（注意，配置文件里的mysql redis 地址要让程序能在docker离找到）
> docker  run -d -p 3000:3000 -v /home/rui/workplace/go/src/login-service/conf:/conf --name=login  mysqlserv.com/common/login-service:6f76db9

## 依赖升级（2026-08）

- Go 版本：1.13 -> 1.25
- echo/v4: v4.9.0 -> v4.15.4
- viper: v1.3.1 -> v1.21.0
- logrus: v1.3.0 -> v1.9.4
- go-sql-driver/mysql: v1.4.1 -> v1.10.0
- golang.org/x/crypto: v0.17.0 -> v0.54.0
- gorm: jinzhu/gorm v1.9.2 -> gorm.io/gorm v1.31.2 (v2)
- redigo: garyburd/redigo (已废弃) -> gomodule/redigo v1.9.3
- validator: go-playground/validator.v9 -> v10

### 登录接口　

+ localhost:3000/account/v1/login 
> 方法　POST
> Header  Content-Type: application/json

body 参数 ：
```json
{"phone":"18312345678","password":"123456"}
```

返回参数

```json
{
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg2NjI1NDAsInN1YiI6MX0.Ntx1fEmoocONYEVu9Ug1r_EJNXiKwslt2x3a9h9YhHI",
    "user": {
        "email": "杨瑞",
        "enable": true,
        "id": 1,
        "name": "rui",
        "phone": "18312345678"
    }
}
```
