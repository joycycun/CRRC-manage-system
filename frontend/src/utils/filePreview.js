export function createLocalFileId() {
  const random = Math.floor(Math.random() * 100000)
  return Number(`${Date.now()}${String(random).padStart(5, '0')}`)
}

export function readFileAsDataURL(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result || ''))
    reader.onerror = () => reject(reader.error || new Error('读取文件失败'))
    reader.readAsDataURL(file)
  })
}

export async function buildUploadFilePayload(file, existingFileId = 0) {
  if (!file) return {}
  return {
    fileId: existingFileId || createLocalFileId(),
    fileName: file.name,
    fileContentType: file.type || 'application/octet-stream',
    fileData: await readFileAsDataURL(file),
    fileUrl: URL.createObjectURL(file)
  }
}

export function getFilePreviewUrl(fileId) {
  return fileId ? `/api/files/${fileId}/preview` : ''
}

export function getFileDownloadUrl(fileId) {
  return fileId ? `/api/files/${fileId}/download` : ''
}

export function openLocalFilePreview(item) {
  const url = item?.fileUrl || getFilePreviewUrl(item?.fileId)
  if (!url) {
    alert('当前文件暂无可预览内容')
    return
  }
  window.open(url, '_blank')
}

export function downloadLocalFile(item, fallbackName = '文件') {
  const url = item?.downloadUrl || item?.fileUrl || getFileDownloadUrl(item?.fileId)
  if (!url) {
    alert('当前文件暂无可下载内容')
    return
  }

  const link = document.createElement('a')
  link.href = url
  link.download = item?.fileName || fallbackName
  link.click()
}
