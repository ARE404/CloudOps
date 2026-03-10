import request from './request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  refresh_token: string
}

export const authApi = {
  login: (data: LoginParams): Promise<LoginResult> =>
    request.post('/api/auth/login', data),
  logout: (): Promise<void> =>
    request.post('/api/auth/logout'),
}
