// 文本处理工具

/**
 * 去除 Markdown 语法标记，返回纯文本
 */
export function stripMarkdown(md) {
  if (!md) return ''
  return md
    // 代码块
    .replace(/```[\s\S]*?```/g, '')
    // 行内代码
    .replace(/`([^`]+)`/g, '$1')
    // 图片
    .replace(/!\[.*?\]\(.*?\)/g, '')
    // 链接
    .replace(/\[([^\]]+)\]\(.*?\)/g, '$1')
    // 标题
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/\s+#{1,6}\s+/g, ' ')
    // 粗体/斜体
    .replace(/\*\*(.+?)\*\*/g, '$1')
    .replace(/\*(.+?)\*/g, '$1')
    .replace(/__(.+?)__/g, '$1')
    .replace(/_(.+?)_/g, '$1')
    // 删除线
    .replace(/~~(.+?)~~/g, '$1')
    // 引用
    .replace(/^>\s+/gm, '')
    // 无序列表
    .replace(/^[\s]*[-*+]\s+/gm, '')
    // 有序列表
    .replace(/^[\s]*\d+\.\s+/gm, '')
    // 水平线
    .replace(/^[-*_]{3,}\s*$/gm, '')
    // HTML 标签
    .replace(/<[^>]+>/g, '')
    // 多余空行
    .replace(/\n{2,}/g, '\n')
    .trim()
}

/**
 * 截断文本到指定长度
 */
export function truncate(text, maxLen = 50) {
  if (!text) return ''
  const clean = stripMarkdown(text)
  if (clean.length <= maxLen) return clean
  return clean.slice(0, maxLen) + '…'
}
