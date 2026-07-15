import request from '@/utils/request'

export function getBoardCompositions(params) {
  return request({
    url: '/board-compositions',
    method: 'get',
    params
  })
}

export function createBoardCompositions(data) {
  return request({
    url: '/board-compositions',
    method: 'post',
    data
  })
}

export function deleteBoardComposition(id) {
  return request({
    url: `/board-compositions/${id}`,
    method: 'delete'
  })
}
