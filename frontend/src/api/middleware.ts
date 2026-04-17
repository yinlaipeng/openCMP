import request from '@/utils/request'

// Kafka Types
export interface KafkaInstance {
  id: string
  name: string
  label: string
  status: string
  version: string
  storage: string
  bandwidth: string
  endpoint: string
  retention: string
  billing_method: string
  platform: string
  account_name: string
  project_id: string
  region_id: string
  zone_id: string
  tags: Record<string, string>
  created_at: string
}

export interface KafkaConfig {
  account_id: number
  name: string
  version: string
  storage?: number
  bandwidth?: number
  retention?: number
  vpc_id: string
  subnet_id: string
  zone_id?: string
  tags?: Record<string, string>
}

export interface KafkaFilter {
  account_id: number
  instance_id?: string
  status?: string
  version?: string
  region_id?: string
}

// Elasticsearch Types
export interface ElasticsearchInstance {
  id: string
  name: string
  label: string
  status: string
  instance_type: string
  config: string
  version: string
  storage: string
  billing_method: string
  platform: string
  account_name: string
  project_id: string
  region_id: string
  zone_id: string
  tags: Record<string, string>
  created_at: string
}

export interface ElasticsearchConfig {
  account_id: number
  name: string
  version: string
  instance_type: string
  node_count?: number
  storage?: number
  vpc_id: string
  subnet_id: string
  zone_id?: string
  tags?: Record<string, string>
}

export interface ElasticsearchFilter {
  account_id: number
  instance_id?: string
  status?: string
  version?: string
  instance_type?: string
  region_id?: string
}

// Kafka API
export function listKafka(filter: KafkaFilter) {
  return request.get<KafkaInstance[]>('/middleware/kafka', { params: filter })
}

export function getKafka(accountId: number, instanceId: string) {
  return request.get<KafkaInstance>(`/middleware/kafka/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function createKafka(config: KafkaConfig) {
  return request.post<KafkaInstance>('/middleware/kafka', config)
}

export function deleteKafka(accountId: number, instanceId: string) {
  return request.delete(`/middleware/kafka/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function kafkaAction(accountId: number, instanceId: string, action: string) {
  return request.post(`/middleware/kafka/${instanceId}/action`, { action }, {
    params: { account_id: accountId }
  })
}

export function resizeKafka(accountId: number, instanceId: string, storage?: number, bandwidth?: number) {
  return request.post(`/middleware/kafka/${instanceId}/resize`, {
    storage,
    bandwidth
  }, {
    params: { account_id: accountId }
  })
}

// Elasticsearch API
export function listElasticsearch(filter: ElasticsearchFilter) {
  return request.get<ElasticsearchInstance[]>('/middleware/elasticsearch', { params: filter })
}

export function getElasticsearch(accountId: number, instanceId: string) {
  return request.get<ElasticsearchInstance>(`/middleware/elasticsearch/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function createElasticsearch(config: ElasticsearchConfig) {
  return request.post<ElasticsearchInstance>('/middleware/elasticsearch', config)
}

export function deleteElasticsearch(accountId: number, instanceId: string) {
  return request.delete(`/middleware/elasticsearch/${instanceId}`, {
    params: { account_id: accountId }
  })
}

export function elasticsearchAction(accountId: number, instanceId: string, action: string) {
  return request.post(`/middleware/elasticsearch/${instanceId}/action`, { action }, {
    params: { account_id: accountId }
  })
}

export function resizeElasticsearch(accountId: number, instanceId: string, nodeCount?: number, storage?: number) {
  return request.post(`/middleware/elasticsearch/${instanceId}/resize`, {
    node_count: nodeCount,
    storage
  }, {
    params: { account_id: accountId }
  })
}