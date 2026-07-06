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
