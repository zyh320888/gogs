# !!! NEVER EVER MODIFY THIS FILE !!!
# !!! PLEASE MAKE CHANGES ON CORRESPONDING CUSTOM CONFIG FILE !!!
# !!! IF YOU ARE PACKAGING PROVIDER, PLEASE MAKE OWN COPY OF IT !!!

; 应用程序的品牌名称，可以是您的公司或团队名称
BRAND_NAME = Gogs
; 运行应用程序的系统用户。在Windows上无效，
; 其他系统上应与$USER环境变量值匹配
RUN_USER = git
; 应用程序运行模式，可以是"dev"、"prod"或"test"
RUN_MODE = dev

[server]
; 应用程序对外公开的URL
EXTERNAL_URL = %(PROTOCOL)s://%(DOMAIN)s:%(HTTP_PORT)s/
; 应用程序对外公开的域名
DOMAIN = localhost
; 用于直接访问应用程序的协议
; 当前支持"http"、"https"、"fcgi"和"unix"
PROTOCOL = http
; 应用程序监听的地址
HTTP_ADDR = 0.0.0.0
; 应用程序监听的端口号
HTTP_PORT = 3000
; 生成步骤:
; $ ./gogs cert -ca=true -duration=8760h0m0s -host=myhost.example.com
;
; 或者从Windows证书存储导出的.pfx文件(不要
; 忘记导出私钥):
; $ openssl pkcs12 -in cert.pfx -out cert.pem -nokeys
; $ openssl pkcs12 -in cert.pfx -out key.pem -nocerts -nodes
CERT_FILE = custom/https/cert.pem
KEY_FILE = custom/https/key.pem
; 允许的最低TLS版本，当前支持"TLS10"、"TLS11"、"TLS12"和"TLS13"
TLS_MIN_VERSION = TLS12
; 通过Unix域套接字服务时的文件权限
UNIX_SOCKET_PERMISSION = 666
; 工作节点(如SSH更新)访问Web服务的本地(DMZ)URL
; 大多数情况下无需修改默认值
; 仅当SSH服务器节点与HTTP节点不同时才需要修改
LOCAL_ROOT_URL = %(PROTOCOL)s://%(HTTP_ADDR)s:%(HTTP_PORT)s/

; 是否禁用静态文件的CDN使用
OFFLINE_MODE = false
; 是否禁用路由日志
DISABLE_ROUTER_LOG = true
; 是否启用应用级GZIP压缩
ENABLE_GZIP = false

; 存储应用程序特定数据的路径
APP_DATA_PATH = data
; 是否从磁盘而非嵌入式bindata加载资源(如"conf"、"templates"、"public")
LOAD_ASSETS_FROM_DISK = false

; 匿名用户的着陆页URL，该值不应包含
; 由反向代理处理的子路径
LANDING_URL = /

; 是否完全禁用对应用程序的SSH访问
DISABLE_SSH = false
; SSH克隆URL中暴露的域名
SSH_DOMAIN = %(DOMAIN)s
; SSH克隆URL中暴露的端口号
SSH_PORT = 22
; SSH根目录路径，默认为"$HOME/.ssh"
SSH_ROOT_PATH =
; ssh-keygen路径，默认为"ssh-keygen"并由shell决定调用哪个
SSH_KEYGEN_PATH = ssh-keygen
; 使用ssh-keygen测试公钥时创建临时文件的目录，
; 默认为系统临时目录
SSH_KEY_TEST_PATH =
; 是否检查对应类型公钥的最小尺寸
MINIMUM_KEY_SIZE_CHECK = false
; 是否在启动时重写"~/.ssh/authorized_keys"文件，使用内置SSH服务器时忽略
REWRITE_AUTHORIZED_KEYS_AT_START = false
; 是否启动内置SSH服务器
START_SSH_SERVER = false
; 内置SSH服务器监听的网络接口
SSH_LISTEN_HOST = 0.0.0.0
; 内置SSH服务器监听的端口号
SSH_LISTEN_PORT = %(SSH_PORT)s
; 内置SSH服务器连接接受的加密算法列表
SSH_SERVER_CIPHERS = aes128-ctr, aes192-ctr, aes256-ctr, aes128-gcm@openssh.com, arcfour256, arcfour128
; 内置SSH服务器连接接受的MAC算法列表
SSH_SERVER_MACS = hmac-sha2-256-etm@openssh.com, hmac-sha2-256, hmac-sha1
; 内置SSH服务器连接接受的密钥交换算法列表
SSH_SERVER_ALGORITHMS = rsa, ecdsa, ed25519

