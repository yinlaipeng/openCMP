# 云资源创建弹窗设计文档

> 日期: 2026-04-12
> 模块: 前端云资源页面完善
> 状态: 待实现

## 概述

为 openCMP 多云管理平台添加云资源创建功能，包括：
- 创建虚拟机弹窗（分步向导）
- 创建 VPC 弹窗（单页表单）
- 创建 Subnet 弹窗（单页表单）

## 1. 创建虚拟机弹窗 (CreateVMModal.vue)

### 1.1 组件位置
`frontend/src/components/vm/CreateVMModal.vue`

### 1.2 组件结构
```
CreateVMModal.vue
├── Steps 导航条（显示当前步骤）
├── 步骤内容区域
│   ├── Step 1: 基本配置
│   ├── Step 2: 计算配置
│   ├── Step 3: 网络配置
│   ├── Step 4: 存储配置
│   └── Step 5: 确认与创建
├── 底部按钮（上一步/下一步/取消/创建）
```

### 1.3 Props 和 Emits

```typescript
interface Props {
  visible: boolean      // 弹窗显示控制 (v-model)
  accountId?: number    // 可预设云账号
  templateId?: string   // 可预设模板 ID
}

interface Emits {
  'update:visible': (val: boolean) => void  // 关闭弹窗
  'success': (vm: VirtualMachine) => void   // 创建成功回调
}
```

### 1.4 步骤详细内容

#### Step 1: 基本配置

| 字段 | 类型 | 说明 |
|-----|------|------|
| 云账号选择 | Select | 必选，从已配置的 CloudAccount 列表选择 |
| 创建方式 | Radio | **使用模板** / **自定义配置** |
| 主机模板 | Select | 仅当选择"使用模板"时显示 |
| 虚拟机名称 | Input | 必填，1-64字符 |
| 数量 | InputNumber | 创建数量，默认 1，范围 1-100 |

**联动逻辑**：
- 选择云账号后 → 加载 Regions、Images、VPCs 等数据，清空后续依赖字段
- 选择模板后 → 自动填充 Step 2-4 配置（用户仍可调整）

#### Step 2: 计算配置

| 字段 | 类型 | 说明 |
|-----|------|------|
| 区域 | Select | 必选，从云账号可用区域列表选择 |
| 可用区 | Select | 必选，依赖区域选择 |
| 镜像 | Select | 必选，系统镜像列表，支持搜索过滤 |
| 实例规格 | Select | 必选，CPU/内存规格组合 |
| 密钥对 | Select | 可选，SSH 登录密钥 |

#### Step 3: 网络配置

| 字段 | 类型 | 说明 |
|-----|------|------|
| VPC | Select | 必选，从已有 VPC 选择 |
| Subnet | Select | 必选，依赖 VPC 选择 |
| 安全组 | MultiSelect | 至少选择一个安全组 |
| 公网 IP | Switch | 是否分配公网 IP |
| 带宽 | InputNumber | 仅当公网 IP 开启时显示，Mbps |

#### Step 4: 存储配置

| 字段 | 类型 | 说明 |
|-----|------|------|
| 系统盘大小 | InputNumber | 默认镜像推荐大小，可调整 |
| 系统盘类型 | Select | cloud_ssd / cloud_efficiency 等 |
| 数据盘 | 动态列表 | 可添加/删除多块数据盘 |

数据盘项结构：
```typescript
{
  size: number    // GB
  type: string    // 磁盘类型
}
```

#### Step 5: 确认与创建

- 显示所有配置参数汇总表格
- 显示操作提示（创建时间预估）
- 确认按钮提交创建请求

### 1.5 内部状态

```typescript
interface VMCreateState {
  // 基本配置
  accountId: number | null
  createMode: 'template' | 'custom'
  templateId: string | null
  name: string
  count: number
  
  // 计算配置
  regionId: string
  zoneId: string
  imageId: string
  instanceType: string
  keypairId: string
  
  // 网络配置
  vpcId: string
  subnetId: string
  securityGroups: string[]
  enablePublicIp: boolean
  bandwidth: number
  
  // 存储配置
  systemDiskSize: number
  systemDiskType: string
  dataDisks: Array<{ size: number; type: string }>
}
```

### 1.6 API 调用时机

