import React, { createContext, useState, useContext, useEffect } from "react";
import { login as apiLogin, register as apiRegister } from "../services/auth";

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check if user is logged in by token in localStorage
        const token = localStorage.getItem("token");
        const userData = localStorage.getItem("user");

        if (token && userData) {
            setUser(JSON.parse(userData));
        }

        setLoading(false);
    }, []);

    const login = async (username, password) => {
        try {
            const { token, user } = await apiLogin(username, password);

            localStorage.setItem("token", token);
            localStorage.setItem("user", JSON.stringify(user));

            setUser(user);
            return { success: true };
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || "Login failed",
            };
        }
    };

    const register = async (username, password, role) => {
        try {
            const newUser = await apiRegister(username, password, role);
            return { success: true, user: newUser };
        } catch (error) {
            return {
                success: false,
                error: error.response?.data?.error || "Registration failed",
            };
        }
    };

    const logout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        setUser(null);
    };

    const value = {
        user,
        loading,
        login,
        register,
        logout,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
