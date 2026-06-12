import request from '@/utils/request'

export function getInventory(params) {
  return request({
    url: '/inventory',
    method: 'get',
    params
  })
}

export function updateInventory(id, data) {
  return request({
    url: `/inventory/${id}`,
    method: 'put',
    data
  })
}