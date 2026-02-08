import { betterAuth, MiddlewareInputContext, MiddlewareOptions } from "better-auth";
import { BetterAuthPlugin } from "better-auth/types";
import { Pool } from "pg";
import { jwt, emailOTP, phoneNumber, bearer } from "better-auth/plugins";
import { createMailer } from "./mailer";
import { createSmsSender } from "./sms-sender";
import { checkUser } from "./check-user";
import { env } from "../config";
import { EMAIL_VERIFICATION_BODY, EMAIL_VERIFICATION_SUBJECT, SMS_OTP_BODY } from "./templates";

const plugins: BetterAuthPlugin[] = [
    // Authorization header
    bearer(),
    // Use JWT tokens for auth
    jwt(),
    // Endpoint to check if user exists by using email or phone
    checkUser()
];


if (env.AUTH_METHODS.includes("email")) {
    const mailer = createMailer({
        user: env.SMTP_USER!,
        pass: env.SMTP_PASSWORD!,
    });

    plugins.push(
        emailOTP({
            overrideDefaultEmailVerification: true,
            sendVerificationOnSignUp: true,
            async sendVerificationOTP({ email, otp, type }) {
                if (type === "email-verification") {
                    await mailer.sendEmail(
                        email,
                        EMAIL_VERIFICATION_SUBJECT,
                        EMAIL_VERIFICATION_BODY(otp)
                    );
                }
            },
        })
    );
}

if (env.AUTH_METHODS.includes("phone")) {
    const smsSender = createSmsSender({
        accountSid: env.TWILIO_ACCOUNT_SID!,
        authToken: env.TWILIO_AUTH_TOKEN!,
        phoneNumber: env.TWILIO_PHONE_NUMBER!,
    });

    plugins.push(
        phoneNumber({
            sendOTP: async ({ phoneNumber, code }) => {
                await smsSender.sendSms(phoneNumber, SMS_OTP_BODY(code));
            },
            signUpOnVerification: {
                getTempEmail: (phoneNumber) => {
                    return `${phoneNumber}@app.internal`;
                },
                getTempName: (phoneNumber) => {
                    return phoneNumber;
                },
            },
        })
    );
}

const socialProviders: any = {};
if (env.AUTH_METHODS.includes("google")) {
    socialProviders.google = {
        clientId: env.GOOGLE_CLIENT_ID!,
        clientSecret: env.GOOGLE_CLIENT_SECRET!,
    };
}

export const auth = betterAuth({
    baseURL: env.BETTER_AUTH_URL,
    trustedOrigins: env.BETTER_AUTH_TRUSTED_ORIGINS,
    database: new Pool({
        connectionString: env.DATABASE_URL,
    }),
    socialProviders,
    emailAndPassword: {
        enabled: env.AUTH_METHODS.includes("email"),
        requireEmailVerification: true,
    },
    plugins,
    user: {
        modelName: "users",
        additionalFields: {
            roles: {
                type: "string[]",
                required: true,
                defaultValue: ["customer"],
                input: false,
            },
        },
    },
    session: {
        modelName: "sessions",
    },
    account: {
        modelName: "accounts",
    },
    verification: {
        modelName: "verifications",
    },
    advanced: {
        database: {
            generateId: "uuid",
        },
    },
});