; 定义允许的算法及其最小密钥长度(使用-1禁用某类型)
[ssh.minimum_key_sizes]
ED25519 = 256
ECDSA   = 256
RSA     = 2048
DSA     = 1024

[repository]
; 存储托管仓库的根路径，默认为"~/gogs-repositories"
ROOT =
; 服务器支持的脚本类型，有时可以是"sh"
SCRIPT_TYPE = bash
; 未识别字符集的默认ANSI字符集
ANSI_CHARSET =
; 是否强制每个新仓库为私有
FORCE_PRIVATE = false
; 用户可创建仓库的全局数量限制，-1表示无限制
MAX_CREATION_LIMIT = -1
; 优先显示在列表顶部的许可证
; 名称必须匹配"conf/license"或"custom/conf/license"中的文件名
PREFERRED_LICENSES = Apache License 2.0, MIT License
; 是否禁用通过HTTP/HTTPS协议与仓库的Git交互
DISABLE_HTTP_GIT = false
; 是否启用通过服务器本地路径迁移仓库的能力
ENABLE_LOCAL_PATH_MIGRATION = false
; 是否启用原始文件的渲染模式。存在潜在安全风险
ENABLE_RAW_FILE_RENDER_MODE = false
; 单个获取请求可同时运行的最大goroutine数量
; 通常该值取决于CPU(核心)数量
; 如果值为非正数，则匹配应用程序可用的CPU数量
COMMITS_FETCH_CONCURRENCY = 0
; 创建新仓库时的默认分支名称
DEFAULT_BRANCH = master

[repository.editor]
; 在CodeMirror编辑器中应自动换行的文件扩展名列表
; 用逗号分隔扩展名
LINE_WRAP_EXTENSIONS = .txt,.md,.markdown,.mdown,.mkd
; 具有关联预览API(如"/api/v1/markdown")的有效文件模式
; 用逗号分隔值。如果文件扩展名不匹配，编辑模式中的预览标签页不会显示
PREVIEWABLE_FILE_MODES = markdown

[repository.upload]
; 是否启用仓库文件上传
ENABLED = true
; 临时存储上传文件的路径(每次启动时此路径下的内容会被清除)
TEMP_PATH = data/tmp/uploads
; File types that are allowed to be uploaded, e.g. "image/jpeg|image/png". Leave empty to allow any file type.
ALLOWED_TYPES =
; The maximum size of each file in MB.
FILE_MAX_SIZE = 3
; The maximum number of files per upload.
MAX_FILES = 5

[database]
; 数据库后端，可以是"postgres"、"mysql"、"sqlite3"或"mssql"
; 您可以使用MySQL协议连接TiDB
TYPE = postgres
HOST = 127.0.0.1:5432
NAME = gogs
USER = gogs
PASSWORD =
; For "postgres" only
SCHEMA = public
; 仅适用于"postgres"，可以是"disable"、"require"或"verify-full"
SSL_MODE = disable
; 仅适用于"sqlite3"，请确保使用绝对路径
PATH = data/gogs.db
; 连接池的最大打开连接数
MAX_OPEN_CONNS = 30
; 连接池的最大空闲连接数
MAX_IDLE_CONNS = 30

