import request from '@/utils/request'

// ============= 站内信 =============

export function getMessages(params?: {
  user_id?: number
  read?: boolean
  page?: number
  page_size?: number
}) {
  return request({ url: '/messages', method: 'get', params })
}

export function getMessage(id: number) {
  return request({ url: `/messages/${id}`, method: 'get' })
}

export function markRead(id: number) {
  return request({ url: `/messages/${id}/read`, method: 'put' })
}

export function markAllRead(userId: number) {
  return request({ url: '/messages/mark-all-read', method: 'post', params: { user_id: userId } })
}

export function deleteMessage(id: number) {
  return request({ url: `/messages/${id}`, method: 'delete' })
}

export function getUnreadCount(userId: number) {
  return request({ url: '/messages/unread-count', method: 'get', params: { user_id: userId } })
}

// ============= 通知渠道 =============

export function getNotificationChannels(params?: { type?: string; page?: number; page_size?: number }) {
  return request({ url: '/notification-channels', method: 'get', params })
}

export function getNotificationChannel(id: number) {
  return request({ url: `/notification-channels/${id}`, method: 'get' })
}

export function createNotificationChannel(data: {
  name: string
  type: string
  config: any
  description?: string
  enabled?: boolean
}) {
  return request({ url: '/notification-channels', method: 'post', data })
}

export function updateNotificationChannel(id: number, data: any) {
  return request({ url: `/notification-channels/${id}`, method: 'put', data })
}

export function deleteNotificationChannel(id: number) {
  return request({ url: `/notification-channels/${id}`, method: 'delete' })
}

export function enableNotificationChannel(id: number) {
  return request({ url: `/notification-channels/${id}/enable`, method: 'post' })
}

export function disableNotificationChannel(id: number) {
  return request({ url: `/notification-channels/${id}/disable`, method: 'post' })
}

export function testNotificationChannel(id: number) {
  return request({ url: `/notification-channels/${id}/test`, method: 'post' })
}

// ============= 机器人 =============

export function getRobots(params?: { type?: string; page?: number; page_size?: number }) {
  return request({ url: '/robots', method: 'get', params })
}

export function getRobot(id: number) {
  return request({ url: `/robots/${id}`, method: 'get' })
}

export function createRobot(data: {
  name: string
  type: string
  webhook_url: string
  secret?: string
  description?: string
  enabled?: boolean
}) {
  return request({ url: '/robots', method: 'post', data })
}

export function updateRobot(id: number, data: any) {
  return request({ url: `/robots/${id}`, method: 'put', data })
}

export function deleteRobot(id: number) {
  return request({ url: `/robots/${id}`, method: 'delete' })
}

export function enableRobot(id: number) {
  return request({ url: `/robots/${id}/enable`, method: 'post' })
}

export function disableRobot(id: number) {
  return request({ url: `/robots/${id}/disable`, method: 'post' })
}

export function testRobot(id: number) {
  return request({ url: `/robots/${id}/test`, method: 'post' })
}

// ============= 接收人 =============

export function getReceivers(params?: { page?: number; page_size?: number }) {
  return request({ url: '/receivers', method: 'get', params })
}

export function getReceiver(id: number) {
  return request({ url: `/receivers/${id}`, method: 'get' })
}

export function createReceiver(data: {
  name: string
  email?: string
  phone?: string
  user_id?: number
  domain_id: number
  enabled?: boolean
}) {
  return request({ url: '/receivers', method: 'post', data })
}

export function updateReceiver(id: number, data: any) {
  return request({ url: `/receivers/${id}`, method: 'put', data })
}

export function deleteReceiver(id: number) {
  return request({ url: `/receivers/${id}`, method: 'delete' })
}

export function enableReceiver(id: number) {
  return request({ url: `/receivers/${id}/enable`, method: 'post' })
}

export function disableReceiver(id: number) {
  return request({ url: `/receivers/${id}/disable`, method: 'post' })
}

export function getReceiverChannels(id: number) {
  return request({ url: `/receivers/${id}/channels`, method: 'get' })
}

export function setReceiverChannels(id: number, data: { channel_ids: number[] }) {
  return request({ url: `/receivers/${id}/channels`, method: 'post', data })
}

export function getReceiverWithChannels(id: number) {
  return request({ url: `/receivers/${id}/with-channels`, method: 'get' })
}

// ============= 消息订阅 =============

export function getSubscriptions(params?: { user_id?: number }) {
  return request({ url: '/subscriptions', method: 'get', params })
}

export function getSubscription(id: number) {
  return request({ url: `/subscriptions/${id}`, method: 'get' })
}

export function createSubscription(data: {
  user_id: number
  message_type_id: number
  email?: boolean
  wechat?: boolean
  dingtalk?: boolean
  webhook?: boolean
  inbox?: boolean
}) {
  return request({ url: '/subscriptions', method: 'post', data })
}

export function updateSubscription(id: number, data: any) {
  return request({ url: `/subscriptions/${id}`, method: 'put', data })
}

export function deleteSubscription(id: number) {
  return request({ url: `/subscriptions/${id}`, method: 'delete' })
}

export function getMessageTypes() {
  return request({ url: '/message-types', method: 'get' })
}
