import request from './request'

export interface User {
  id: number
  username: string
  realName: string
  email: string
  phone: string
  status: number
  createdAt: string
}

export interface CreateUserParams {
  username: string
  password: string
  realName: string
  email: string
  phone: string
}

export const userApi = {
  list: (params?: Record<string, unknown>): Promise<{ list: User[]; total: number }> =>
    request.get('/api/users', { params }),
  create: (data: CreateUserParams): Promise<User> =>
    request.post('/api/users', data),
  update: (id: number, data: Partial<CreateUserParams>): Promise<User> =>
    request.put(`/api/users/${id}`, data),
  delete: (id: number): Promise<void> =>
    request.delete(`/api/users/${id}`),
}
