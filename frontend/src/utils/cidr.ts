/**
 * CIDR Validation and Utility Functions
 * Used for validating and parsing CIDR notation for VPC and subnet management
 */

/**
 * Validates CIDR format like "10.0.0.0/16"
 * @param cidr - CIDR string to validate
 * @returns true if valid CIDR format
 */
export function validateCIDR(cidr: string): boolean {
  // Match CIDR pattern: IPv4 address followed by /prefix
  const regex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\/(\d{1,2})$/;
  const match = cidr.match(regex);

  if (!match) {
    return false;
  }

  // Extract IP segments and prefix
  const [, seg1, seg2, seg3, seg4, prefixStr] = match;
  const segments = [seg1, seg2, seg3, seg4].map(Number);
  const prefix = parseInt(prefixStr, 10);

  // Validate each IP segment is between 0-255
  for (const seg of segments) {
    if (seg < 0 || seg > 255) {
      return false;
    }
  }

  // Validate prefix is between 0-32
  if (prefix < 0 || prefix > 32) {
    return false;
  }

  return true;
}

/**
 * Parses CIDR to get IP range start and mask
 * @param cidr - Valid CIDR string
 * @returns Object with start (IP as number) and mask, or null if invalid
 */
export function parseCIDR(cidr: string): { start: number; mask: number } | null {
  if (!validateCIDR(cidr)) {
    return null;
  }

  const regex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\/(\d{1,2})$/;
  const match = cidr.match(regex);

  if (!match) {
    return null;
  }

  const [, seg1, seg2, seg3, seg4, prefixStr] = match;
  const segments = [seg1, seg2, seg3, seg4].map(Number);
  const prefix = parseInt(prefixStr, 10);

  // Convert IP segments to a 32-bit number
  // e.g., 10.0.0.0 -> 0x0A000000
  const start =
    ((segments[0] << 24) >>> 0) +
    ((segments[1] << 16) >>> 0) +
    ((segments[2] << 8) >>> 0) +
    segments[3];

  // Calculate network mask using bitwise operations
  // For prefix N, mask has N leading 1s followed by (32-N) 0s
  // e.g., prefix 16 -> 0xFFFF0000
  const mask = prefix === 0 ? 0 : ((0xFFFFFFFF << (32 - prefix)) >>> 0);

  return { start, mask };
}

/**
 * Checks if subnet CIDR is within VPC CIDR range
 * @param subnetCIDR - Subnet CIDR string (e.g., "10.0.1.0/24")
 * @param vpcCIDR - VPC CIDR string (e.g., "10.0.0.0/16")
 * @returns true if subnet is within VPC range
 */
export function isSubnetInVPC(subnetCIDR: string, vpcCIDR: string): boolean {
  // Both CIDRs must be valid
  const subnetParsed = parseCIDR(subnetCIDR);
  const vpcParsed = parseCIDR(vpcCIDR);

  if (!subnetParsed || !vpcParsed) {
    return false;
  }

  const { start: subnetStart, mask: subnetMask } = subnetParsed;
  const { start: vpcStart, mask: vpcMask } = vpcParsed;

  // Subnet prefix must be greater than or equal to VPC prefix
  // (subnet must be smaller or equal in size to VPC)
  // Actually, for a subnet to be inside a VPC, subnet prefix must be >= vpc prefix
  // (meaning subnet is smaller or equal in size, providing more specific routing)
  const vpcPrefix = getPrefixFromMask(vpcMask);
  const subnetPrefix = getPrefixFromMask(subnetMask);

  if (subnetPrefix < vpcPrefix) {
    return false;
  }

  // Check if subnet's network address is within VPC's range
  // subnet.start & vpc.mask should equal vpc.start
  if ((subnetStart & vpcMask) >>> 0 !== vpcStart) {
    return false;
  }

  return true;
}

/**
 * Get prefix length from mask
 * @param mask - Network mask as 32-bit number
 * @returns Prefix length (0-32)
 */
function getPrefixFromMask(mask: number): number {
  if (mask === 0) return 0;
  let count = 0;
  let m = mask >>> 0;
  while (m & 0x80000000) {
    count++;
    m = (m << 1) >>> 0;
  }
  return count;
}

/**
 * Format CIDR for display
 * @param cidr - CIDR string to format
 * @returns Formatted CIDR string or original if invalid
 */
export function formatCIDR(cidr: string): string {
  if (!validateCIDR(cidr)) {
    return cidr; // Return original if invalid
  }

  const regex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\/(\d{1,2})$/;
  const match = cidr.match(regex);

  if (!match) {
    return cidr;
  }

  const [, seg1, seg2, seg3, seg4, prefix] = match;

  // Normalize the IP segments (remove leading zeros)
  const normalizedSegs = [seg1, seg2, seg3, seg4].map(s => parseInt(s, 10));

  // Return normalized format: x.x.x.x/prefix
  return `${normalizedSegs.join('.')}/${prefix}`;
}

/**
 * Get the network address from an IP address within a CIDR
 * @param ip - IP address string (e.g., "10.0.1.100")
 * @param cidr - CIDR notation (e.g., "10.0.0.0/16")
 * @returns Network address or null if invalid
 */
export function getNetworkAddress(ip: string, cidr: string): string | null {
  const parsed = parseCIDR(cidr);
  if (!parsed) return null;

  const { mask } = parsed;

  // Parse the IP address
  const ipRegex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/;
  const ipMatch = ip.match(ipRegex);
  if (!ipMatch) return null;

  const segments = ipMatch.slice(1, 5).map(Number);
  for (const seg of segments) {
    if (seg < 0 || seg > 255) return null;
  }

  const ipNum =
    ((segments[0] << 24) >>> 0) +
    ((segments[1] << 16) >>> 0) +
    ((segments[2] << 8) >>> 0) +
    segments[3];

  const networkNum = (ipNum & mask) >>> 0;

  // Convert back to IP string
  const n1 = (networkNum >>> 24) & 0xFF;
  const n2 = (networkNum >>> 16) & 0xFF;
  const n3 = (networkNum >>> 8) & 0xFF;
  const n4 = networkNum & 0xFF;

  return `${n1}.${n2}.${n3}.${n4}`;
}

/**
 * Calculate the number of usable IP addresses in a CIDR range
 * @param cidr - CIDR notation
 * @returns Number of usable IPs (total - 2 for network and broadcast) or 0 if invalid
 */
export function getUsableIPCount(cidr: string): number {
  const parsed = parseCIDR(cidr);
  if (!parsed) return 0;

  const { mask } = parsed;
  const prefix = getPrefixFromMask(mask);

  // Total IPs = 2^(32-prefix)
  // Usable IPs = Total - 2 (network and broadcast addresses)
  // For /31 and /32, special rules apply
  const totalIPs = Math.pow(2, 32 - prefix);

  if (prefix >= 31) {
    // /31 has 2 usable IPs (point-to-point link)
    // /32 has 1 usable IP (single host)
    return totalIPs;
  }

  return totalIPs - 2;
}