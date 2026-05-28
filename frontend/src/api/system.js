import request from '@/utils/request'

// 获取系统状态（设备 + 网络）
export function getSystemStatus() {
  return request({
    url: '/system/status',
    method: 'get'
  })
}

// 获取设备状态（SIP）
export function getDeviceStatus() {
  return request({
    url: '/device/status',
    method: 'get'
  })
}

export function updateSystemDescription(description) {
  return request.post('/system/description', {
    description
  })
}