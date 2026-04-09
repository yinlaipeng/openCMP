<template>
  <div class="rule-configuration">
    <div v-for="(rule, index) in localRules" :key="rule._id || index" class="rule-item">
      <div class="rule-header">
        <span class="rule-title">规则 {{ index + 1 }}</span>
        <el-button size="small" type="danger" @click="removeRule(index)" plain>删除</el-button>
      </div>

      <div class="rule-content">
        <div class="form-row">
          <label class="input-label">条件类型</label>
          <el-select v-model="rule.condition_type" placeholder="选择条件类型" class="full-width">
            <el-option label="全部匹配标签" value="all_match" />
            <el-option label="至少一个标签" value="any_match" />
            <el-option label="根据标签key匹配" value="key_match" />
          </el-select>
        </div>

        <div class="form-row">
          <label class="input-label">资源映射</label>
          <el-select v-model="rule.resource_mapping" placeholder="选择资源映射方式" class="full-width">
            <el-option label="指定项目" value="specify_project" />
            <el-option label="指定名称" value="specify_name" />
          </el-select>
        </div>

        <div class="form-row" v-if="rule.resource_mapping === 'specify_project'">
          <label class="input-label">目标项目</label>
          <el-select v-model="rule.target_project_id" placeholder="选择项目" class="full-width">
            <el-option
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
          <div class="field-note">指定项目时：当资源符合上述条件时，同步云账号时资源会归属到该项目</div>
        </div>

        <div class="form-row" v-if="rule.resource_mapping === 'specify_name'">
          <label class="input-label">项目名称</label>
          <el-input v-model="rule.target_project_name" placeholder="请输入项目名称" class="full-width" />
          <div class="field-note">指定项目名称时：根据云账号归属域设置项目，已有同名项目直接使用，否则自动创建</div>
        </div>

        <!-- 标签配置 -->
        <div class="form-row">
          <label class="input-label">标签配置</label>
          <div class="tags-container">
            <el-tag
              v-for="(tag, tagIndex) in rule.tags || []"
              :key="tagIndex"
              closable
              @close="removeTag(index, tagIndex)"
              class="tag-item"
            >
              {{ tag.key }}: {{ tag.value }}
            </el-tag>

            <el-popover
              v-model:visible="tagFormVisible[index]"
              placement="bottom"
              width="450"
              trigger="click"
            >
              <div class="tag-form">
                <div class="tag-input-row">
                  <el-input v-model="tempTag.key" placeholder="标签键" class="tag-input-half" />
                  <el-input v-model="tempTag.value" placeholder="标签值" class="tag-input-half" />
                </div>
                <div class="tag-form-buttons">
                  <el-button
                    size="small"
                    type="primary"
                    @click="addTag(index)"
                    :disabled="!tempTag.key || !tempTag.value"
                  >
                    添加
                  </el-button>
                  <el-button size="small" @click="cancelAddTag(index)">取消</el-button>
                </div>
              </div>
              <template #reference>
                <el-button size="small" type="primary" plain>+ 添加标签</el-button>
              </template>
            </el-popover>
          </div>
        </div>
      </div>
    </div>

    <el-button @click="addRule" type="primary" size="large" class="add-rule-btn">
      <el-icon><Plus /></el-icon>
      添加规则
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { Rule } from '@/types/sync-policy'
import { getProjects } from '@/api/project'
import type { Project } from '@/types'

interface Props {
  modelValue: Rule[]
}

interface Emits {
  (e: 'update:modelValue', value: Rule[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 生成唯一ID帮助器
let idCounter = 0
const generateId = () => ++idCounter

// 本地规则副本
const localRules = ref<Rule[]>(props.modelValue.map(rule => ({
  ...rule,
  tags: [...(rule.tags || [])],
  _id: generateId() // 添加唯一ID避免key重复
})))

const projects = ref<Project[]>([])
const tagFormVisible = ref<boolean[]>([])
const tempTag = reactive({
  key: '',
  value: ''
})

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  // 只有在值真正变化时才更新本地副本
  if (JSON.stringify(newValue) !== JSON.stringify(localRules.value.map(r => ({...r, _id: undefined})))) {
    localRules.value = newValue.map(rule => ({
      ...rule,
      tags: [...(rule.tags || [])],
      _id: generateId()
    }))
    // 重新初始化标签表单可见性
    tagFormVisible.value = Array(localRules.value.length).fill(false)
  }
}, { deep: true, immediate: true })

// 监听本地规则变化并同步到父组件
watch(localRules, (newRules) => {
  // 移除内部使用的 _id 字段后再发送
  const rulesToSend = newRules.map(({ _id, ...rest }) => rest as Rule)
  emit('update:modelValue', rulesToSend)
}, { deep: true })

// 加载项目列表
const loadProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    console.error('Failed to load projects:', e)
  }
}

// 添加规则
const addRule = () => {
  const newRule: Rule = {
    condition_type: 'all_match',
    resource_mapping: 'specify_project',
    tags: [],
    _id: generateId()
  }
  localRules.value.push(newRule)
  tagFormVisible.value.push(false)
}

// 删除规则
const removeRule = (index: number) => {
  localRules.value.splice(index, 1)
  tagFormVisible.value.splice(index, 1)
}

// 添加标签到规则
const addTag = (ruleIndex: number) => {
  if (!tempTag.key || !tempTag.value) return

  if (!localRules.value[ruleIndex].tags) {
    localRules.value[ruleIndex].tags = []
  }

  localRules.value[ruleIndex].tags!.push({
    key: tempTag.key,
    value: tempTag.value
  })

  // 清空临时标签
  tempTag.key = ''
  tempTag.value = ''

  // 关闭表单
  tagFormVisible.value[ruleIndex] = false
}

// 取消添加标签
const cancelAddTag = (ruleIndex: number) => {
  tempTag.key = ''
  tempTag.value = ''
  tagFormVisible.value[ruleIndex] = false
}

// 删除标签
const removeTag = (ruleIndex: number, tagIndex: number) => {
  if (localRules.value[ruleIndex].tags) {
    localRules.value[ruleIndex].tags!.splice(tagIndex, 1)
  }
}

onMounted(() => {
  loadProjects()
  // 初始化标签表单可见性
  tagFormVisible.value = Array(localRules.value.length).fill(false)
})
</script>

<style scoped>
.rule-configuration {
  padding: 20px 0;
}

.rule-item {
  margin-bottom: 30px;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
  background-color: white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.rule-header {
  padding: 15px 20px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rule-title {
  font-weight: 600;
  color: #303133;
  font-size: 16px;
}

.rule-content {
  padding: 25px;
}

.form-row {
  margin-bottom: 25px;
}

.input-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
  font-size: 14px;
}

.field-note {
  margin-top: 8px;
  padding: 10px;
  background-color: #f0f9ff;
  border: 1px solid #e1f0ff;
  border-radius: 4px;
  color: #5d7a96;
  font-size: 13px;
  line-height: 1.5;
}

.full-width {
  width: 100%;
  min-width: 300px; /* 确保最小宽度 */
}

.tags-container {
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background-color: #fafcff;
  min-height: 70px;
}

.tag-item {
  margin-right: 10px;
  margin-bottom: 10px;
}

.tag-form {
  padding: 15px 0;
}

.tag-input-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.tag-input-half {
  flex: 1;
}

.tag-form-buttons {
  text-align: right;
}

.add-rule-btn {
  width: 100%;
  margin-top: 15px;
  height: 48px;
  font-size: 16px;
}
</style>