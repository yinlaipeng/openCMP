import request from '@/utils/request'
import type { VPC, Subnet, SecurityGroup, EIP, Region, Zone, VPCInterconnect, VPCPeering, RouteTable, L2Network } from '@/types'

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
  ipv6_cidr?: string
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

export function deleteSubnet(id: string, account_id: number) {
  return request({
    url: `/network/subnets/${id}`,
    method: 'delete',
    params: { account_id }
  })
}

export function updateSubnet(id: string, account_id: number, data: { name?: string; description?: string; tags?: Record<string, string> }) {
  return request<Subnet>({
    url: `/network/subnets/${id}`,
    method: 'put',
    params: { account_id },
    data
  })
}

export function changeSubnetProject(id: string, account_id: number, project_id: number) {
  return request({
    url: `/network/subnets/${id}/change-project`,
    method: 'post',
    params: { account_id },
    data: { project_id }
  })
}

export function splitSubnet(id: string, account_id: number, new_cidrs: string[]) {
  return request({
    url: `/network/subnets/${id}/split`,
    method: 'post',
    params: { account_id },
    data: { new_cidrs }
  })
}

export function reserveIP(id: string, account_id: number, data: { ips: string[]; reason?: string; reserved_by?: string }) {
  return request({
    url: `/network/subnets/${id}/reserve-ip`,
    method: 'post',
    params: { account_id },
    data
  })
}

export function releaseIP(id: string, account_id: number, ips: string[]) {
  return request({
    url: `/network/subnets/${id}/release-ip`,
    method: 'post',
    params: { account_id },
    data: { ips }
  })
}

export function getReservedIPs(account_id: number, subnet_id: string) {
  return request({
    url: '/network/subnets/reserved-ips',
    method: 'get',
    params: { account_id, subnet_id }
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

export function deleteSecurityGroup(id: string, account_id: number) {
  return request({
    url: `/network/security-groups/${id}`,
    method: 'delete',
    params: { account_id }
  })
}

export function addSecurityGroupRule(id: string, account_id: number, data: {
  direction: string
  protocol: string
  port_range?: string
  cidr: string
  action: string
  description?: string
  priority?: number
}) {
  return request({
    url: `/network/security-groups/${id}/rules`,
    method: 'post',
    params: { account_id },
    data
  })
}

export function deleteSecurityGroupRule(id: string, rule_id: string, account_id: number) {
  return request({
    url: `/network/security-groups/${id}/rules/${rule_id}`,
    method: 'delete',
    params: { account_id }
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

export function deleteEIP(id: string, account_id: number) {
  return request({
    url: `/network/eips/${id}`,
    method: 'delete',
    params: { account_id }
  })
}

export function bindEIP(id: string, account_id: number, data: { resource_id: string; resource_type: string }) {
  return request({
    url: `/network/eips/${id}/bind`,
    method: 'post',
    params: { account_id },
    data
  })
}

export function unbindEIP(id: string, account_id: number) {
  return request({
    url: `/network/eips/${id}/unbind`,
    method: 'post',
    params: { account_id }
  })
}

// Region APIs
export function getRegions(params?: { account_id?: number }) {
  return request<Region[]>({
    url: '/network/regions',
    method: 'get',
    params
  })
}

// Zone APIs
export function getZones(params?: { account_id?: number; region_id?: string }) {
  return request<Zone[]>({
    url: '/network/zones',
    method: 'get',
    params
  })
}

// VPC Interconnect APIs
export function getVPCInterconnects(params?: { account_id?: number; region_id?: string }) {
  return request<VPCInterconnect[]>({
    url: '/network/vpc-interconnects',
    method: 'get',
    params
  })
}

export function createVPCInterconnect(data: {
  account_id: number
  name: string
  type: string
  bandwidth: number
  region_id: string
  peer_region: string
  description?: string
}) {
  return request<VPCInterconnect>({
    url: '/network/vpc-interconnects',
    method: 'post',
    data
  })
}

// VPC Peering APIs
export function getVPCPeerings(params?: { account_id?: number; local_vpc_id?: string }) {
  return request<VPCPeering[]>({
    url: '/network/vpc-peerings',
    method: 'get',
    params
  })
}

export function createVPCPeering(data: {
  account_id: number
  name: string
  local_vpc_id: string
  peer_vpc_id: string
  peer_account: string
  description?: string
}) {
  return request<VPCPeering>({
    url: '/network/vpc-peerings',
    method: 'post',
    data
  })
}

// Route Table APIs
export function getRouteTables(params?: { account_id?: number; vpc_id?: string }) {
  return request<RouteTable[]>({
    url: '/network/route-tables',
    method: 'get',
    params
  })
}

export function createRouteTable(data: {
  account_id: number
  name: string
  vpc_id: string
  routes: any[] // Using any for routes initially
  description?: string
}) {
  return request<RouteTable>({
    url: '/network/route-tables',
    method: 'post',
    data
  })
}

// L2 Network APIs
export function getL2Networks(params?: { account_id?: number; vpc_id?: string }) {
  return request<L2Network[]>({
    url: '/network/l2-networks',
    method: 'get',
    params
  })
}

export function createL2Network(data: {
  account_id: number
  name: string
  vlan_id: number
  vpc_id: string
  subnets: string[]
  description?: string
}) {
  return request<L2Network>({
    url: '/network/l2-networks',
    method: 'post',
    data
  })
}
