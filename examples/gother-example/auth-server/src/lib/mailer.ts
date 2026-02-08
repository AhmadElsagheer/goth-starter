import * as nodemailer from "nodemailer";

interface MailerConfig {
    user: string;
    pass: string;
}

export class Mailer {
    private transporter: nodemailer.Transporter;
    private readonly from = "Pulse <noreply@pulse.app>";

    constructor(config: MailerConfig) {
        if (!config.user || !config.pass) {
            throw new Error("invalid mailer config");
        }

        this.transporter = nodemailer.createTransport({
            host: "smtp.gmail.com",
            port: 587,
            secure: false, // true for 465, false for other ports. TODO(sagheer): revist this
            auth: {
                user: config.user,
                pass: config.pass,
            },
        });
    }

    async sendEmail(to: string, subject: string, body: string): Promise<void> {
        await this.transporter.sendMail({
            from: this.from,
            to,
            subject,
            html: body,
        });
    }
}

export const createMailer = (config: MailerConfig) => {
    return new Mailer(config);
};
