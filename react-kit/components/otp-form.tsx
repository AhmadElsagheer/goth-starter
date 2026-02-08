'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { CircleAlert, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { otpSchema, OtpFormData } from '@/modules/auth/schemas';

interface OtpFormProps {
    onSubmit: (data: OtpFormData) => void;
    isLoading: boolean;
    identifier: string;
}

export function OtpForm({ onSubmit, isLoading, identifier }: OtpFormProps) {
    const form = useForm<OtpFormData>({
        resolver: zodResolver(otpSchema),
    });

    return (
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 flex flex-col items-start">
            <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Enter the code we sent over SMS to <span className="text-foreground font-medium">{identifier}</span>.</p>
            </div>

            <div className="border border-input rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary focus-within:border-primary transition-all relative bg-background">
                <Input
                    {...form.register('otp')}
                    type="text"
                    maxLength={6}
                    className="border-0 rounded-none focus-visible:ring-0 shadow-none h-auto py-4 px-6 text-base tracking-[0.5em] font-mono max-w-40"
                    placeholder="------"
                />
            </div>
            {form.formState.errors.otp && (
                <p className="text-red-500 text-xs flex items-center gap-1">
                    <CircleAlert className='size-3' />
                    {form.formState.errors.otp.message}
                </p>
            )}

            <div className="flex justify-between text-xs font-medium">
                <button type="button" className="underline hover:text-primary">Resend code</button>
            </div>

            <Button
                type="submit"
                disabled={isLoading}
                className="w-full h-12 text-md font-semibold bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg"
            >
                {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Continue'}
            </Button>
        </form>
    );
}
