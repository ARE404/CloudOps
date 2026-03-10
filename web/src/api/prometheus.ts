import request from './request'

export interface ScrapePool {
  id: number
  name: string
  clusterId: number
  labels: string
  createdAt: string
}

export interface AlertRule {
  id: number
  name: string
  expr: string
  severity: string
  summary: string
  createdAt: string
}

export interface AlertEvent {
  id: number
  alertName: string
  severity: string
  status: string
  summary: string
  startAt: string
}

export interface SendGroup {
  id: number
  name: string
  type: string
  webhook: string
  createdAt: string
}

export const prometheusApi = {
  listScrapePools: (params?: Record<string, unknown>): Promise<{ list: ScrapePool[]; total: number }> =>
    request.get('/api/prometheus/scrape-pools', { params }),
  createScrapePool: (data: Partial<ScrapePool>): Promise<ScrapePool> =>
    request.post('/api/prometheus/scrape-pools', data),
  updateScrapePool: (id: number, data: Partial<ScrapePool>): Promise<ScrapePool> =>
    request.put(`/api/prometheus/scrape-pools/${id}`, data),
  deleteScrapePool: (id: number): Promise<void> =>
    request.delete(`/api/prometheus/scrape-pools/${id}`),

  listAlertRules: (params?: Record<string, unknown>): Promise<{ list: AlertRule[]; total: number }> =>
    request.get('/api/prometheus/alert-rules', { params }),
  createAlertRule: (data: Partial<AlertRule>): Promise<AlertRule> =>
    request.post('/api/prometheus/alert-rules', data),
  updateAlertRule: (id: number, data: Partial<AlertRule>): Promise<AlertRule> =>
    request.put(`/api/prometheus/alert-rules/${id}`, data),
  deleteAlertRule: (id: number): Promise<void> =>
    request.delete(`/api/prometheus/alert-rules/${id}`),

  listAlertEvents: (params?: Record<string, unknown>): Promise<{ list: AlertEvent[]; total: number }> =>
    request.get('/api/prometheus/alert-events', { params }),

  listSendGroups: (params?: Record<string, unknown>): Promise<{ list: SendGroup[]; total: number }> =>
    request.get('/api/prometheus/send-groups', { params }),
  createSendGroup: (data: Partial<SendGroup>): Promise<SendGroup> =>
    request.post('/api/prometheus/send-groups', data),
  updateSendGroup: (id: number, data: Partial<SendGroup>): Promise<SendGroup> =>
    request.put(`/api/prometheus/send-groups/${id}`, data),
  deleteSendGroup: (id: number): Promise<void> =>
    request.delete(`/api/prometheus/send-groups/${id}`),
}
