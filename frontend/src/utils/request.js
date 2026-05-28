import axios from 'axios'

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
