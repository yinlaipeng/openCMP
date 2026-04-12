# openCMP 开发进度看板

> 由各 Agent 协作维护，完成一项更新一项。
> 图例：✅ 完成 | 🚧 进行中 | ⬜ 待开始

## IAM 模块

| 模块 | 设计 | 后端 | 前端 | 测试 | 状态 |
|------|------|------|------|------|------|
| 域 (Domain) | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 项目 (Project) | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 组 (Group) | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 用户 (User) | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 角色 (Role) | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 权限 (Permission) | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 认证源 (Auth Source) | ✅ | ✅ | ✅ | ✅ | 已完成（LDAP bind 待依赖库）|
| 策略 (Policy) | ✅ | ✅ | ✅ | ✅ | 已完成 |

## 平台功能模块

| 模块 | 设计 | 后端 | 前端 | 测试 | 状态 |
|------|------|------|------|------|------|
| 站内信 | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 消息订阅 | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 通知渠道 | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 接受人管理 | ✅ | ✅ | ✅ | ✅ | 已完成 |
| 机器人管理 | ✅ | ✅ | ✅ | ✅ | 已完成 |

## 多云管理模块

| 模块 | 设计 | 后端 | 前端 | 测试 | 状态 |
|------|------|------|------|------|------|
| 云账号增强 | ✅ | ✅ | ✅ | ⬜ | 待测试 |
| 同步策略 | ✅ | ✅ | ✅ | ⬜ | 前端已完成，待测试 |
| 定时同步任务 | ✅ | ✅ | ✅ | ⬜ | 待测试 |

**详细说明：**
- **同步策略**：全部完成（handler/service/model + API + 前端页面）
- **定时同步任务**：全部完成（handler/service/model + API + 前端页面）
- **资源同步规则**：已整合到同步策略的 Rule 模型中

## 云厂商适配器

| 云厂商 | 主机 | 网络 | 存储 | 数据库 | 中间件 | 状态 |
|--------|------|------|------|--------|--------|------|
| 阿里云 | ✅ | ✅ | ✅ | ⬜ | ⬜ | 核心完成 |
| 腾讯云 | ✅ | ✅ | ⬜ | ⬜ | ⬜ | 核心完成 |
| AWS | ✅ | ✅ | ⬜ | ⬜ | ⬜ | 核心完成 |
| Azure | 🚧 | 🚧 | ⬜ | ⬜ | ⬜ | 骨架 |

## 最近更新

- 2026-04-12: 前端云资源创建弹窗全部完成（CreateVMModal 5步向导、CreateVPCModal、CreateSubnetModal、CloudAccountSelector、CIDR 校验工具）
- 2026-04-12: AWS 适配器实现完成（Compute: VM 镜像管理 + Network: VPC/Subnet/SecurityGroup/EIP）
- 2026-04-12: 腾讯云适配器实现完成（Compute: VM 镜像管理 + Network: VPC/Subnet/SecurityGroup/EIP）
- 2026-04-12: 同步策略前端页面补全完成 2026-04-05: 消息中心测试用例全部完成（message/notification_channel/robot/receiver/subscription 服务层和 handler 层测试覆盖）
- 2026-04-05: 认证源模块完善：操作列改为"更多"下拉、LDAP 认证流程骨架、域绑定登录、JWT 携带 domain_id、service/handler 双层单元测试（14 个测试用例）
- 2026-04-05: 产品设计 Agent 规范文档创建（docs/superpowers/specs/2026-04-05-product-design-agent-guide.md）
- 2026-04-05: 消息中心后端 handler 层全部完成（message/notification_channel/robot/receiver/subscription 5 个 handler + 路由注册）
- 2026-04-05: 消息中心前端全部完成（api/message.ts + 4个新页面 + 站内信完善 + 路由 + 侧边栏）
- 2026-04-05: 修复 PolicyStatement 未添加到 AutoMigrate 的问题
- 2026-04-03: 修复后端 Go 编译错误（PolicyStatement 模型、GroupID 问题、PolicyID 类型、scripts main 重复、context 导入）
- 2026-04-03: IAM 模块核心功能已全部完成（域、项目、用户、组、角色、权限、认证源、策略）