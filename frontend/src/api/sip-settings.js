import request from '@/utils/request'

export function setSipAccount(data) {
  return request({
    url: '/sip-settings/account',
    method: 'post',
    data
  })
}

export function getSipAccount() {
  return request({
    url: '/sip-settings/account',
    method: 'get'
  })
}

export function getSipCodec() {
  return request({
    url: '/sip-settings/audio-codes',
    method: 'get'
  })
}

export function setSipCodec(data) {
  return request({
    url: '/sip-settings/audio-codes',
    method: 'post',
    data
  })
}

export function setSipAdvanceInfo(data) {
  return request({
    url: '/sip-settings/sip-advance-settings',
    method: 'post',
    data
  })
}

export function getSipAdvanceInfo() {
  return request({
    url: '/sip-settings/sip-advance-settings',
    method: 'get'
  })
}

