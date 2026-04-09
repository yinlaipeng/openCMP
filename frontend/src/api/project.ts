import request from '@/utils/request'
import type { Project } from '@/types'

export function getProjects(params?: { page?: number; page_size?: number; domain_id?: number }) {
  return request<{ items: Project[], total: number }>({
    url: '/projects',
    method: 'get',
    params
  })
}

export function getProject(id: number) {
  return request<Project>({
    url: `/projects/${id}`,
    method: 'get'
  })
}

export function createProject(data: { name: string; description?: string; domain_id: number; manager_id?: number }) {
  return request<Project>({
    url: '/projects',
    method: 'post',
    data
  })
}

export function updateProject(id: number, data: Partial<{ name: string; description?: string; enabled: boolean }>) {
  return request<Project>({
    url: `/projects/${id}`,
    method: 'put',
    data
  })
}

export function deleteProject(id: number) {
  return request({
    url: `/projects/${id}`,
    method: 'delete'
  })
}

export function enableProject(id: number) {
  return request({
    url: `/projects/${id}/enable`,
    method: 'post'
  })
}

export function disableProject(id: number) {
  return request({
    url: `/projects/${id}/disable`,
    method: 'post'
  })
}