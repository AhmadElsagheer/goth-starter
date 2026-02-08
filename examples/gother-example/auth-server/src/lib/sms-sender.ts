import twilio from "twilio";

interface SmsSenderConfig {
    accountSid: string;
    authToken: string;
    phoneNumber: string;
}

export class SmsSender {
    private client: twilio.Twilio;
    private from: string;

    constructor(config: SmsSenderConfig) {
        this.client = twilio(config.accountSid, config.authToken);
        this.from = config.phoneNumber;
    }

    async sendSms(to: string, body: string): Promise<void> {
        await this.client.messages.create({
            body,
            from: this.from,
            to,
        });
    }
}

export const createSmsSender = (config: SmsSenderConfig) => {
    return new SmsSender(config);
};
