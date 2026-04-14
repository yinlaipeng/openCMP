<template>
  <el-config-provider :locale="locale">
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import zhCn from 'element-plus/es/locale/lang/zh-cn'
const locale = zhCn
</script>

<style>
/* Import Fira fonts from Google Fonts */
@import url('https://fonts.googleapis.com/css2?family=Fira+Code:wght@400;500;600;700&family=Fira+Sans:wght@300;400;500;600;700&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: var(--font-sans);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100vh;
  color: var(--color-foreground);
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--color-background);
}

/* Ensure Element Plus uses our font stack */
.el-button,
.el-input__inner,
.el-select__wrapper,
.el-form-item__label,
.el-table,
.el-dialog__title {
  font-family: var(--font-sans);
}

/* Technical data elements use Fira Code */
.el-table .el-table__cell[data-type="ip"],
.el-table .el-table__cell[data-type="cidr"],
.el-table .el-table__cell[data-type="id"],
.el-descriptions__content,
.el-input-number__wrapper .el-input__inner {
  font-family: var(--font-mono);
}
</style>
