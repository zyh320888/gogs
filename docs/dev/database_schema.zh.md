# 表 "access"

```
  字段  |  列名  |   POSTGRESQL    |         MYSQL         |     SQLITE3       
--------+---------+-----------------+-----------------------+-------------------
  ID    | id      | BIGSERIAL       | BIGINT AUTO_INCREMENT | INTEGER           
  用户ID | user_id | BIGINT NOT NULL | BIGINT NOT NULL       | INTEGER NOT NULL  
  仓库ID | repo_id | BIGINT NOT NULL | BIGINT NOT NULL       | INTEGER NOT NULL  
  模式   | mode    | BIGINT NOT NULL | BIGINT NOT NULL       | INTEGER NOT NULL  

主键: id
索引: 
	"access_user_repo_unique" 唯一键 (user_id, repo_id)
```

# 表 "access_token"

```
      字段      |    列名    |         POSTGRESQL          |            MYSQL            |           SQLITE3            
----------------+--------------+-----------------------------+-----------------------------+------------------------------
  ID           | id           | BIGSERIAL                   | BIGINT AUTO_INCREMENT       | INTEGER                      
  用户ID       | uid          | BIGINT                      | BIGINT                      | INTEGER                      
  名称         | name         | TEXT                        | LONGTEXT                    | TEXT                         
  Sha1         | sha1         | VARCHAR(40) UNIQUE          | VARCHAR(40) UNIQUE          | VARCHAR(40) UNIQUE           
  SHA256       | sha256       | VARCHAR(64) NOT NULL UNIQUE | VARCHAR(64) NOT NULL UNIQUE | VARCHAR(64) NOT NULL UNIQUE  
  创建时间Unix | created_unix | BIGINT                      | BIGINT                      | INTEGER                      
  更新时间Unix | updated_unix | BIGINT                      | BIGINT                      | INTEGER                      

主键: id
索引: 
	"idx_access_token_user_id" (uid)
```

# 表 "action"

```
       字段       |     列名     |           POSTGRESQL           |             MYSQL              |            SQLITE3              
------------------+----------------+--------------------------------+--------------------------------+---------------------------------
  ID             | id             | BIGSERIAL                      | BIGINT AUTO_INCREMENT          | INTEGER                         
  用户ID         | user_id        | BIGINT                         | BIGINT                         | INTEGER                         
  操作类型       | op_type        | BIGINT                         | BIGINT                         | INTEGER                         
  操作用户ID     | act_user_id    | BIGINT                         | BIGINT                         | INTEGER                         
  操作用户名     | act_user_name  | TEXT                           | LONGTEXT                       | TEXT                            
  仓库ID         | repo_id        | BIGINT                         | BIGINT                         | INTEGER                         
  仓库用户名     | repo_user_name | TEXT                           | LONGTEXT                       | TEXT                            
  仓库名称       | repo_name      | TEXT                           | LONGTEXT                       | TEXT                            
  引用名称       | ref_name       | TEXT                           | LONGTEXT                       | TEXT                            
  是否私有       | is_private     | BOOLEAN NOT NULL DEFAULT FALSE | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE  
  内容           | content        | TEXT                           | LONGTEXT                       | TEXT                            
  创建时间Unix   | created_unix   | BIGINT                         | BIGINT                         | INTEGER                         

主键: id
索引: 
	"idx_action_repo_id" (repo_id)
	"idx_action_user_id" (user_id)
```

# 表 "email_address"

```
      字段      |    列名    |           POSTGRESQL           |             MYSQL              |            SQLITE3              
----------------+--------------+--------------------------------+--------------------------------+---------------------------------
  ID           | id           | BIGSERIAL                      | BIGINT AUTO_INCREMENT          | INTEGER                         
  用户ID       | uid          | BIGINT NOT NULL                | BIGINT NOT NULL                | INTEGER NOT NULL                
  邮箱         | email        | VARCHAR(254) NOT NULL          | VARCHAR(254) NOT NULL          | TEXT NOT NULL                   
  是否激活     | is_activated | BOOLEAN NOT NULL DEFAULT FALSE | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE  

主键: id
索引: 
	"email_address_user_email_unique" 唯一键 (uid, email)
	"idx_email_address_user_id" (uid)
```

# 表 "follow"

```
   字段   |  列名   |   POSTGRESQL    |         MYSQL         |     SQLITE3       
----------+-----------+-----------------+-----------------------+-------------------
  ID      | id        | BIGSERIAL       | BIGINT AUTO_INCREMENT | INTEGER           
  用户ID  | user_id   | BIGINT NOT NULL | BIGINT NOT NULL       | INTEGER NOT NULL  
  关注ID  | follow_id | BIGINT NOT NULL | BIGINT NOT NULL       | INTEGER NOT NULL  

主键: id
索引: 
	"follow_user_follow_unique" 唯一键 (user_id, follow_id)
```

# 表 "lfs_object"

```
     字段    |    列名    |      POSTGRESQL      |        MYSQL         |      SQLITE3       
-------------+------------+----------------------+----------------------+--------------------
  仓库ID     | repo_id    | BIGINT               | BIGINT               | INTEGER            
  OID        | oid        | TEXT                 | VARCHAR(191)         | TEXT               
  大小       | size       | BIGINT NOT NULL      | BIGINT NOT NULL      | INTEGER NOT NULL   
  存储       | storage    | TEXT NOT NULL        | LONGTEXT NOT NULL    | TEXT NOT NULL      
  创建时间   | created_at | TIMESTAMPTZ NOT NULL | DATETIME(3) NOT NULL | DATETIME NOT NULL  

主键: repo_id, oid
```

# 表 "login_source"

```
      字段      |    列名    |    POSTGRESQL    |         MYSQL         |     SQLITE3       
----------------+--------------+------------------+-----------------------+-------------------
  ID           | id           | BIGSERIAL        | BIGINT AUTO_INCREMENT | INTEGER           
  类型         | type         | BIGINT           | BIGINT                | INTEGER           
  名称         | name         | TEXT UNIQUE      | VARCHAR(191) UNIQUE   | TEXT UNIQUE       
  是否激活     | is_actived   | BOOLEAN NOT NULL | BOOLEAN NOT NULL      | NUMERIC NOT NULL  
  是否默认     | is_default   | BOOLEAN          | BOOLEAN               | NUMERIC           
  配置         | cfg          | TEXT             | TEXT                  | TEXT              
  创建时间Unix | created_unix | BIGINT           | BIGINT                | INTEGER           
  更新时间Unix | updated_unix | BIGINT           | BIGINT                | INTEGER           

主键: id
```

# 表 "notice"

```
      字段      |    列名    | POSTGRESQL |         MYSQL         | SQLITE3  
----------------+--------------+------------+-----------------------+----------
  ID           | id           | BIGSERIAL  | BIGINT AUTO_INCREMENT | INTEGER  
  类型         | type         | BIGINT     | BIGINT                | INTEGER  
  描述         | description  | TEXT       | TEXT                  | TEXT     
  创建时间Unix | created_unix | BIGINT     | BIGINT                | INTEGER  

主键: id