| API | 调用时机 |
|-----|---------|
| `getCloudAccounts()` | 弹窗打开时 |
| `getHostTemplates()` | Step 1 选择"使用模板"时 |
| `getRegions(accountId)` | Step 1 选择云账号后 |
| `getZones(accountId, regionId)` | Step 2 选择区域后 |
| `getImages(accountId)` | 进入 Step 2 时 |
| `getVPCs(accountId, regionId)` | 进入 Step 3 时 |
| `getSubnets(accountId, vpcId)` | Step 3 选择 VPC 后 |
| `getSecurityGroups(accountId, vpcId)` | Step 3 选择 VPC 后 |
| `createVM(formData)` | Step 5 点击确认 |

### 1.7 错误处理

| 场景 | 处理方式 |
|-----|---------|
| 云账号连接失败 | 显示错误提示，允许重新选择 |
| API 加载失败 | 对应字段显示"加载失败"状态，可点击重试 |
| 创建请求失败 | 显示具体错误信息，不关闭弹窗 |
| 必填字段未填 | 禁用"下一步"按钮，显示提示 |

### 1.8 创建成功后

- 显示成功提示（ElMessage.success）
- 关闭弹窗
- 触发 `@success` 事件
- 父组件刷新 VM 列表

---

## 2. 创建 VPC 弹窗 (CreateVPCModal.vue)

### 2.1 组件位置
`frontend/src/components/network/CreateVPCModal.vue`

### 2.2 表单结构（单页）

| 字段 | 类型 | 说明 |
|-----|------|------|
| 云账号 | Select | 必选 |
| 名称 | Input | 必填，1-64字符 |
| IPv4 CIDR | Input | 必填，格式校验（如 10.0.0.0/16） |
| IPv6 CIDR | Input | 可选 |
| 描述 | TextArea | 可选，最多 256字符 |

### 2.3 Props 和 Emits

```typescript
interface Props {
  visible: boolean
  accountId?: number
}

interface Emits {
  'update:visible': (val: boolean) => void
  'success': (vpc: VPC) => void
}
```

### 2.4 校验规则

- IPv4 CIDR：正则校验 `/^\d{1,3}(\.\d{1,3}){3}\/\d{1,2}$/`
- 名称：不允许特殊字符

---

## 3. 创建 Subnet 弹窗 (CreateSubnetModal.vue)

### 3.1 组件位置
`frontend/src/components/network/CreateSubnetModal.vue`

### 3.2 表单结构（单页）

| 字段 | 类型 | 说明 |
|-----|------|------|
| 云账号 | Select | 必选 |
| VPC | Select | 必选，依赖云账号 |
| 名称 | Input | 必填，1-64字符 |
| CIDR | Input | 必填，需在 VPC CIDR 范围内 |
| 可用区 | Select | 必选 |
| 描述 | TextArea | 可选 |

### 3.3 校验规则

- Subnet CIDR 必须属于所选 VPC 的 CIDR 范围
- 需校验 CIDR 格式有效性

---

## 4. 公共组件与工具函数

### 4.1 CloudAccountSelector 组件
可复用的云账号选择器组件：
- Props: `value`, `disabled`
- Emits: `change`
- 功能: 加载云账号列表，显示状态标识

### 4.2 CIDR 校验函数
```typescript
// utils/cidr.ts
export function validateCIDR(cidr: string): boolean
export function isSubnetInVPC(subnetCIDR: string, vpcCIDR: string): boolean
```

---

## 5. 文件清单

| 文件 | 类型 | 说明 |
|-----|------|------|
| `components/vm/CreateVMModal.vue` | 新增 | 创建 VM 分步弹窗 |
| `components/network/CreateVPCModal.vue` | 新增 | 创建 VPC 弹窗 |
| `components/network/CreateSubnetModal.vue` | 新增 | 创建 Subnet 弹窗 |
| `components/common/CloudAccountSelector.vue` | 新增 | 云账号选择器 |
| `utils/cidr.ts` | 新增 | CIDR 校验工具 |
| `views/compute/vms/index.vue` | 修改 | 引入 CreateVMModal |
| `views/network/vpcs/index.vue` | 修改 | 添加创建按钮，引入 CreateVPCModal |
| `views/network/subnets/index.vue` | 修改 | 添加创建按钮，引入 CreateSubnetModal |

---

## 6. 实现优先级

1. **P0 - 核心**: CreateVMModal.vue（分步向导）
2. **P1 - 重要**: CreateVPCModal.vue、CreateSubnetModal.vue
3. **P2 - 优化**: CloudAccountSelector 组件复用、CIDR 校验工具

---

## 7. 测试要点

- 分步导航：前进、后退、跳转限制
- 字段联动：云账号切换清空依赖数据
- 表单校验：必填项、格式校验、范围校验
- API 错误处理：加载失败、创建失败
- 成功回调：弹窗关闭、列表刷新