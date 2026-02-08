import { LoginRequest, LoginResponse, SignupRequest, VerifyEmailRequest } from "./types";

const API_BASE_URL = 'http://localhost:8081'; // TODO(sagheer): move to env

export const checkUser = async (params: { email?: string; phone?: string }): Promise<{ exists: boolean }> => {
    const query = new URLSearchParams();
    if (params.email) query.set('email', params.email);
    if (params.phone) query.set('phoneNumber', params.phone);

    const res = await fetch(`${API_BASE_URL}/auth/check-user?${query.toString()}`);
    if (!res.ok) throw new Error('Failed to check user');
    return res.json();
};

export const signup = async (data: SignupRequest): Promise<void> => {
    const res = await fetch(`${API_BASE_URL}/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!res.ok) {
        const errorBody = await res.text();
        throw new Error(errorBody || 'Signup failed');
    }
};

export const verifyEmail = async (data: VerifyEmailRequest): Promise<void> => {
    const res = await fetch(`${API_BASE_URL}/auth/verify-email`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!res.ok) {
        const errorBody = await res.text();
        throw new Error(errorBody || 'Verification failed');
    }
}

export const login = async (data: LoginRequest): Promise<LoginResponse> => {
    const res = await fetch(`${API_BASE_URL}/_local/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });

    if (!res.ok) {
        const errorBody = await res.text();
        throw new Error(errorBody || 'Login failed');
    }

    return res.json();
};
