<template>
  <div class="finance-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">预算管理</span>
          <el-button type="primary">新建预算</el-button>
        </div>
      </template>
      <el-empty description="预算管理功能开发中..." />
    </el-card>
  </div>
</template>

<style scoped>
.finance-page {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
</style>