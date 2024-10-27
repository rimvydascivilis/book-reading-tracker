import axios from "axios";
import {getToken} from "../service/TokenService";

const api = axios.create({
  baseURL: "http://localhost:8081/api",
  timeout: 10000, // 10 seconds
});

api.interceptors.request.use(
  config => {
    const token = getToken();
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  },
);

export default api;
