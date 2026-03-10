import request from './request'

export interface Menu {
  id: number
  parentId: number
  name: string
  path: string
  icon: string
  sort: number
  type: number
  children?: Menu[]
}

export const menuApi = {
  list: (): Promise<Menu[]> =>
    request.get('/api/menus'),
  create: (data: Partial<Menu>): Promise<Menu> =>
    request.post('/api/menus', data),
  update: (id: number, data: Partial<Menu>): Promise<Menu> =>
    request.put(`/api/menus/${id}`, data),
  delete: (id: number): Promise<void> =>
    request.delete(`/api/menus/${id}`),
}
