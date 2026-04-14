<template>
  <el-dialog
    v-model="visible"
    :title="`VNC 远程桌面 - ${vmName || '虚拟机'}`"
    width="80%"
    :destroy-on-close="true"
    :close-on-click-modal="false"
  >
    <div class="vnc-container" v-loading="connecting">
      <div v-if="showConnectionInfo" class="connection-info">
        <el-alert
          title="连接信息"
          type="info"
          :closable="false"
          show-icon
        >
          <p>正在建立远程连接，请稍候...</p>
          <p v-if="vncInfo.console_url">控制台地址: {{ vncInfo.console_url }}</p>
          <p>虚拟机: {{ vmName }}</p>
        </el-alert>
      </div>

      <div v-if="hasConnection" class="vnc-controls">
        <el-button-group>
          <el-button size="small" @click="sendCtrlAltDel" :disabled="!connected">
            <el-icon><Key /></el-icon>
            发送 Ctrl+Alt+Del
          </el-button>
          <el-button size="small" @click="toggleScale" :disabled="!connected">
            <el-icon><FullScreen /></el-icon>
            {{ scaled ? '取消缩放' : '适应屏幕' }}
          </el-button>
          <el-button size="small" @click="disconnect" :disabled="!connected">
            <el-icon><SwitchButton /></el-icon>
            断开连接
          </el-button>
        </el-button-group>
      </div>

      <!-- VNC 显示容器 -->
      <div v-if="!vncError" ref="vncContainerRef" class="vnc-display"></div>

      <div v-if="vncError" class="vnc-error">
        <el-result icon="error" title="连接失败" :sub-title="vncError">
          <template #extra>
            <el-button type="primary" @click="retryConnection">重试连接</el-button>
            <el-button @click="close">关闭</el-button>
          </template>
        </el-result>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">关闭</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Key, FullScreen, SwitchButton } from '@element-plus/icons-vue'
import { getVNCInfo } from '@/api/compute'

interface Props {
  modelValue: boolean
  vmId: string
  vmName: string
  accountId: number
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const connecting = ref(true)
const connected = ref(false)
const scaled = ref(false)
const vncError = ref('')
const vncInfo = ref<any>({})
const vncContainerRef = ref<HTMLDivElement>()
let rfb: any = null

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value)
})

const showConnectionInfo = computed(() => connecting.value && !vncError.value)
const hasConnection = computed(() => connected.value && !vncError.value)

// 方法
const close = () => {
  disconnect()
  visible.value = false
  emit('close')
}

const disconnect = () => {
  if (rfb) {
    rfb.disconnect()
    connected.value = false
    connecting.value = false
  }
}

const toggleScale = () => {
  scaled.value = !scaled.value
  if (rfb) {
    rfb.scaleViewport(scaled.value)
    rfb.fitWindow(scaled.value)
  }
}

const sendCtrlAltDel = () => {
  if (rfb && connected.value) {
    rfb.sendCtrlAltDel()
    ElMessage.success('已发送 Ctrl+Alt+Del')
  }
}

const loadVNCInfo = async () => {
  try {
    vncInfo.value = await getVNCInfo(props.vmId, props.accountId)
    // 初始化 VNC 连接
    await initVNC()
  } catch (error) {
    console.error('Failed to get VNC info:', error)
    vncError.value = error instanceof Error ? error.message : '获取VNC连接信息失败'
    connecting.value = false
  }
}

const initVNC = async () => {
  // 在实际应用中，这里应该集成 noVNC 库
  // 由于前端依赖问题，这里先模拟连接过程
  connecting.value = true

  // 模拟连接延迟
  setTimeout(() => {
    connecting.value = false
    connected.value = true
    vncError.value = ''

    // 在真实场景中，会在这里初始化 noVNC 连接
    // 例如：
    /*
    import RFB from '@novnc/novnc/core/rfb'
    const vncUrl = vncInfo.value.ws_url || `ws://${window.location.hostname}:6080/websockify?token=${props.vmId}`
    rfb = new RFB(vncContainerRef.value!, vncUrl, {
      credentials: { password: vncInfo.value.password || '' }
    })

    rfb.addEventListener('connect', () => {
      connected.value = true
      connecting.value = false
    })

    rfb.addEventListener('disconnect', (e) => {
      connected.value = false
      if (e.detail.clean) {
        ElMessage.info('已断开连接')
      } else {
        vncError.value = '连接意外断开'
      }
    })

    rfb.addEventListener('credentialsrequired', (e) => {
      // 如果需要凭据，则处理它们
    })

    rfb.addEventListener('securityfailure', (e) => {
      vncError.value = e.detail.reason || '安全验证失败'
      connecting.value = false
    })
    */

    ElMessage.success('VNC 连接已建立 (模拟)')
  }, 2000)
}

const retryConnection = () => {
  vncError.value = ''
  loadVNCInfo()
}

// 监听可见性变化
watch(visible, async (newValue) => {
  if (newValue && props.vmId && props.accountId) {
    await loadVNCInfo()
  } else if (!newValue) {
    disconnect()
  }
})

onUnmounted(() => {
  disconnect()
})
</script>

<style scoped>
.vnc-container {
  min-height: 500px;
}

.connection-info {
  margin-bottom: 20px;
}

.vnc-controls {
  margin-bottom: 20px;
  text-align: center;
}

.vnc-display {
  width: 100%;
  height: 500px;
  border: 1px solid #dcdfe6;
  background-color: #000;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.vnc-display::before {
  content: 'VNC 控制台显示区域';
  color: #999;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.vnc-error {
  text-align: center;
  padding: 40px 0;
}

.dialog-footer {
  text-align: right;
}
</style>