// frontend/src/api/networkSync.ts
import request from '@/utils/request'

// ========== 安全组同步 API ==========

export interface SecurityGroup {
  id: number | string
  cloud_account_id: number
  security_group_id: string
  name: string
  description: string
  vpc_id: string
  vpc: string
  guest_cnt: number
  public_scope: string
  is_public: boolean
  is_system: boolean
  is_default_vpc: boolean
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface SecurityGroupListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export interface SecurityGroup {
  id: number | string
  cloud_account_id: number
  security_group_id: string
  name: string
  description: string
  vpc_id: string
  vpc: string
  guest_cnt: number
  public_scope: string
  is_public: boolean
  is_system: boolean
  is_default_vpc: boolean
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export function getSecurityGroups(params?: SecurityGroupListParams) {
  return request<{ items: SecurityGroup[]; total: number; page: number; page_size: number }>({
    url: '/network/security-groups',
    method: 'get',
    params
  })
}

export function getSecurityGroup(id: number | string) {
  return request<SecurityGroup>({
    url: `/network/security-groups/${id}`,
    method: 'get'
  })
}

export function createSecurityGroup(data: Partial<SecurityGroup>) {
  return request<SecurityGroup>({
    url: '/network/security-groups',
    method: 'post',
    data
  })
}

export function deleteSecurityGroup(id: number | string) {
  return request<{ message: string }>({
    url: `/network/security-groups/${id}`,
    method: 'delete'
  })
}

export function batchDeleteSecurityGroups(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/security-groups/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== IP子网同步 API ==========

export interface Subnet {
  id: number | string
  cloud_account_id: number
  subnet_id: string
  name: string
  vpc_id: string
  vpc: string
  cidr: string
  guest_ip_start: string
  guest_ip_end: string
  guest_ip_mask: number
  guest_ip6_mask: string
  guest_gateway: string
  dns: string
  is_auto_alloc: boolean
  is_classic: boolean
  schedtag: string
  ports: number
  ports_used: number
  ports6_used: number
  vnics: number
  zone_id: string
  zone: string
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface SubnetListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
  vpc_id?: string
}

export function getSubnets(params?: SubnetListParams) {
  return request<{ items: Subnet[]; total: number; page: number; page_size: number }>({
    url: '/network/subnets',
    method: 'get',
    params
  })
}

export function getSubnet(id: number | string) {
  return request<Subnet>({
    url: `/network/subnets/${id}`,
    method: 'get'
  })
}

export function createSubnet(data: Partial<Subnet>) {
  return request<Subnet>({
    url: '/network/subnets',
    method: 'post',
    data
  })
}

export function deleteSubnet(id: number | string) {
  return request<{ message: string }>({
    url: `/network/subnets/${id}`,
    method: 'delete'
  })
}

export function batchDeleteSubnets(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/subnets/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 弹性公网IP同步 API ==========

export interface EIP {
  id: number | string
  cloud_account_id: number
  eip_id: string
  name: string
  address: string
  bandwidth: number
  billing_method: string
  resource_id: string
  resource_type: string
  resource_name: string
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface EIPListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getEIPs(params?: EIPListParams) {
  return request<{ items: EIP[]; total: number; page: number; page_size: number }>({
    url: '/network/eips',
    method: 'get',
    params
  })
}

export function getEIP(id: number | string) {
  return request<EIP>({
    url: `/network/eips/${id}`,
    method: 'get'
  })
}

export function createEIP(data: Partial<EIP>) {
  return request<EIP>({
    url: '/network/eips',
    method: 'post',
    data
  })
}

export function bindEIP(id: number | string, resource_id: string) {
  return request<{ message: string }>({
    url: `/network/eips/${id}/bind`,
    method: 'post',
    data: { resource_id }
  })
}

export function unbindEIP(id: number | string) {
  return request<{ message: string }>({
    url: `/network/eips/${id}/unbind`,
    method: 'post'
  })
}

export function deleteEIP(id: number | string) {
  return request<{ message: string }>({
    url: `/network/eips/${id}`,
    method: 'delete'
  })
}

export function batchDeleteEIPs(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/eips/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== SSH密钥同步 API ==========

export interface KeyPair {
  id: number | string
  cloud_account_id: number
  keypair_id: string
  name: string
  description: string
  status: string
  provider_type: string
  account_name: string
  public_key: string
  fingerprint: string
  keypair_type: string
  public_scope: string
  guest_cnt: number
  region_id: string
  region: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface KeyPairListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getKeyPairs(params?: KeyPairListParams) {
  return request<{ items: KeyPair[]; total: number; page: number; page_size: number }>({
    url: '/network/keypairs',
    method: 'get',
    params
  })
}

export function getKeyPair(id: number | string) {
  return request<KeyPair>({
    url: `/network/keypairs/${id}`,
    method: 'get'
  })
}

export function createKeyPair(data: Partial<KeyPair> & { auto_generate?: boolean }) {
  return request<KeyPair>({
    url: '/network/keypairs',
    method: 'post',
    data
  })
}

export function deleteKeyPair(id: number | string) {
  return request<{ message: string }>({
    url: `/network/keypairs/${id}`,
    method: 'delete'
  })
}

export function batchDeleteKeyPairs(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/keypairs/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 区域 API ==========

export interface Region {
  id: string
  name: string
  enabled: boolean
  zone_count: number
  vpc_count: number
  server_count: number
  provider_type: string
  created_at: string
}

export interface RegionListParams {
  page?: number
  page_size?: number
  name?: string
  enabled?: string
  platform?: string
}

export function getRegions(params?: RegionListParams) {
  return request<{ items: Region[]; total: number; page: number; page_size: number }>({
    url: '/network/regions',
    method: 'get',
    params
  })
}

// ========== 可用区 API ==========

export interface Zone {
  id: string
  name: string
  region_id: string
  l2_network_count: number
  host_count: number
  enabled_host_count: number
  physical_host_count: number
  enabled_physical_host_count: number
  status: string
  created_at: string
}

export interface ZoneListParams {
  page?: number
  page_size?: number
  name?: string
  region_id?: string
}

export function getZones(params?: ZoneListParams) {
  return request<{ items: Zone[]; total: number; page: number; page_size: number }>({
    url: '/network/zones',
    method: 'get',
    params
  })
}

export function createZone(data: Partial<Zone>) {
  return request<Zone>({
    url: '/network/zones',
    method: 'post',
    data
  })
}

export function deleteZone(id: string) {
  return request<{ message: string }>({
    url: `/network/zones/${id}`,
    method: 'delete'
  })
}

export function batchDeleteZones(ids: string[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/zones/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== VPC API ==========

export interface VPC {
  id: number | string
  cloud_account_id: number
  vpc_id: string
  name: string
  cidr: string
  ipv6_cidr: string
  allow_external_access: boolean
  network_count: number
  status: string
  provider_type: string
  account_name: string
  owner_domain: string
  region: string
  description: string
  created_at: string
}

export interface VPCListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getVPCs(params?: VPCListParams) {
  return request<{ items: VPC[]; total: number; page: number; page_size: number }>({
    url: '/network/vpcs',
    method: 'get',
    params
  })
}

export function createVPC(data: Partial<VPC>) {
  return request<VPC>({
    url: '/network/vpcs',
    method: 'post',
    data
  })
}

export function deleteVPC(id: number | string) {
  return request<{ message: string }>({
    url: `/network/vpcs/${id}`,
    method: 'delete'
  })
}

export function batchDeleteVPCs(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/vpcs/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 全局VPC API ==========

export interface GlobalVPC {
  id: number | string
  cloud_account_id: number
  name: string
  description: string
  vpc_count: number
  status: string
  provider_type: string
  account_name: string
  owner_domain: string
  tags: { key: string; value: string }[]
  created_at: string
}

export interface GlobalVPCListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
}

export function getGlobalVPCs(params?: GlobalVPCListParams) {
  return request<{ items: GlobalVPC[]; total: number; page: number; page_size: number }>({
    url: '/network/global-vpcs',
    method: 'get',
    params
  })
}

export function createGlobalVPC(data: Partial<GlobalVPC>) {
  return request<GlobalVPC>({
    url: '/network/global-vpcs',
    method: 'post',
    data
  })
}

export function deleteGlobalVPC(id: number | string) {
  return request<{ message: string }>({
    url: `/network/global-vpcs/${id}`,
    method: 'delete'
  })
}

export function batchDeleteGlobalVPCs(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/global-vpcs/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== VPC互联 API ==========

export interface VPCInterconnect {
  id: number | string
  name: string
  vpc_count: number
  attribution_scope: string
  interconnect_type: string
  bandwidth: number
  status: string
  provider_type: string
  account_name: string
  description: string
  tags: { key: string; value: string }[]
  created_at: string
}

export interface VPCInterconnectListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
}

export function getVPCInterconnects(params?: VPCInterconnectListParams) {
  return request<{ items: VPCInterconnect[]; total: number; page: number; page_size: number }>({
    url: '/network/vpc-interconnects',
    method: 'get',
    params
  })
}

export function deleteVPCInterconnect(id: number | string) {
  return request<{ message: string }>({
    url: `/network/vpc-interconnects/${id}`,
    method: 'delete'
  })
}

export function batchDeleteVPCInterconnects(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/vpc-interconnects/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 二层网络 API ==========

export interface L2Network {
  id: number | string
  name: string
  vpc_id: string
  vpc_name: string
  bandwidth: number
  vlan_id: string
  mtu: number
  network_count: number
  status: string
  provider_type: string
  owner_domain: string
  region: string
  tags: { key: string; value: string }[]
  created_at: string
}

export interface L2NetworkListParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getL2Networks(params?: L2NetworkListParams) {
  return request<{ items: L2Network[]; total: number; page: number; page_size: number }>({
    url: '/network/l2-networks',
    method: 'get',
    params
  })
}

export function createL2Network(data: Partial<L2Network>) {
  return request<L2Network>({
    url: '/network/l2-networks',
    method: 'post',
    data
  })
}

export function deleteL2Network(id: number | string) {
  return request<{ message: string }>({
    url: `/network/l2-networks/${id}`,
    method: 'delete'
  })
}

export function batchDeleteL2Networks(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/l2-networks/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 路由表 API ==========

export interface RouteTable {
  id: number | string
  name: string
  vpc_id: string
  vpc_name: string
  status: string
  provider_type: string
  account_name: string
  owner_domain: string
  region: string
  created_at: string
}

export interface RouteTableListParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getRouteTables(params?: RouteTableListParams) {
  return request<{ items: RouteTable[]; total: number; page: number; page_size: number }>({
    url: '/network/route-tables',
    method: 'get',
    params
  })
}

// ========== NAT网关 API ==========

export interface NATGateway {
  id: number | string
  cloud_account_id: number
  nat_gateway_id: string
  name: string
  description: string
  nat_type: string // 公网NAT/私网NAT
  specification: string
  billing_method: string // Postpaid/Prepaid
  vpc_id: string
  vpc_name: string
  subnet_id: string
  subnet_name: string
  eip_id: string
  eip_address: string
  snat_table_entries: number
  dnat_table_entries: number
  owner_domain: string
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface NATGatewayListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  region?: string
  nat_type?: string
}

export function getNATGateways(params?: NATGatewayListParams) {
  return request<{ items: NATGateway[]; total: number; page: number; page_size: number }>({
    url: '/network/nat-gateways',
    method: 'get',
    params
  })
}

export function getNATGateway(id: number | string) {
  return request<{ nat_gateway: NATGateway; rules: NATRule[] }>({
    url: `/network/nat-gateways/${id}`,
    method: 'get'
  })
}

export function createNATGateway(data: Partial<NATGateway>) {
  return request<NATGateway>({
    url: '/network/nat-gateways',
    method: 'post',
    data
  })
}

export function updateNATGateway(id: number | string, data: Partial<NATGateway>) {
  return request<NATGateway>({
    url: `/network/nat-gateways/${id}`,
    method: 'put',
    data
  })
}

export function deleteNATGateway(id: number | string) {
  return request<{ message: string }>({
    url: `/network/nat-gateways/${id}`,
    method: 'delete'
  })
}

export function batchDeleteNATGateways(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/nat-gateways/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== NAT规则 API ==========

export interface NATRule {
  id: number | string
  nat_gateway_id: number
  rule_id: string
  rule_type: string // SNAT/DNAT
  name: string
  external_ip: string
  external_port: string
  internal_ip: string
  internal_port: string
  protocol: string // TCP/UDP/ALL
  status: string
  created_at: string
  updated_at: string
}

export interface NATRuleListParams {
  page?: number
  page_size?: number
}

export function getNATRules(natGatewayId: number | string, params?: NATRuleListParams) {
  return request<{ items: NATRule[]; total: number; page: number; page_size: number }>({
    url: `/network/nat-gateways/${natGatewayId}/rules`,
    method: 'get',
    params
  })
}

export function createNATRule(natGatewayId: number | string, data: Partial<NATRule>) {
  return request<NATRule>({
    url: `/network/nat-gateways/${natGatewayId}/rules`,
    method: 'post',
    data
  })
}

export function deleteNATRule(natGatewayId: number | string, ruleId: number | string) {
  return request<{ message: string }>({
    url: `/network/nat-gateways/${natGatewayId}/rules/${ruleId}`,
    method: 'delete'
  })
}

// ========== IPv6网关 API ==========

export interface IPv6Gateway {
  id: number | string
  cloud_account_id: number
  ipv6_gateway_id: string
  name: string
  vpc_id: string
  vpc_name: string
  specification: string
  ipv6_cidr: string
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface IPv6GatewayListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  region?: string
}

export function getIPv6Gateways(params?: IPv6GatewayListParams) {
  return request<{ items: IPv6Gateway[]; total: number; page: number; page_size: number }>({
    url: '/network/ipv6-gateways',
    method: 'get',
    params
  })
}

export function getIPv6Gateway(id: number | string) {
  return request<IPv6Gateway>({
    url: `/network/ipv6-gateways/${id}`,
    method: 'get'
  })
}

export function createIPv6Gateway(data: Partial<IPv6Gateway>) {
  return request<IPv6Gateway>({
    url: '/network/ipv6-gateways',
    method: 'post',
    data
  })
}

export function deleteIPv6Gateway(id: number | string) {
  return request<{ message: string }>({
    url: `/network/ipv6-gateways/${id}`,
    method: 'delete'
  })
}

export function batchDeleteIPv6Gateways(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/ipv6-gateways/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== DNS Zone API ==========

export interface DNSZone {
  id: number | string
  cloud_account_id: number
  dns_zone_id: string
  name: string
  vpc_count: number
  attribution_scope: string
  region_id: string
  region: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface DNSZoneListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  region?: string
}

export function getDNSZones(params?: DNSZoneListParams) {
  return request<{ items: DNSZone[]; total: number; page: number; page_size: number }>({
    url: '/network/dns-zones',
    method: 'get',
    params
  })
}

export function getDNSZone(id: number | string) {
  return request<{ dns_zone: DNSZone; records: DNSRecord[] }>({
    url: `/network/dns-zones/${id}`,
    method: 'get'
  })
}

export function createDNSZone(data: Partial<DNSZone>) {
  return request<DNSZone>({
    url: '/network/dns-zones',
    method: 'post',
    data
  })
}

export function deleteDNSZone(id: number | string) {
  return request<{ message: string }>({
    url: `/network/dns-zones/${id}`,
    method: 'delete'
  })
}

export function batchDeleteDNSZones(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/dns-zones/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== DNS Record API ==========

export interface DNSRecord {
  id: number | string
  dns_zone_id: number
  record_id: string
  name: string
  type: string // A/CNAME/MX/TXT/NS等
  value: string
  ttl: number
  priority: number
  status: string
  created_at: string
  updated_at: string
}

export interface DNSRecordListParams {
  page?: number
  page_size?: number
}

export function getDNSRecords(dnsZoneId: number | string, params?: DNSRecordListParams) {
  return request<{ items: DNSRecord[]; total: number; page: number; page_size: number }>({
    url: `/network/dns-zones/${dnsZoneId}/records`,
    method: 'get',
    params
  })
}

export function createDNSRecord(dnsZoneId: number | string, data: Partial<DNSRecord>) {
  return request<DNSRecord>({
    url: `/network/dns-zones/${dnsZoneId}/records`,
    method: 'post',
    data
  })
}

export function deleteDNSRecord(dnsZoneId: number | string, recordId: number | string) {
  return request<{ message: string }>({
    url: `/network/dns-zones/${dnsZoneId}/records/${recordId}`,
    method: 'delete'
  })
}

// ========== 负载均衡实例 API ==========

export interface LBInstance {
  id: number | string
  cloud_account_id: number
  lb_id: string
  name: string
  address: string
  address_type: string
  specification: string
  bandwidth: number
  vpc_id: string
  vpc: string
  security_group_id: string
  security_group: string
  charging_method: string
  listener_count: number
  backend_count: number
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  region: string
  status: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface LBInstanceListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getLBInstances(params?: LBInstanceListParams) {
  return request<{ items: LBInstance[]; total: number; page: number; page_size: number }>({
    url: '/network/lb-instances',
    method: 'get',
    params
  })
}

export function getLBInstance(id: number | string) {
  return request<LBInstance>({
    url: `/network/lb-instances/${id}`,
    method: 'get'
  })
}

export function createLBInstance(data: Partial<LBInstance>) {
  return request<LBInstance>({
    url: '/network/lb-instances',
    method: 'post',
    data
  })
}

export function syncLBStatus(id: number | string) {
  return request<{ message: string }>({
    url: `/network/lb-instances/${id}/sync`,
    method: 'post'
  })
}

export function batchSyncLBStatus(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/lb-instances/batch-sync',
    method: 'post',
    data: { ids }
  })
}

export function deleteLBInstance(id: number | string) {
  return request<{ message: string }>({
    url: `/network/lb-instances/${id}`,
    method: 'delete'
  })
}

export function batchDeleteLBInstances(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/lb-instances/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 访问控制 API ==========

export interface LBACL {
  id: number | string
  cloud_account_id: number
  acl_id: string
  name: string
  address_source: string
  remarks: string
  listener_count: number
  status: string
  provider_type: string
  account_name: string
  shared_scope: string
  region: string
  project_id: number
  project_name: string
  created_at: string
  updated_at: string
}

export interface LBACLListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getLBACLs(params?: LBACLListParams) {
  return request<{ items: LBACL[]; total: number; page: number; page_size: number }>({
    url: '/network/lb-acls',
    method: 'get',
    params
  })
}

export function getLBACL(id: number | string) {
  return request<LBACL>({
    url: `/network/lb-acls/${id}`,
    method: 'get'
  })
}

export function createLBACL(data: Partial<LBACL>) {
  return request<LBACL>({
    url: '/network/lb-acls',
    method: 'post',
    data
  })
}

export function deleteLBACL(id: number | string) {
  return request<{ message: string }>({
    url: `/network/lb-acls/${id}`,
    method: 'delete'
  })
}

export function batchDeleteLBACLs(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/lb-acls/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// ========== 证书 API ==========

export interface LBCertificate {
  id: number | string
  cloud_account_id: number
  cert_id: string
  name: string
  domain_name: string
  expiration: string
  subject_alternative_names: string
  listener_count: number
  status: string
  shared_scope: string
  project_id: number
  project_name: string
  provider_type: string
  account_name: string
  region: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface LBCertificateListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  region?: string
}

export function getLBCertificates(params?: LBCertificateListParams) {
  return request<{ items: LBCertificate[]; total: number; page: number; page_size: number }>({
    url: '/network/lb-certificates',
    method: 'get',
    params
  })
}

export function getLBCertificate(id: number | string) {
  return request<LBCertificate>({
    url: `/network/lb-certificates/${id}`,
    method: 'get'
  })
}

export function createLBCertificate(data: Partial<LBCertificate>) {
  return request<LBCertificate>({
    url: '/network/lb-certificates',
    method: 'post',
    data
  })
}

export function deleteLBCertificate(id: number | string) {
  return request<{ message: string }>({
    url: `/network/lb-certificates/${id}`,
    method: 'delete'
  })
}

export function batchDeleteLBCertificates(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/lb-certificates/batch-delete',
    method: 'post',
    data: { ids }
  })
}

export function changeCertificateProject(id: number | string, project_id: number) {
  return request<{ message: string }>({
    url: `/network/lb-certificates/${id}/change-project`,
    method: 'post',
    data: { project_id }
  })
}

export function setupCertificateSharing(id: number | string, shared_scope: string) {
  return request<{ message: string }>({
    url: `/network/lb-certificates/${id}/setup-sharing`,
    method: 'post',
    data: { shared_scope }
  })
}

// ========== CDN域名 API ==========

export interface CDNDomain {
  id: number | string
  cloud_account_id: number
  cdn_id: string
  name: string
  cname: string
  area: string
  service_type: string
  origin_type: string
  origin_address: string
  status: string
  provider_type: string
  account_name: string
  project_id: number
  project_name: string
  tags: { key: string; value: string }[]
  created_at: string
  updated_at: string
}

export interface CDNDomainListParams {
  page?: number
  page_size?: number
  cloud_account_id?: number
  name?: string
  status?: string
  platform?: string
  area?: string
}

export function getCDNDomains(params?: CDNDomainListParams) {
  return request<{ items: CDNDomain[]; total: number; page: number; page_size: number }>({
    url: '/network/cdn-domains',
    method: 'get',
    params
  })
}

export function getCDNDomain(id: number | string) {
  return request<CDNDomain>({
    url: `/network/cdn-domains/${id}`,
    method: 'get'
  })
}

export function createCDNDomain(data: Partial<CDNDomain>) {
  return request<CDNDomain>({
    url: '/network/cdn-domains',
    method: 'post',
    data
  })
}

export function syncCDNStatus(id: number | string) {
  return request<{ message: string }>({
    url: `/network/cdn-domains/${id}/sync`,
    method: 'post'
  })
}

export function batchSyncCDNStatus(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/cdn-domains/batch-sync',
    method: 'post',
    data: { ids }
  })
}

export function deleteCDNDomain(id: number | string) {
  return request<{ message: string }>({
    url: `/network/cdn-domains/${id}`,
    method: 'delete'
  })
}

export function batchDeleteCDNDomains(ids: number[]) {
  return request<{ total: number; success: number; failed: number; message: string }>({
    url: '/network/cdn-domains/batch-delete',
    method: 'post',
    data: { ids }
  })
}
