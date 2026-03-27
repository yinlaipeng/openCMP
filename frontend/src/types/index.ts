export interface CloudAccount {
  id: number
  name: string
  provider_type: string
  credentials: Record<string, string>
  status: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateCloudAccountRequest {
  name: string
  provider_type: string
  credentials: Record<string, string>
  description?: string
}

export interface VirtualMachine {
  id: string
  name: string
  status: string
  instance_type: string
  image_id: string
  vpc_id: string
  subnet_id: string
  private_ip: string
  public_ip: string
  region_id: string
  zone_id: string
  created_at: string
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
