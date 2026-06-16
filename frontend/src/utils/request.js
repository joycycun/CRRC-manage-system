import axios from 'axios'
import { getStoredRoles } from '@/utils/permission'

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api',      
  timeout: 10000
})

// 请求拦截器（自动带 token）
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = 'Bearer ' + token
    }

    const userText = localStorage.getItem('user')
    if (userText) {
      try {
        const user = JSON.parse(userText)
        config.headers['X-User-Id'] = user.id || ''
        config.headers['X-User-Name'] = encodeURIComponent(user.realName || user.username || '')
      } catch (err) {
        console.warn('读取请求用户信息失败：', err)
      }
    }

    config.headers['X-User-Roles'] = getStoredRoles().join(',')

    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器
request.interceptors.response.use(
  response => response,
  error => {
    console.error('request error:', error)
    return Promise.reject(error)
  }
)

export default request
