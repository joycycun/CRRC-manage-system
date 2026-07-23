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

export function submitInventoryScrap(id) {
  return request({
    url: `/inventory/${id}/scrap-submit`,
    method: 'post'
  })
}

export function auditInventoryScrap(id, data) {
  return request({
    url: `/inventory/${id}/scrap-audit`,
    method: 'post',
    data
  })
}

export function returnInventoryTypeToBoardInbound(data) {
  return request({
    url: '/inventory/return-to-board-inbound',
    method: 'post',
    data
  })
}
