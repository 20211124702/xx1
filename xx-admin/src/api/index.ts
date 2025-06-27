import axios from 'axios'

const service = axios.create({
  baseURL: '/api',
  timeout: 5000
})

service.interceptors.request.use(config => {
  // 可添加 token 等
  return config
})

service.interceptors.response.use(
  response => response.data,
  error => Promise.reject(error)
)

export default service 