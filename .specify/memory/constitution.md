# openCMP Constitution

> openCMP 是开源多云管理平台，通过统一界面和适配器模式实现多云厂商（阿里云、腾讯云、AWS、Azure）资源的统一管理和快速访问。核心模块包括 IAM 身份管理、云资源管理（主机/网络/存储/数据库）、消息中心、多云同步管理，支持租户隔离（域）、资源分组（项目）、细粒度权限控制（RBAC）。

---

## Context Detection

**Ralph Loop Mode** (started by ralph-loop*.sh):
- Pick highest priority incomplete spec from `specs/`
- Implement, test, commit, push
- Output `<promise>DONE</promise>` only when 100% complete
- Output `<promise>ALL_DONE</promise>` when no work remains

**Interactive Mode** (normal conversation):
- Be helpful, guide decisions, create specs

---

## Core Principles

1. **真实API对接** - 确保后端API真实调用云厂商SDK，前端真实对接后端API，不使用Mock数据
2. **模块化架构** - 使用标准adapter模式（ICompute/INetwork/IStorage/IDatabase），支持多云厂商扩展
3. **功能完整性优先** - 先完成核心流程（云账号添加→资源同步→资源管理），再迭代完善细节
4. **UI一致性** - 页面样式统一，遵循 ui-ux-pro-max 规范（spacing-scale 8dp、consistent structure）

---

## Technical Stack

**Backend**: Go 1.21 + Gin 1.9 + GORM 1.25
**Frontend**: Vue 3 + TypeScript + Element Plus
**Database**: MySQL 8.0 / PostgreSQL
**Cache**: Redis 7.x
**Scheduler**: robfig/cron
**Cloud SDKs**: 
- Alibaba: aliyun-go-sdk
- Tencent: tencentcloud-sdk-go
- AWS: aws-sdk-go-v2
- Azure: azure-sdk-for-go

**Project Structure**:
- `backend/cmd/server/main.go` - Entry point
- `backend/internal/handler/` - HTTP handlers
- `backend/internal/service/` - Business logic
- `backend/internal/model/` - Database models
- `backend/pkg/cloudprovider/` - Cloud adapter interfaces
- `frontend/src/views/` - Vue pages
- `frontend/src/api/` - API clients

---

## Autonomy

YOLO Mode: ENABLED (execute commands, modify files, run tests without asking each time)
Git Autonomy: ENABLED (commit and push automatically after each spec completion)

---

## Specs

Specs live in `specs/` as markdown files. Pick the highest priority incomplete spec (lower number = higher priority). A spec is incomplete if it lacks `## Status: COMPLETE`.

Spec template: https://raw.githubusercontent.com/github/spec-kit/refs/heads/main/templates/spec-template.md

---

## Project-Specific Rules

### Backend Development
- Always use real cloud SDK calls in adapters (no mock returns)
- Services should query local database (CloudVM/CloudVPC etc.) for resource lists
- PermissionMiddleware and ProjectIsolationMiddleware are registered in main.go
- Resource sync: Incremental (INSERT new) vs Full (INSERT new, UPDATE existing, MARK terminated)

### Frontend Development
- Use CloudAccountSelector component for account selection (not manual ID input)
- Page structure: `.xxx-container { padding: 20px }`, `.page-header`, `.filter-card`, `.pagination`
- Table must have `row-key="id"`
- Follow host-templates/index.vue as reference style

### Testing & Verification
- Run `go test -v ./...` for backend tests
- Run `npx vite build` for frontend build verification
- Check git status after changes

---

## Ralph Wiggum Version

Commit: d205125cc33745116cce22d883417461174dcde5
Source: https://github.com/fstandhartinger/ralph-wiggum