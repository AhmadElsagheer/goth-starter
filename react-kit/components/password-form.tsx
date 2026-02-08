'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { CircleAlert, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { passwordSchema, PasswordFormData } from '@/modules/auth/schemas';

interface PasswordFormProps {
    onSubmit: (data: PasswordFormData) => void;
    isLoading: boolean;
    identifier: string;
}

export function PasswordForm({ onSubmit, isLoading, identifier }: PasswordFormProps) {
    const form = useForm<PasswordFormData>({
        resolver: zodResolver(passwordSchema),
    });

    return (
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-2 pb-4">
                <h2 className="text-2xl font-semibold tracking-tight">Welcome back</h2>
                <p className="text-sm text-muted-foreground">
                    Enter your password below to access your account.
                </p>
                <div className="mt-2 inline-flex items-center px-4 py-1.5 rounded-full bg-secondary/50 text-secondary-foreground text-sm font-medium">
                    {identifier}
                </div>
            </div>

            <div className="border border-input rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary focus-within:border-primary transition-all relative bg-background">
                <Label className="absolute top-2 left-3 text-xs text-muted-foreground z-10 pointer-events-none">Password</Label>
                <Input
                    {...form.register('password')}
                    type="password"
                    className="border-0 rounded-none focus-visible:ring-0 shadow-none h-auto pt-6 pb-2 px-3 text-base"
                    placeholder=" "
                />
            </div>
            {form.formState.errors.password && (
                <p className="text-red-500 text-xs flex items-center gap-1">
                    <CircleAlert className='size-3' />
                    {form.formState.errors.password.message}
                </p>
            )}

            <Button
                type="submit"
                disabled={isLoading}
                className="w-full h-12 text-md font-semibold bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg"
            >
                {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Log in'}
            </Button>

            <div className="text-center">
                <a href="#" className="text-sm font-medium underline hover:text-primary transition-colors">Forgot password?</a>
            </div>
        </form>
    );
}
