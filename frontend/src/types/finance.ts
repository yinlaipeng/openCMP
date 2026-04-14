// frontend/src/types/finance.ts

// 账单记录
export interface Bill {
  id: number
  cloud_account_id: number
  billing_cycle: string
  product_type: string
  product_name: string
  instance_id: string
  usage_amount: number
  unit_price: number
  total_cost: number
  currency: string
  billing_method: string
  status: string
  provider_type: string
  created_at: string
}

// 订单记录
export interface Order {
  id: number
  cloud_account_id: number
  order_id: string
  order_type: string
  product_type: string
  product_name: string
  instance_id: string
  amount: number
  currency: string
  status: string
  payment_time?: string
  effective_time: string
  expire_time?: string
  provider_type: string
  created_at: string
}

// 预算配置
export interface Budget {
  id: number
  cloud_account_id: number
  name: string
  type: string
  amount: number
  alert_threshold: number
  current_usage: number
  status: string
  created_at: string
  updated_at: string
}

// 成本异常
export interface CostAnomaly {
  id: number
  cloud_account_id: number
  anomaly_type: string
  detected_at: string
  period: string
  expected_cost: number
  actual_cost: number
  deviation_rate: number
  severity: string
  status: string
  resolution?: string
  created_at: string
}

// 续费资源
export interface RenewalResource {
  id: number
  cloud_account_id: number
  instance_id: string
  instance_name: string
  product_type: string
  expire_time: string
  days_remaining: number
  renewal_price: number
  status: string
}

// 成本分析数据
export interface CostAnalysisData {
  period: string
  total_cost: number
  product_costs: Record<string, number>
  trend: number
}