[security]
; 是否显示安装页面，设为"true"可跳过
INSTALL_LOCK = false
; 用于加密cookie值、2FA代码等的密钥
; !!修改此值以保护您的用户数据安全!!
SECRET_KEY = !#@FDEWREWR&*(
; 自动登录的记住天数
LOGIN_REMEMBER_DAYS = 7
; 存储自动登录信息的cookie名称
COOKIE_REMEMBER_NAME = gogs_incredible
; 存储登录用户名的cookie名称
COOKIE_USERNAME = gogs_awesome
; 是否设置安全cookie
COOKIE_SECURE = false
; 是否设置cookie来指示用户登录状态
ENABLE_LOGIN_STATUS_COOKIE = false
; 存储用户登录状态的cookie名称
LOGIN_STATUS_COOKIE_NAME = login_status
; 本地网络内明确允许访问的主机名逗号分隔列表
; 使用"*"允许所有主机名
LOCAL_NETWORK_ALLOWLIST =

[email]
; 是否启用邮件服务
ENABLED = false
; 邮件主题前缀
SUBJECT_PREFIX = `[%(BRAND_NAME)s] `
; SMTP服务器及其端口，例如 smtp.mailgun.org:587, smtp.gmail.com:587, smtp.qq.com:465
; 如果端口以"465"结尾，将使用SMTPS。根据RFC 6409建议在587端口使用STARTTLS
; 如果服务器支持STARTTLS，将始终使用
HOST = smtp.mailgun.org:587
; 发件人地址(RFC 5322)。可以是纯邮箱地址或`"名称" <email@example.com>`格式
FROM = noreply@gogs.localhost
; 登录用户名
USER = noreply@gogs.localhost
; 登录密码
PASSWORD =

; 当主机名不同时是否禁用HELO操作
DISABLE_HELO =
; HELO操作的自定义主机名，默认为系统主机名
HELO_HOSTNAME =

; 是否跳过服务器证书验证。仅用于自签名证书
SKIP_VERIFY = false
; 是否使用客户端证书
USE_CERTIFICATE = false
CERT_FILE = custom/email/cert.pem
KEY_FILE = custom/email/key.pem

; 是否使用"text/plain"作为内容格式
USE_PLAIN_TEXT = false
; 发送HTML邮件时是否附加纯文本替代内容
; 用于支持旧版邮件客户端并提高垃圾邮件过滤通过率
ADD_PLAIN_TEXT_ALT = false

[auth]
; 激活码的有效期(分钟)
ACTIVATE_CODE_LIVES = 180
; 重置密码码的有效期(分钟)
RESET_PASSWORD_CODE_LIVES = 180
; 是否要求验证新增邮箱地址
; 启用此选项也会要求用户在注册时验证邮箱
REQUIRE_EMAIL_CONFIRMATION = false
; 是否禁止匿名用户访问网站
REQUIRE_SIGNIN_VIEW = false
; 是否禁用自助注册。禁用后账号必须由管理员创建
DISABLE_REGISTRATION = false
; 是否启用注册验证码
ENABLE_REGISTRATION_CAPTCHA = true

; 是否启用通过HTTP头的反向代理认证
ENABLE_REVERSE_PROXY_AUTHENTICATION = false
; 是否为反向代理认证自动创建新用户
ENABLE_REVERSE_PROXY_AUTO_REGISTRATION = false
; 用于反向代理认证的用户名HTTP头
REVERSE_PROXY_AUTHENTICATION_HEADER = X-WEBAUTH-USER

[user]
; 是否启用用户邮件通知
ENABLE_EMAIL_NOTIFICATION = false

[session]
; 会话存储提供者，可选"memory"、"file"或"redis"
PROVIDER = memory
; 各提供者的配置：
; - memory: 目前不需要配置
; - file: 会话文件路径，如`data/sessions`
; - redis: 网络=tcp,地址=:6379,密码=macaron,数据库=0,连接池大小=100,空闲超时=180,tls=是
PROVIDER_CONFIG = data/sessions
; 存储会话标识符的cookie名称
COOKIE_NAME = i_like_gogs
; 是否仅在HTTPS下设置cookie
COOKIE_SECURE = false
; 会话数据的垃圾回收间隔(秒)
GC_INTERVAL = 3600
; 会话的最大生命周期(秒)
MAX_LIFE_TIME = 86400
; CSRF令牌的cookie名称
CSRF_COOKIE_NAME = _csrf

