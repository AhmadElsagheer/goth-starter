import { createAuthEndpoint } from "better-auth/api";
import { phoneNumber } from "better-auth/plugins";
import { BetterAuthPlugin } from "better-auth/types";
import { email, z } from "zod";

export const checkUser = () => {
    return {
        id: "check-user",
        endpoints: {
            checkUser: createAuthEndpoint("/check-user/check", {
                method: "POST",
                body: z.object({
                    email: z.string().optional(),
                    phoneNumber: z.string().optional()
                }),
            }, async (ctx) => {
                const { email, phoneNumber } = ctx.body;
                const { internalAdapter, adapter } = ctx.context

                if (email) {
                    let user = await internalAdapter.findUserByEmail(email);
                    if (user) return ctx.json({ exists: true });
                }

                if (phoneNumber) {
                    let user = await adapter.findOne({
                        model: "user",
                        where: [
                            {
                                field: "phoneNumber",
                                value: phoneNumber,
                            }
                        ]
                    })
                    if (user) return ctx.json({ exists: true });
                }

                return ctx.json({
                    exists: false
                });
            })
        }
    } satisfies BetterAuthPlugin
}
