export const EMAIL_VERIFICATION_SUBJECT = "Your Email Verification OTP";
export const EMAIL_VERIFICATION_BODY = (otp: string) => `<p>Your OTP is <strong>${otp}</strong></p>`;
export const SMS_OTP_BODY = (code: string) => `Your OTP is ${code}`;
