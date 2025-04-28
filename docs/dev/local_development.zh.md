# 设置开发环境

Gogs使用[Go](https://golang.org/)语言编写，如果您还未学习过Go，请先完成[A Tour of Go](https://tour.golang.org/)教程！

## 目录

- [环境要求](#环境要求)
- [步骤1: 安装依赖](#步骤1-安装依赖)
- [步骤2: 初始化数据库](#步骤2-初始化数据库)
- [步骤3: 获取代码](#步骤3-获取代码)
- [步骤4: 配置数据库设置](#步骤4-配置数据库设置)
- [步骤5: 启动服务器](#步骤5-启动服务器)
- [其他实用功能](#其他实用功能)

## 环境要求

Gogs构建为单一可执行文件，设计为跨平台运行。因此，您可以在任何主流平台上开发Gogs。

## 步骤1: 安装依赖

Gogs需要以下依赖项：

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) (v1.8.3或更高)
- [Go](https://golang.org/doc/install) (v1.20或更高)
- [Less.js](http://lesscss.org/usage/#command-line-usage-installing)
- [Task](https://github.com/go-task/task) (v3)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [go-mockgen](https://github.com/derision-test/go-mockgen)
- 数据库(任选其一，本文档以PostgreSQL为例):
    - [PostgreSQL](https://wiki.postgresql.org/wiki/Detailed_installation_guides) (v9.6或更高)
    - [MySQL](https://dev.mysql.com/downloads/mysql/) 使用`ENGINE=InnoDB` (v5.7或更高)
    - [SQLite3](https://www.sqlite.org/index.html)
    - [TiDB](https://github.com/pingcap/tidb)

### macOS

1. 安装[Homebrew](https://brew.sh/)
1. 安装依赖:

    ```bash
    brew install go postgresql git npm go-task/tap/go-task
    npm install -g less
    npm install -g less-plugin-clean-css
    go install github.com/derision-test/go-mockgen/cmd/go-mockgen@v1.3.3
    go install golang.org/x/tools/cmd/goimports@latest
    ```

1. 配置PostgreSQL开机自启:

    ```bash
    brew services start postgresql
    ```

1. 确保PostgreSQL命令行客户端`psql`在`$PATH`中。
   Homebrew默认不会将其添加到PATH。Homebrew会在`brew info postgresql`的"Caveats"部分给出添加`psql`到PATH的命令。或者您可以使用下面的命令。根据您的Homebrew前缀(下面的`/usr/local`)和shell(bash)，可能需要调整。

    ```bash
    hash psql || { echo 'export PATH="/usr/local/opt/postgresql/bin:$PATH"' >> ~/.bash_profile }
    source ~/.bash_profile
    ```

### Ubuntu

1. 添加软件包仓库:

    ```bash
    curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
    ```

1. 更新仓库:

    ```bash
    sudo apt-get update
    ```

1. 安装依赖:

    ```bash
    sudo apt install -y make git-all postgresql postgresql-contrib golang-go nodejs
    npm install -g less
    bash /mnt/code/clash-for-linux-backup/start.sh
    source /etc/profile.d/clash.sh
    go install github.com/go-task/task/v3/cmd/task@latest
    go install github.com/derision-test/go-mockgen/cmd/go-mockgen@v1.3.3
    go install golang.org/x/tools/cmd/goimports@latest
    ```

1. 配置开机自启服务:

    ```bash
    sudo systemctl enable postgresql
    ```

## 步骤2: 初始化数据库

您需要一个全新的Postgres数据库和一个对该数据库拥有完全权限的用户。

1. 为当前Unix用户创建数据库:

    ```bash
    # Linux用户需要先切换到postgres用户
    sudo su - postgres
    ```

    ```bash
    createdb
    ```

2. 创建Gogs用户和密码:

    ```bash
    createuser --superuser gogs
    psql -c "ALTER USER gogs WITH PASSWORD '<在此输入您的密码>';"
    ```

3. 创建Gogs数据库

    ```bash
    createdb --owner=gogs --encoding=UTF8 --template=template0 gogs
    ```

## 步骤3: 获取代码

通常您不需要完整克隆，所以设置`--depth`为`10`:

```bash
git clone --depth 10 https://github.com/gogs/gogs.git
```

**注意** 该仓库已启用Go模块，请克隆到`$GOPATH`之外的目录。

## 步骤4: 配置数据库设置

在仓库内创建`custom/conf/app.ini`文件并添加以下配置(`custom/`目录下的所有文件都会覆盖默认文件且被`.gitignore`排除):

```ini
[database]
TYPE = postgres
HOST = 127.0.0.1:5432
NAME = gogs
USER = gogs
PASSWORD = <在此输入您的密码>
SSL_MODE = disable
```

## 步骤5: 启动服务器

以下命令将启动Web服务器并在任何Go文件更改时自动重新编译和重启服务器:

```bash
task web --watch
```

**注意** 如果您修改了`conf/`、`template/`或`public/`目录下的任何文件，请确保之后运行`task generate`!

## 其他实用功能

### 从磁盘加载HTML模板和静态文件

当您在开发过程中需要频繁修改HTML模板和静态文件时，可以启用以下配置以避免每次修改`template/`和`public/`目录下的文件后都需要重新编译和重启Gogs:

```ini
RUN_MODE = dev

[server]
LOAD_ASSETS_FROM_DISK = true
```

### 离线开发

有时您可能想在飞机、火车或海滩等没有WiFi的环境下开发Gogs。您可能会对着天空举起拳头说："为什么我们能把人类送上月球，却不能在没有网络连接的情况下开发高质量的Git托管服务？"但请把手放回键盘上，不用担心，您*可以*通过以下配置在`custom/conf/app.ini`中启用离线模式开发Gogs:

```ini
[server]
OFFLINE_MODE = true
