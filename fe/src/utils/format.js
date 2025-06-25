export const getStatusType = (status) => {
  switch (status) {
    case 0: return 'info'    // 待处理
    case 1: return 'warning' // 处理中
    case 2: return 'success' // 已完成
    case 3: return 'danger'  // 失败
    default: return 'info'
  }
}

export const getStatusText = (status) => {
  switch (status) {
    case 0: return '待处理'
    case 1: return '处理中'
    case 2: return '已完成'
    case 3: return '失败'
    default: return '未知'
  }
}

export const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}