import axios from "axios";

const API_BASE_URL = "http://localhost:8080";

export const registerUser = async (username: string, password: string) => {
    return axios.post(`${API_BASE_URL}/register`, { username, password });
};

export const loginUser = async (username: string, password: string) => {
    const response = await axios.post(`${API_BASE_URL}/login`, { username, password });
    
    // ✅ Store token and role in localStorage
    localStorage.setItem("token", response.data.token);
    localStorage.setItem("role", response.data.role);

    return response.data;
};

export const logoutUser = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
};
