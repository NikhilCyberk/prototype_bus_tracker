import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Create axios instance with base configuration
const api = axios.create({
    baseURL: API_URL,
    timeout: 5000, // 5 second timeout
    headers: {
        'Content-Type': 'application/json'
    }
});

// Add request interceptor for adding auth token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);


// Add response interceptor
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.code === 'ERR_NETWORK') {
            console.error('Network error - is the backend running?');
            // You could return mock data here for development
            return Promise.resolve({ data: [] });
        }
        return Promise.reject(error);
    }
);

// Bus related API calls
export const getBuses = () => api.get('/api/buses').then(res => res.data);
export const getBusById = (id) => api.get(`/api/buses/${id}`).then(res => res.data);
export const updateBusLocation = (id, latitude, longitude) =>
    api.put(`/api/buses/${id}/location`, { latitude, longitude }).then(res => res.data);
export const updateBusStatus = (id, status) =>
    api.put(`/api/buses/${id}/status`, { status }).then(res => res.data);

// Route related API calls
export const getRoutes = () => api.get('/api/routes').then(res => res.data);
export const getRouteById = (id) => api.get(`/api/routes/${id}`).then(res => res.data);
export const getRouteStops = (id) => api.get(`/api/routes/${id}/stops`).then(res => res.data);

// Stop related API calls
export const getStops = () => api.get('/api/stops').then(res => res.data);
export const getStopById = (id) => api.get(`/api/stops/${id}`).then(res => res.data);
export const getStopArrivals = (id) => api.get(`/api/stops/${id}/arrivals`).then(res => res.data);

export default api;