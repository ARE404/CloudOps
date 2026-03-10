import request from './request'

export interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

export interface KnowledgeItem {
  id: number
  name: string
  size: number
  createdAt: string
}

export const aiApi = {
  chat: (messages: ChatMessage[]): Promise<{ content: string }> =>
    request.post('/api/ai/chat', { messages }),

  listKnowledge: (): Promise<KnowledgeItem[]> =>
    request.get('/api/ai/knowledge'),
  uploadKnowledge: (formData: FormData): Promise<KnowledgeItem> =>
    request.post('/api/ai/knowledge', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
  deleteKnowledge: (id: number): Promise<void> =>
    request.delete(`/api/ai/knowledge/${id}`),
}
