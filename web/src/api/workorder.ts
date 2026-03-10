import request from './request'

export interface WorkorderTemplate {
  id: number
  name: string
  description: string
  steps: string
  createdAt: string
}

export interface WorkorderInstance {
  id: number
  templateId: number
  templateName: string
  title: string
  status: string
  creatorName: string
  createdAt: string
}

export interface WorkorderFlow {
  id: number
  instanceId: number
  action: string
  comment: string
  operatorName: string
  createdAt: string
}

export const workorderApi = {
  listTemplates: (params?: Record<string, unknown>): Promise<{ list: WorkorderTemplate[]; total: number }> =>
    request.get('/api/workorder/templates', { params }),
  createTemplate: (data: Partial<WorkorderTemplate>): Promise<WorkorderTemplate> =>
    request.post('/api/workorder/templates', data),
  updateTemplate: (id: number, data: Partial<WorkorderTemplate>): Promise<WorkorderTemplate> =>
    request.put(`/api/workorder/templates/${id}`, data),
  deleteTemplate: (id: number): Promise<void> =>
    request.delete(`/api/workorder/templates/${id}`),

  listInstances: (params?: Record<string, unknown>): Promise<{ list: WorkorderInstance[]; total: number }> =>
    request.get('/api/workorder/instances', { params }),
  createInstance: (data: Partial<WorkorderInstance>): Promise<WorkorderInstance> =>
    request.post('/api/workorder/instances', data),

  listFlows: (instanceId: number): Promise<WorkorderFlow[]> =>
    request.get(`/api/workorder/instances/${instanceId}/flows`),
  doAction: (instanceId: number, data: { action: string; comment: string }): Promise<void> =>
    request.post(`/api/workorder/instances/${instanceId}/action`, data),
}
