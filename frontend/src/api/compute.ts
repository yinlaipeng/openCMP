import request from '@/utils/request'
import type { VirtualMachine, Image } from '@/types'

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
    url: '/compute/images',
    method: 'get',
    params
  })
}
