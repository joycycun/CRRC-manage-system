import request from '@/utils/request'

export function getUsageGuides() {
  return request.get('/usage-guides')
}

export function createUsageGuide(data) {
  return request.post('/usage-guides', data, { timeout: 60000 })
}

export function deleteUsageGuide(id) {
  return request.delete(`/usage-guides/${id}`)
}
