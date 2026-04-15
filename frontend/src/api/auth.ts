import request from '@/utils/request'

// Login
export function login(username: string, password: string) {
  return request<{ token: string; user: any }>({
    url: '/auth/login',
    method: 'post',
    data: { username, password }
  })
}

// Get current user info
export function getCurrentUser() {
  return request<any>({
    url: '/auth/me',
    method: 'get'
  })
}

// Change password
export function changePassword(oldPassword: string, newPassword: string) {
  return request<{ message: string }>({
    url: '/auth/change-password',
    method: 'post',
    data: { old_password: oldPassword, new_password: newPassword }
  })
}

// Update profile
export function updateProfile(data: { display_name?: string; email?: string; phone?: string; remark?: string }) {
  return request<{ message: string; user: any }>({
    url: '/auth/profile',
    method: 'put',
    data
  })
}