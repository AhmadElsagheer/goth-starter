import * as z from 'zod';
import { isValidPhoneNumber } from 'libphonenumber-js';

export const emailSchema = z.object({
    email: z.email('Invalid email address'),
});

export const phoneSchema = z.object({
    countryCode: z.string().min(1, 'Country code is required'),
    phoneNumber: z.string().min(1, 'Phone number is required'),
}).superRefine((data, ctx) => {
    const fullNumber = `${data.countryCode}${data.phoneNumber}`;
    if (!isValidPhoneNumber(fullNumber)) {
        ctx.addIssue({
            code: "custom",
            message: 'Invalid phone number',
            path: ['phoneNumber'],
        });
    }
});

export const passwordSchema = z.object({
    password: z.string().min(1, 'Password is required'),
});

export const otpSchema = z.object({
    otp: z.string().length(6, 'Code must be 6 digits'),
});

export type EmailFormData = z.infer<typeof emailSchema>;
export type PhoneFormData = z.infer<typeof phoneSchema>;
export type PasswordFormData = z.infer<typeof passwordSchema>;
export type OtpFormData = z.infer<typeof otpSchema>;
