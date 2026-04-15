import request from '@/utils/request'
import type { ScheduledTask, CreateScheduledTaskRequest } from '@/types'

export function getScheduledTasks(params?: { page?: number; page_size?: number; cloud_account_id?: number }) {
  return request<{ items: ScheduledTask[], total: number }>({
    url: '/scheduled-tasks',
    method: 'get',
    params
  })
}

export function getScheduledTask(id: number) {
  return request<ScheduledTask>({
    url: `/scheduled-tasks/${id}`,
    method: 'get'
  })
}

export function createScheduledTask(data: CreateScheduledTaskRequest) {
  return request<ScheduledTask>({
    url: '/scheduled-tasks',
    method: 'post',
    data
  })
}

export function updateScheduledTask(id: number, data: Partial<CreateScheduledTaskRequest>) {
  return request<ScheduledTask>({
    url: `/scheduled-tasks/${id}`,
    method: 'put',
    data
  })
}

export function deleteScheduledTask(id: number) {
  return request({
    url: `/scheduled-tasks/${id}`,
    method: 'delete'
  })
}

export function updateScheduledTaskStatus(id: number, status: 'active' | 'inactive') {
  return request({
    url: `/scheduled-tasks/${id}/status`,
    method: 'post',
    data: { status }
  })
}

// 执行定时任务（手动触发同步）
export function executeScheduledTask(id: number) {
  return request<{ message: string; task_id: number; task_name: string; statistics: Record<string, number> }>({
    url: `/scheduled-tasks/${id}/execute`,
    method: 'post'
  })
}