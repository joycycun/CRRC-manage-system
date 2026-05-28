import request from '@/utils/request'

export function GetAudioInfo() {
  return request({
    url: '/system-settings/audio',
    method: 'get'
  })
}

export function getVolume() {
  return request({
    url: '/system-settings/volume-control',
    method: 'get'
  })
}

export function setVolume(data) {
  return request({
    url: '/system-settings/volume-control',
    method: 'post',
    data
  })
}

export function setIntercom(data) {
  return request({
    url: 'system-settings/intercom',
    method: 'post',
    data
  })
}

export function getIntercom() {
  return request({
    url: 'system-settings/intercom',
    method: 'get'
  })
}

