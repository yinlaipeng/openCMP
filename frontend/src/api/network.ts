import request from '@/utils/request'
import type { VPC, Subnet, SecurityGroup, EIP } from '@/types'

// VPC APIs
export function getVPCs(params?: { account_id?: number; region_id?: string }) {
  return request<VPC[]>({
    url: '/network/vpcs',
    method: 'get',
    params
  })
}

export function createVPC(data: {
  account_id: number
  name: string
  cidr: string
  description?: string
}) {
  return request<VPC>({
    url: '/network/vpcs',
    method: 'post',
    data
  })
}

export function deleteVPC(id: string, account_id: number) {
  return request({
    url: `/network/vpcs/${id}`,
    method: 'delete',
    params: { account_id }
  })
}

// Subnet APIs
export function getSubnets(params?: { account_id?: number; vpc_id?: string; zone_id?: string }) {
  return request<Subnet[]>({
    url: '/network/subnets',
    method: 'get',
    params
  })
}

export function createSubnet(data: {
  account_id: number
  name: string
  vpc_id: string
  cidr: string
  zone_id: string
  description?: string
}) {
  return request<Subnet>({
    url: '/network/subnets',
    method: 'post',
    data
  })
}

// Security Group APIs
export function getSecurityGroups(params?: { account_id?: number; vpc_id?: string }) {
  return request<SecurityGroup[]>({
    url: '/network/security-groups',
    method: 'get',
    params
  })
}

export function createSecurityGroup(data: {
  account_id: number
  name: string
  vpc_id: string
  description?: string
}) {
  return request<SecurityGroup>({
    url: '/network/security-groups',
    method: 'post',
    data
  })
}

// EIP APIs
export function getEIPs(params?: { account_id?: number; region_id?: string; status?: string }) {
  return request<EIP[]>({
    url: '/network/eips',
    method: 'get',
    params
  })
}

export function createEIP(data: {
  account_id: number
  bandwidth: number
  region_id: string
}) {
  return request<EIP>({
    url: '/network/eips',
    method: 'post',
    data
  })
}
