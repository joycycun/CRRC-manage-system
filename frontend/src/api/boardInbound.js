import request from '@/utils/request'

export function getBoardInboundRecords(params) {
  return request({
    url: '/board-inbound',
    method: 'get',
    params
  })
}

export function importBoardInboundRecords(data) {
  return request({
    url: '/board-inbound',
    method: 'post',
    data
  })
}

export function updateBoardInboundRecord(data) {
  return request({
    url: '/board-inbound',
    method: 'put',
    data
  })
}

export function deleteBoardInboundRecord(id) {
  return request({
    url: '/board-inbound',
    method: 'delete',
    params: { id }
  })
}
