'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { CircleAlert, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { emailSchema, EmailFormData } from '@/modules/auth/schemas';

interface EmailFormProps {
    onSubmit: (data: EmailFormData) => void;
    isLoading: boolean;
}

export function EmailForm({ onSubmit, isLoading }: EmailFormProps) {
    const form = useForm<EmailFormData>({
        resolver: zodResolver(emailSchema),
    });

    return (
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="border border-input rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary focus-within:border-primary transition-all relative bg-background">
                <Label className="absolute top-2 left-3 text-xs text-muted-foreground z-10 pointer-events-none">Email</Label>
                <Input
                    {...form.register('email')}
                    type="email"
                    className="border-0 rounded-none focus-visible:ring-0 shadow-none h-auto pt-6 pb-2 px-3 text-base"
                    placeholder=" "
                />
            </div>
            {form.formState.errors.email && (
                <p className="text-red-500 text-xs flex items-center gap-1">
                    <CircleAlert className='size-3' />
                    {form.formState.errors.email.message}
                </p>
            )}

            <Button
                type="submit"
                disabled={isLoading}
                className="w-full h-12 text-md font-semibold bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg mt-2"
            >
                {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Continue'}
            </Button>
        </form>
    );
}
