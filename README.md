# Blog System Project

![Go](https://img.shields.io/badge/Go-1.24.2-blue) ![Vue](https://img.shields.io/badge/Vue-3-green) ![TypeScript](https://img.shields.io/badge/TypeScript-5.0-orange) ![Elasticsearch](https://img.shields.io/badge/Elasticsearch-8.17-yellow) ![MySQL](https://img.shields.io/badge/MySQL-8.0-blue) ![Redis](https://img.shields.io/badge/Redis-7.0-red)

## 项目结构
├── server/ # Go后端核心代码

│ ├── main # 编译后的Linux可执行文件

│ └── ... # 其他Go源码

├── web/ # Vue3前端

│ ├── src/ # 前端源码(TypeScript)

│ └── ...

├── session_create/ # 加密模块

│ └── session.go # 会话加密实现

└── dockerbuild_testing/ # Docker整合实验

└── Dockerfile.wip # 多服务镜像构建文件(开发中)


## 技术说明

### 后端服务 (server/)
- **Go语言**开发的后端逻辑
- 集成服务：
  - 🐬 MySQL - 主数据库
  - 🔍 Elasticsearch 8.17 - 全文搜索
  - 🗃️ Redis - 缓存管理
- 编译命令：
  ```bash
  cd server && go build -o main .
前端界面 (web/)
Vue3 + TypeScript开发

主要依赖：

缓存：redis

数据库：mysql、sql语法、参数化查询

认证：双token、jwt

查询：elasticsearch

UI框架：Element Plus

编辑器：MdeditorV3

开发运行：

bash
cd web && npm install && npm run dev

加密模块 (session_create/)




```go linenums="1"
// 基于SHA-256的会话令牌生成
func GenerateSecureToken(user string) string {
    h := sha256.New()
    h.Write([]byte(user + salt))
    return hex.EncodeToString(h.Sum(nil))
}
```
容器化部署 (实验阶段)

## 构建多服务镜像(开发中)
## 目标：部署到Claw Cloud免费服务

开发进度
✅ 已完成功能：

后端基础API

前端管理界面

会话加密模块

🚧 进行中：

多服务Docker镜像整合

Claw Cloud适配

压力测试


⚠️ 注意：dockerbuild_testing/ 下的为实验性代码
