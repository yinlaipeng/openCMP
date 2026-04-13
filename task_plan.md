# Task Plan: openCMP 项目落地实现

## Goal
实现 openCMP 多云管理平台的完整功能落地，前端页面与后端 API 真实对接，云厂商适配器调用真实 SDK，实现从云账号添加 → 资源同步 → 资源管理的完整业务流程。

## Current Phase
Phase 9 - 云账号完整流程实现

## Phases

### Phase 8: 后端 API 完善与云厂商适配器补全
- [x] **Azure 适配器实现 (已完成)**
  - Azure SDK 安装 (azure-sdk-for-go: azidentity, armcompute/v5, armnetwork/v4)
  - Azure Compute (VM 创建/删除/启停/列表) - 真实 SDK 调用
  - Azure Network (VNet/Subnet/SecurityGroup/PublicIP) - 真实 SDK 调用
  - 编译验证成功
- [ ] **阿里云 Database/Middleware 适配器**
  - RDS MySQL/PostgreSQL 支持
  - Redis 实例支持
- [ ] **腾讯云 Database 适配器**
  - TencentDB MySQL 支持
- [ ] **AWS Database 适配器**
  - RDS MySQL/PostgreSQL 支持
- **Status:** in_progress (Azure 完成，Database 待实现)

### Phase 9: 云账号完整流程实现
- [x] **云账号同步功能实现**
  - SyncResources 方法：同步所有资源类型
  - /sync API 端点已添加
  - 前端已有 sync 对话框，已正确对接
- [x] **凭证验证增强**
  - VerifyCredentials 方法：通过实际 API 调用验证
  - /verify-credentials API 端点已添加
  - 阿里云：验证区域列表；其他云：验证镜像列表
- [ ] **定时同步执行器**
  - 执行 scheduled_task 触发同步
  - 使用 sync_policy 规则映射资源
- [ ] **云账号状态管理完善**
  - 连接状态检测（定时）
  - 同步状态显示
  - 异常告警机制
- **Status:** in_progress

### Phase 10: 资源管理业务流程实现
- [ ] **虚拟机全生命周期**
  - 创建：前端参数 → 后端验证 → 云厂商 SDK 调用 → 状态跟踪
  - 操作：启动/停止/重启/删除 真实调用
  - 详情：实时获取云厂商数据 + 本地缓存
- [ ] **网络资源管理**
  - VPC 创建/删除 流程验证
  - Subnet 创建/删除 流程验证
  - SecurityGroup 规则配置
  - EIP 申请/释放/绑定
- [ ] **存储资源管理**
  - Block Storage 创建/挂载/卸载/删除
  - Object Storage Bucket 创建/管理
- [ ] **数据库资源管理**
  - RDS 实例创建/管理
  - Redis 实例创建/管理
- **Status:** pending

### Phase 11: 前后端联调与功能验证
- [ ] **核心功能联调**
  - 云账号添加 → 验证 → 同步 完整流程
  - VM 创建 → 状态跟踪 → 操作 完整流程
  - VPC/Subnet 创建 → 资源关联 完整流程
- [ ] **异常处理完善**
  - 云厂商 API 错误友好提示
  - 网络异常重试机制
  - 权限不足处理
- [ ] **UI 交互优化**
  - 加载状态反馈
  - 操作进度显示
  - 成功/失败消息提示
- **Status:** pending

### Phase 12: 项目交付准备
- [ ] **部署配置**
  - Docker 容器化
  - 环境变量配置
  - 数据库初始化脚本
- [ ] **文档完善**
  - API 文档更新
  - 部署指南
  - 用户手册
- [ ] **最终测试**
  - 功能完整性测试
  - 多云场景测试
  - 性能测试
- **Status:** pending

## Key Architecture Points

### 真实数据流
```
前端页面 → API 请求 → Handler → Service → CloudProvider Adapter → 云厂商 SDK → 真实云资源
                                ↓
                           数据库存储（账号/资源元数据/同步状态）
```

### 云厂商适配器状态
| 云厂商 | Compute | Network | Storage | Database | Middleware |
|--------|---------|---------|---------|----------|------------|
| 阿里云 | ✅ SDK  | ✅ SDK  | ✅ SDK  | ⚠️ 待实现 | ⚠️ 待实现 |
| 腾讯云 | ✅ SDK  | ✅ SDK  | ❌      | ❌       | ❌         |
| AWS    | ✅ SDK  | ✅ SDK  | ❌      | ❌       | ❌         |
| Azure  | ✅ SDK  | ✅ SDK  | ❌      | ❌       | ❌         |

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| 优先实现阿里云完整功能 | 国内主要用户群，SDK成熟，作为标杆实现 |
| Azure 作为后续扩展 | 国际需求相对较低，SDK较复杂 |
| 采用同步+缓存混合模式 | 云厂商 API 有延迟，本地缓存提升响应速度 |

## Notes
- 当前后端 service 层已经正确调用 cloudprovider.GetProvider() 获取适配器
- 前端 API 调用已准备好，只需确保后端返回真实数据
- 关键是完成云厂商适配器的真实 SDK 调用实现