#!/bin/bash

# # 停止并删除旧的容器
# docker stop d8d-gogs
# docker rm d8d-gogs

# # 删除旧的 Docker 镜像
# docker rmi registry.cn-beijing.aliyuncs.com/d8dcloud/d8d-gogs:latest

# 定义变量
GOGS_DIR="/gogs"

# 创建临时目录
temp_dir=$(mktemp -d)
cd "$temp_dir"

# 将 app.prd.ini 复制到临时目录
cp "/docker/codeserver/project/test/gogs/app.prd.ini" .

# 将 gogs 文件 复制到临时目录
cp "/docker/codeserver/project/test/gogs/gogs" .

# 编写 Dockerfile
cat << EOF > Dockerfile
# 基于你现有的 debian 镜像
FROM debian:latest

# 创建并写入阿里云镜像源配置
RUN echo "deb http://mirrors.aliyun.com/debian/ bookworm main non-free contrib" > /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian/ bookworm main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian-security/ bookworm-security main" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian-security/ bookworm-security main" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian/ bookworm-updates main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian/ bookworm-updates main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian/ bookworm-backports main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian/ bookworm-backports main non-free contrib" >> /etc/apt/sources.list

# 创建git 用户
RUN useradd -m -s /bin/bash git

# 切换到 git 用户
USER git

# 安装必要的依赖
RUN apt-get update && apt-get install -y curl wget unzip git

RUN echo 'alias ll="ls -alF"' >> ~/.bashrc


# 清除 apt 缓存和无用文件
RUN apt-get clean && rm -rf /var/lib/apt/lists/*


# 创建必要的目录
RUN mkdir -p ${GOGS_DIR}
RUN mkdir -p ${GOGS_DIR}/custom/conf

# 复制 app.ini 包
COPY app.prd.ini ${GOGS_DIR}/custom/conf/app.ini

# 复制 gogs 包
COPY gogs ${GOGS_DIR}/

# 授予执行权限
RUN chmod +x ${GOGS_DIR}/gogs

# 暴露 Gogs 端口
EXPOSE 80 443 22

# 启动 Gogs
CMD ["${GOGS_DIR}/gogs", "web"]
EOF


# 构建新的 Docker 镜像
docker build -t registry.cn-beijing.aliyuncs.com/d8dcloud/d8d-gogs:latest .

# 运行新的容器
# docker run -d -p 23962:23958 \
#   -v /tmp/gogs-repositories:/mnt/gogs/gogs-repositories \
#   -v /tmp/gogs-log:/mnt/gogs/log \
#   registry.cn-beijing.aliyuncs.com/d8dcloud/d8d-gogs:latest

# docker push registry.cn-beijing.aliyuncs.com/d8dcloud/d8d-gogs:latest

# 删除临时目录
rm -rf "$temp_dir"