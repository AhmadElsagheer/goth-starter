import { z } from 'zod';

const envSchema = z.object({
    SERVER_PORT: z.coerce.number(),
    DATABASE_URL: z.string(),
    BETTER_AUTH_URL: z.url(),
    BETTER_AUTH_TRUSTED_ORIGINS: z.string().transform((str) => str.split(",")),

    BETTER_AUTH_SECRET: z.string().min(32),
    AUTH_METHODS: z.string().transform((str) => str.split(",")),

    SMTP_USER: z.email().optional(),
    SMTP_PASSWORD: z.string().optional(),

    TWILIO_ACCOUNT_SID: z.string().optional(),
    TWILIO_AUTH_TOKEN: z.string().optional(),
    TWILIO_PHONE_NUMBER: z.string().optional(),

    GOOGLE_CLIENT_ID: z.string().optional(),
    GOOGLE_CLIENT_SECRET: z.string().optional(),
}).superRefine((data, ctx) => {
    if (data.AUTH_METHODS.includes("email")) {
        if (!data.SMTP_USER) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: "SMTP_USER is required for email auth",
                path: ["SMTP_USER"],
            });
        }
        if (!data.SMTP_PASSWORD) {
            ctx.addIssue({
                code: "custom",
                message: "SMTP_PASSWORD is required for email auth",
                path: ["SMTP_PASSWORD"],
            });
        }
    }
    if (data.AUTH_METHODS.includes("phone")) {
        if (!data.TWILIO_ACCOUNT_SID) {
            ctx.addIssue({
                code: "custom",
                message: "TWILIO_ACCOUNT_SID is required for phone auth",
                path: ["TWILIO_ACCOUNT_SID"],
            });
        }
        if (!data.TWILIO_AUTH_TOKEN) {
            ctx.addIssue({
                code: "custom",
                message: "TWILIO_AUTH_TOKEN is required for phone auth",
                path: ["TWILIO_AUTH_TOKEN"],
            });
        }
        if (!data.TWILIO_PHONE_NUMBER) {
            ctx.addIssue({
                code: "custom",
                message: "TWILIO_PHONE_NUMBER is required for phone auth",
                path: ["TWILIO_PHONE_NUMBER"],
            });
        }
    }
    if (data.AUTH_METHODS.includes("google")) {
        if (!data.GOOGLE_CLIENT_ID) {
            ctx.addIssue({
                code: "custom",
                message: "GOOGLE_CLIENT_ID is required for google auth",
                path: ["GOOGLE_CLIENT_ID"],
            });
        }
        if (!data.GOOGLE_CLIENT_SECRET) {
            ctx.addIssue({
                code: "custom",
                message: "GOOGLE_CLIENT_SECRET is required for google auth",
                path: ["GOOGLE_CLIENT_SECRET"],
            });
        }
    }
});

export const env = envSchema.parse(process.env);
