import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const createTask = async (title: string, description: string) => {
    const token = localStorage.getItem("token"); // ✅ Get JWT token
    if (!token) {
        throw new Error("User is not authenticated"); // ✅ Prevent unauthorized requests
    }

    return axios.post(`${API_BASE_URL}/tasks`, { title, description }, {
        headers: { 
            Authorization: `Bearer ${token}`, // ✅ Send token in request
            "Content-Type": "application/json"
        }
    });
};
