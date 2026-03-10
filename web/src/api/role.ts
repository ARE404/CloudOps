import request from './request'

export interface Role {
  id: number
  name: string
  description: string
  createdAt: string
}

export const roleApi = {
  list: (params?: Record<string, unknown>): Promise<{ list: Role[]; total: number }> =>
    request.get('/api/roles', { params }),
  create: (data: Partial<Role>): Promise<Role> =>
    request.post('/api/roles', data),
  update: (id: number, data: Partial<Role>): Promise<Role> =>
    request.put(`/api/roles/${id}`, data),
  delete: (id: number): Promise<void> =>
    request.delete(`/api/roles/${id}`),
}