[cache]
; 缓存适配器，可选"memory"、"redis"或"memcache"
ADAPTER = memory
; 仅对"memory"有效，垃圾回收间隔(秒)
INTERVAL = 60
; 对于"redis"和"memcache"，连接主机地址：
; - redis: 网络=tcp,地址=:6379,密码=macaron,数据库=0,连接池大小=100,空闲超时=180
; - memcache: `127.0.0.1:11211`
HOST =

[http]
; "Access-Control-Allow-Origin"头的值，默认为不设置
ACCESS_CONTROL_ALLOW_ORIGIN =

[lfs]
; 新对象上传的存储后端
STORAGE = local
; 本地文件系统中存储LFS对象的根路径
OBJECTS_PATH = data/lfs-objects

[attachment]
; 是否启用通用附件上传功能
ENABLED = true
; 文件系统中存储附件的路径
PATH = data/attachments
; 允许上传的文件类型，如"image/jpeg|image/png"。留空则允许所有文件类型
ALLOWED_TYPES = image/jpeg|image/png
; 每个文件的最大大小(MB)
MAX_SIZE = 4
; 每次上传的最大文件数量
MAX_FILES = 5

[release.attachment]
; 是否启用发布版本附件上传功能
ENABLED = true
; 允许上传的文件类型，如"image/jpeg|image/png"。留空则允许所有文件类型
ALLOWED_TYPES = */*
; 每个文件的最大大小(MB)
MAX_SIZE = 32
; 每次上传的最大文件数量
MAX_FILES = 10

[time]
; 指定完整日期输出的格式
; 可选值包括:
; ANSIC, UnixDate, RubyDate, RFC822, RFC822Z, RFC850, RFC1123, RFC1123Z, RFC3339, RFC3339Nano, Kitchen, Stamp, StampMilli, StampMicro 和 StampNano
; 更多格式信息请参考 http://golang.org/pkg/time/#pkg-constants
FORMAT = RFC1123

[picture]
; 文件系统中存储用户头像的路径
AVATAR_UPLOAD_PATH = data/avatars
; 文件系统中存储仓库头像的路径
REPOSITORY_AVATAR_UPLOAD_PATH = data/repo-avatars
; 中国用户可以使用自定义头像源，如 http://cn.gravatar.com/avatar/
GRAVATAR_SOURCE = gravatar
; 是否禁用Gravatar，离线模式下此值将强制为true
DISABLE_GRAVATAR = false
; 是否启用联合头像查找，使用DNS发现与邮箱关联的头像
; 详情参见 https://www.libravatar.org
; 离线模式或禁用Gravatar时此值将强制为false
ENABLE_FEDERATED_AVATAR = false

[markdown]
; 是否启用硬换行扩展
ENABLE_HARD_LINE_BREAK = false
; 渲染Markdown时允许作为链接的自定义URL方案列表
; 例如："git"(对应"git://")和"magnet"(对应"magnet://")
CUSTOM_URL_SCHEMES =
; 应作为Markdown渲染/编辑的文件扩展名列表
; 用逗号分隔扩展名。要将无扩展名文件渲染为markdown，只需放一个逗号
FILE_EXTENSIONS = .md,.markdown,.mdown,.mkd

[smartypants]
; 是否启用Smartypants扩展
ENABLED = false
FRACTIONS = true  ; 转换分数
DASHES = true     ; 转换破折号
LATEX_DASHES = true  ; 转换LaTeX破折号
ANGLED_QUOTES = true  ; 转换尖括号引号

[admin]
; 是否禁止普通(非管理员)用户创建组织
DISABLE_REGULAR_ORG_CREATION = false

[webhook]
; 用户可用的Webhook类型列表，可选"gogs"、"slack"、"discord"、"dingtalk"
TYPES = gogs, slack, discord, dingtalk
; 投递超时时间(秒)
DELIVER_TIMEOUT = 15
; 是否允许不安全的证书
SKIP_TLS_VERIFY = false
; 每页显示的历史信息数量
PAGING_NUM = 10

; General settings of loggers.
[log]
; 所有日志文件的根路径，默认为"log/"子目录
ROOT_PATH =
; 可选"console"、"file"、"slack"和"discord"
; 使用逗号分隔多种模式，例如"console, file"
MODE = console
; 通道缓冲区长度，如果不了解请保持默认
BUFFER_LEN = 100
; 可选"Trace"、"Info"、"Warn"、"Error"、"Fatal"，默认为"Trace"
LEVEL = Trace

; For "console" mode only
[log.console]
; 注释掉则继承上级设置
; LEVEL =

; For "file" mode only
[log.file]
; Comment out to inherit
; LEVEL =
; 是否启用自动日志轮转(控制以下选项的开关)
LOG_ROTATE = true
; 是否按天分割日志文件
DAILY_ROTATE = true
; 单个文件的最大大小位移，默认28表示1<<28=256MB
MAX_SIZE_SHIFT = 28
; 单个文件的最大行数
MAX_LINES = 1000000
; 日志文件的过期天数(超过最大天数后删除)
MAX_DAYS = 7

; For "slack" mode only
[log.slack]
; Comment out to inherit
; LEVEL =
; Webhook URL地址
URL =

[log.discord]
; Comment out to inherit
; LEVEL =
; Webhook URL地址
URL =
; 通知中显示的用户名
USERNAME = %(BRAND_NAME)s

[log.xorm]
; 启用文件轮转
ROTATE = true
; 每天轮转
ROTATE_DAILY = true
; 当文件大小超过x MB时轮转
MAX_SIZE = 100
; 日志文件保留的最大天数
MAX_DAYS = 3

[log.gorm]
; 是否启用文件轮转
ROTATE = true
; 是否每天轮转文件
ROTATE_DAILY = true
; 下次轮转前的最大文件大小(MB)
MAX_SIZE = 100
; The maximum days to keep files.
MAX_DAYS = 3

[cron]
; 启用定期运行cron任务
ENABLED = true
; Gogs启动时运行cron任务
RUN_AT_START = false

[cron.update_mirrors]
; 定义镜像同步器检查是否需要同步的频率(基于镜像更新间隔)
SCHEDULE = @every 10m

; 仓库健康检查
[cron.repo_health_check]
SCHEDULE = @every 24h
TIMEOUT = 60s
; 'git fsck'命令的参数，例如 "--unreachable --tags"
; 更多信息见 http://git-scm.com/docs/git-fsck/1.7.5
ARGS =

; 检查仓库统计信息
[cron.check_repo_stats]
RUN_AT_START = true
SCHEDULE = @every 24h

; 清理仓库归档文件
[cron.repo_archive_cleanup]
RUN_AT_START = false
SCHEDULE = @every 24h
; 检查归档文件是否需要清理的时间范围
OLDER_THAN = 24h

[git]
; 禁用添加和删除变更的高亮显示
DISABLE_DIFF_HIGHLIGHT = false
; diff视图中显示的最大文件数
MAX_GIT_DIFF_FILES = 100
; diff视图中单个文件允许的最大行数
MAX_GIT_DIFF_LINES = 1000
; diff视图中单行允许的最大字符数
MAX_GIT_DIFF_LINE_CHARACTERS = 2000
; 'git gc'命令的参数，例如 "--aggressive --auto"
; 更多信息见 http://git-scm.com/docs/git-gc/1.7.5
GC_ARGS =

; 操作超时时间(秒)
[git.timeout]
MIGRATE = 600
MIRROR = 300
CLONE = 300
PULL = 300
DIFF = 60
GC = 60

[mirror]
; 定义镜像下次同步的默认间隔时间(小时)(在成功同步后)
; 可以在设置中为每个镜像仓库单独覆盖此值
DEFAULT_INTERVAL = 8

[api]
; 每页返回的最大项目数
MAX_RESPONSE_ITEMS = 50

[ui]
; 探索页面显示的仓库数量
EXPLORE_PAGING_NUM = 20
; 每页显示的问题数量
ISSUE_PAGING_NUM = 10
; 活动动态中显示的最大提交数
FEED_MAX_COMMIT_NUM = 5
; "theme-color"元标签的值，用于Android >= 5.0
; 无效颜色如"none"或"disable"将使用默认样式
; 更多信息: https://developers.google.com/web/updates/2014/11/Support-for-theme-color-in-Chrome-39-for-Android
THEME_COLOR_META_TAG = `#ff5343`
; 可显示文件的最大字节数(默认为8MB)
MAX_DISPLAY_FILE_SIZE = 8388608

[ui.admin]
; 每页显示的用户数量
USER_PAGING_NUM = 50
; 每页显示的仓库数量
REPO_PAGING_NUM = 50
; 每页显示的通知数量
NOTICE_PAGING_NUM = 25
; 每页显示的组织数量
ORG_PAGING_NUM = 50

[ui.user]
; 每页显示的仓库数量
REPO_PAGING_NUM = 15
; 每页显示的动态消息数量
NEWS_FEED_PAGING_NUM = 20
; 每页显示的提交数量
COMMITS_PAGING_NUM = 30

[prometheus]
; 是否启用Prometheus指标
ENABLED = true
; 是否启用HTTP基本认证来保护指标数据
ENABLE_BASIC_AUTH = false
; HTTP基本认证的用户名
BASIC_AUTH_USERNAME =
; HTTP基本认证的密码
BASIC_AUTH_PASSWORD =

; 扩展名到高亮类的映射
; 例如 .toml=ini
[highlight.mapping]

[i18n]
LANGS = en-US,zh-CN,zh-HK,zh-TW,de-DE,fr-FR,nl-NL,lv-LV,ru-RU,ja-JP,es-ES,pt-BR,pl-PL,bg-BG,it-IT,fi-FI,tr-TR,cs-CZ,sr-SP,sv-SE,ko-KR,gl-ES,uk-UA,en-GB,hu-HU,sk-SK,id-ID,fa-IR,vi-VN,pt-PT,mn-MN,ro-RO
NAMES = English,简体中文,繁體中文（香港）,繁體中文（臺灣）,Deutsch,français,Nederlands,latviešu,русский,日本語,español,português do Brasil,polski,български,italiano,suomi,Türkçe,čeština,српски,svenska,한국어,galego,українська,English (United Kingdom),Magyar,Slovenčina,Indonesian,Persian,Vietnamese,Português,Монгол,Română

; Used for jQuery DateTimePicker,
; list of supported languages in https://xdsoft.net/jqplugins/datetimepicker/#lang
[i18n.datelang]
en-US = en
zh-CN = zh
zh-HK = zh-TW
zh-TW = zh-TW
de-DE = de
fr-FR = fr
nl-NL = nl
lv-LV = lv
ru-RU = ru
ja-JP = ja
es-ES = es
pt-BR = pt-BR
pl-PL = pl
bg-BG = bg
it-IT = it
fi-FI = fi
tr-TR = tr
cs-CZ = cs-CZ
sr-SP = sr
sv-SE = sv
ko-KR = ko
gl-ES = gl
uk-UA = uk
en-GB = en-GB
hu-HU = hu
sk-SK = sk
id-ID = id
fa-IR = fa
vi-VN = vi
pt-PT = pt
mn-MN = mn
ro-RO = ro

[other]
SHOW_FOOTER_BRANDING = false
; 在页脚显示模板执行时间
SHOW_FOOTER_TEMPLATE_LOAD_TIME = true
