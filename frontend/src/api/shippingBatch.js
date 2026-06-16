import request from '@/utils/request'

export function getShippingBatches(params) {
  return request({
    url: '/shipping-batches',
    method: 'get',
    params
  })
}

export function createShippingBatch(data) {
  return request({
    url: '/shipping-batches',
    method: 'post',
    data
  })
}

export function submitShippingBatch(id) {
  return request({
    url: `/shipping-batches/${id}/submit`,
    method: 'post'
  })
}

export function auditShippingBatch(id, data) {
  return request({
    url: `/shipping-batches/${id}/audit`,
    method: 'post',
    data
  })
}

export function deleteShippingBatch(id) {
  return request({
    url: `/shipping-batches/${id}`,
    method: 'delete'
  })
}

export function createProductionRequest(data) {
  return request({
    url: '/production-requests',
    method: 'post',
    data
  })
}

export function confirmProductionRequest(id, data) {
  return request({
    url: `/production-requests/${id}/confirm`,
    method: 'post',
    data
  })
}
