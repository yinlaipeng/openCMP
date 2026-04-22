import request from '@/utils/request'
import type { VirtualMachine, Image, HostTemplate } from '@/types'

export function getVMs(params?: { account_id?: number; vpc_id?: string; region_id?: string }) {
  return request<VirtualMachine[]>({
    url: '/compute/vms',
    method: 'get',
    params
  })
}

export function getVM(id: string, account_id: number) {
  return request<VirtualMachine>({
    url: `/compute/vms/${id}`,
    method: 'get',
    params: { account_id }
  })
}

export function createVM(data: {
  account_id: number
  name: string
  instance_type: string
  image_id: string
  vpc_id: string
  subnet_id: string
  security_groups?: string[]
  disk_size?: number
}) {
  return request<VirtualMachine>({
    url: '/compute/vms',
    method: 'post',
    data
  })
}

export function deleteVM(id: string, account_id: number) {
  return request({
    url: `/compute/vms/${id}`,
    method: 'delete',
    params: { account_id }
  })
}

export function vmAction(id: string, account_id: number, action: 'start' | 'stop' | 'reboot') {
  return request({
    url: `/compute/vms/${id}/action`,
    method: 'post',
    params: { account_id },
    data: { action }
  })
}

export function getImages(params?: { account_id?: number; platform?: string }) {
  return request<Image[]>({
    url: '/images',
    method: 'get',
    params
  })
}

// 获取虚拟机详细信息
export function getVMDetails(id: string, account_id: number) {
  return request<VirtualMachine>({
    url: `/compute/vms/${id}/details`,
    method: 'get',
    params: { account_id }
  })
}

// 获取虚拟机关联的安全组
export function getVMSecurityGroups(id: string, account_id: number) {
  return request<any>({
    url: `/compute/vms/${id}/security-groups`,
    method: 'get',
    params: { account_id }
  })
}

// 获取虚拟机网络信息
export function getVMNetworkInfo(id: string, account_id: number) {
  return request<any>({
    url: `/compute/vms/${id}/networks`,
    method: 'get',
    params: { account_id }
  })
}

// 获取虚拟机关联的磁盘
export function getVMDiskInfo(id: string, account_id: number) {
  return request<any>({
    url: `/compute/vms/${id}/disks`,
    method: 'get',
    params: { account_id }
  })
}

// 获取虚拟机相关的快照
export function getVMSnapshots(id: string, account_id: number) {
  return request<any>({
    url: `/compute/vms/${id}/snapshots`,
    method: 'get',
    params: { account_id }
  })
}

// 获取虚拟机操作日志
export function getVMOperationLogs(id: string, account_id: number) {
  return request<any>({
    url: `/compute/vms/${id}/logs`,
    method: 'get',
    params: { account_id }
  })
}

// 获取VNC连接信息
export function getVNCInfo(id: string, account_id: number) {
  return request<any>({
    url: `/compute/vms/${id}/vnc`,
    method: 'get',
    params: { account_id }
  })
}

// 重置密码
export function resetPassword(id: string, account_id: number, data: { username?: string, new_password?: string }) {
  return request({
    url: `/compute/vms/${id}/action`,
    method: 'post',
    params: { account_id },
    data: { ...data, action: 'reset_password' }
  })
}

// 更新虚拟机配置
export function updateVMConfig(id: string, account_id: number, data: { instance_type?: string, name?: string }) {
  return request({
    url: `/compute/vms/${id}/action`,
    method: 'post',
    params: { account_id },
    data: { ...data, action: 'update_config' }
  })
}

// Host Template API functions
export function getHostTemplates(params?: { project_id?: string; page?: number; page_size?: number }) {
  return request<{ items: HostTemplate[]; pagination: any }>({
    url: '/host-templates',
    method: 'get',
    params
  })
}

export function getHostTemplate(id: string) {
  return request<HostTemplate>({
    url: `/host-templates/${id}`,
    method: 'get'
  })
}

export function createHostTemplate(data: Partial<HostTemplate>) {
  return request<HostTemplate>({
    url: '/host-templates',
    method: 'post',
    data
  })
}

export function updateHostTemplate(id: string, data: Partial<HostTemplate>) {
  return request<HostTemplate>({
    url: `/host-templates/${id}`,
    method: 'put',
    data
  })
}

export function deleteHostTemplate(id: string) {
  return request({
    url: `/host-templates/${id}`,
    method: 'delete'
  })
}

// Autoscaling Group API functions
export interface AutoscalingGroup {
  id: string
  name: string
  description: string
  status: string
  host_template_id: string
  current_capacity: number
  desired_capacity: number
  min_size: number
  max_size: number
  platform: string
  project_id: string
  region_id: string
  zone_id: string
  tags: string
  created_at: string
  updated_at: string
}

export function getAutoscalingGroups(params?: { project_id?: string; page?: number; page_size?: number }) {
  return request<{ items: AutoscalingGroup[]; pagination: any }>({
    url: '/autoscaling-groups',
    method: 'get',
    params
  })
}

export function getAutoscalingGroup(id: string) {
  return request<AutoscalingGroup>({
    url: `/autoscaling-groups/${id}`,
    method: 'get'
  })
}

export function createAutoscalingGroup(data: Partial<AutoscalingGroup>) {
  return request<AutoscalingGroup>({
    url: '/autoscaling-groups',
    method: 'post',
    data
  })
}

export function updateAutoscalingGroup(id: string, data: Partial<AutoscalingGroup>) {
  return request<AutoscalingGroup>({
    url: `/autoscaling-groups/${id}`,
    method: 'put',
    data
  })
}

export function deleteAutoscalingGroup(id: string) {
  return request({
    url: `/autoscaling-groups/${id}`,
    method: 'delete'
  })
}

// 批量虚拟机操作
export function batchVMAction(data: {
  vm_ids: string[]
  account_id: number
  action: 'start' | 'stop' | 'reboot' | 'delete'
}) {
  return request<{
    total: number
    success: number
    failed: number
    results: Array<{ vm_id: string; status: string; error?: string }>
    message: string
  }>({
    url: '/compute/vms/batch-action',
    method: 'post',
    data
  })
}
