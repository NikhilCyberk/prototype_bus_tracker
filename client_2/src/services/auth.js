//File : src/services/auth.js
import api from './api';

export const login = async (username, password) => {
    const response = await api.post('/api/auth/login', { username, password });
    return response.data;
};

export const register = async (username, password, role) => {
    const response = await api.post('/api/auth/register', { username, password, role });
    return response.data;
};