export interface ApiResponse<T = unknown> {
  code: number
  data: T
  message: string
}

export interface PageResult<T> {
  list: T[]
  total: number
}

export interface PageParams {
  page: number
  pageSize: number
}
