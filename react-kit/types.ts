export interface SignupRequest {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    countryCode: string;
    phoneNumber: string;
    contactMethods: string[];
}

export interface VerifyEmailRequest {
    email: string;
    otp: string;
    goal: 'Signup' | 'ResetPassword';
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
}
