export interface CloudAccountWizardStep {
  id: number
  title: string
  description: string
}

export interface CloudAccountProvider {
  id: string
  name: string
  displayName: string
  icon: string
}

export interface AlibabaCloudConfig {
  name: string
  remarks?: string
  accountType: 'public' | 'finance' // finance暂不支持
  accessKeyId: string
  accessKeySecret: string
  resourceAssignmentMethod: string[]
  syncStrategy?: {
    policy: string
    scope: 'resource_tags' | 'project_tags'
    defaultProject: number
  }
  cloudProjectMapping?: {
    defaultProject: number
  }
  cloudSubscriptionMapping?: {
    defaultProject: number
  }
  specificProject?: {
    defaultProject: number
  }
  syncRegions: string[]
  syncSchedule?: {
    name: string
    type: 'sync_cloud_account'
    frequency: 'once' | 'daily' | 'weekly' | 'monthly' | 'custom'
    triggerTime: string // HH:mm format
    validFrom?: Date
    validUntil?: Date
  }
}