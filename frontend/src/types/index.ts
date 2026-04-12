export interface CloudAccount {
  id: number
  name: string
  provider_type: string
  credentials: Record<string, string>
  status: string
  description: string
  created_at: string
  updated_at: string
  remarks?: string
  enabled: boolean
  health_status?: string
  balance?: number
  account_number?: string
  last_sync?: string
  sync_time?: string
  domain_id?: number
  resource_assignment_method?: string
}

export interface CreateCloudAccountRequest {
  name: string
  provider_type: string
  credentials: Record<string, string>
  description?: string
  remarks?: string
  enabled?: boolean
  health_status?: string
  balance?: number
  account_number?: string
  last_sync?: string
  sync_time?: string
  domain_id?: number
  resource_assignment_method?: string
}

export interface ScheduledTask {
  id: number
  name: string
  type: string
  frequency: 'once' | 'daily' | 'weekly' | 'monthly' | 'custom'
  triggerTime: string
  validFrom?: string
  validUntil?: string
  status: 'active' | 'inactive'
  created_at: string
  updated_at: string
}

export interface CreateScheduledTaskRequest {
  name: string
  type: string
  frequency: 'once' | 'daily' | 'weekly' | 'monthly' | 'custom'
  triggerTime: string
  validFrom?: string
  validUntil?: string
  status?: 'active' | 'inactive'
}

export interface VirtualMachine {
  id: string
  name: string
  status: string
  instance_type: string
  image_id: string
  os_name?: string                 // Operating system name
  billing_method?: string         // Billing method (pay-as-you-go, subscription, etc.)
  platform?: string              // Cloud platform (alibaba, tencent, aws, azure)
  project_id?: string            // Associated project ID
  vpc_id: string
  subnet_id: string
  private_ip: string
  public_ip: string
  security_groups?: string[]      // Security group IDs
  security_group_names?: string[] // Names of security groups
  zone_id: string
  created_at: string
  region_id: string
}

export interface VPC {
  id: string
  name: string
  cidr: string
  description: string
  status: string
  region_id: string
  created_at: string
}

export interface Subnet {
  id: string
  name: string
  vpc_id: string
  cidr: string
  zone_id: string
  status: string
  created_at: string
}

export interface SecurityGroup {
  id: string
  name: string
  description: string
  vpc_id: string
  created_at: string
}

export interface EIP {
  id: string
  address: string
  bandwidth: number
  status: string
  region_id: string
  created_at: string
}

export interface Image {
  id: string
  name: string
  description: string
  os_name: string
  status: string
  size: number
}

export interface HostTemplate {
  id: string
  name: string
  description: string
  status: string
  instance_type: string
  cpu_arch: string
  memory_size: number
  cpu_count: number
  disk_size: number
  image_id: string
  os_name: string
  os_version: string
  vpc_id: string
  subnet_id: string
  billing_method: string
  platform: string
  project_id: string
  region_id: string
  zone_id: string
  tags: Record<string, string>
  created_at: string
  updated_at: string
}

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

// Geographic types
export interface Region {
  id: string
  name: string
  status: string
  available_zone_count?: number
  vpc_count?: number
  virtual_machine_count?: number
  platform?: string
  created_at: string
}

export interface Zone {
  id: string
  name: string
  region_id: string
  status: string
  l2_network_count?: number
  host_count?: number
  available_host_count?: number
  physical_host_count?: number
  available_physical_host_count?: number
  created_at: string
}

// VPC Interconnect
export interface VPCInterconnect {
  id: string
  name: string
  status: string
  type: string // 专线, VPN
  bandwidth: number // Mbps
  region_id: string
  peer_region: string
  vpc_count?: number
  domain?: string
  platform?: string
  account?: string
  created_at: string
}

// VPC Peering
export interface VPCPeering {
  id: string
  name: string
  status: string
  local_vpc_id: string
  peer_vpc_id: string
  peer_account: string
  domain?: string
  platform?: string
  account?: string
  created_at: string
}

// Route Table
export interface Route {
  destination_cidr: string
  target: string // VPC Gateway, Internet Gateway, VPC Peering, etc.
  status: string
}

export interface RouteTable {
  id: string
  name: string
  vpc_id: string
  status: string
  routes?: Route[]
  domain?: string
  platform?: string
  account?: string
  region_id?: string
  created_at: string
}

// L2 Network
export interface L2Network {
  id: string
  name: string
  vlan_id: number
  status: string
  vpc_id: string
  bandwidth: number // Mbps
  network_count?: number
  platform?: string
  domain?: string
  region_id?: string
  created_at: string
}
