import request from './index'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    email: string
    role: {
      id: number
      name: string
    }
  }
}

// 登录
export function login(data: LoginRequest) {
  return request<LoginResponse>({
    url: '/auth/login',
    method: 'POST',
    data
  })
}

// 登出
export function logout() {
  return request({
    url: '/auth/logout',
    method: 'POST'
  })
}

// 获取用户资料
export function getProfile() {
  return request({
    url: '/auth/profile',
    method: 'GET'
  })
}

// 注册
export function register(data: { username: string; password: string; email: string }) {
  return request({
    url: '/auth/register',
    method: 'POST',
    data
  })
} 