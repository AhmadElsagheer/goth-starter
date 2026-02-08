'use client';

import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { CircleAlert, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { phoneSchema, PhoneFormData } from '@/modules/auth/schemas';
import { COUNTRIES } from '@/lib/static';

interface PhoneFormProps {
    onSubmit: (data: PhoneFormData) => void;
    isLoading: boolean;
}

export function PhoneForm({ onSubmit, isLoading }: PhoneFormProps) {
    const form = useForm<PhoneFormData>({
        resolver: zodResolver(phoneSchema),
        defaultValues: {
            countryCode: '+20',
        },
    });

    return (
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="border border-input rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary focus-within:border-primary transition-all">
                <div className="border-b border-input relative">
                    <Label className="absolute top-2 left-3 text-xs text-muted-foreground z-10 w-full cursor-text pointer-events-none">Country Code</Label>
                    <Controller
                        control={form.control}
                        name="countryCode"
                        render={({ field }) => (
                            <Select onValueChange={field.onChange} defaultValue={field.value}>
                                <SelectTrigger className="w-full border-0 rounded-none h-auto pt-10 pb-4 px-3 focus:ring-0 shadow-none font-normal [&_svg:not([class*='text-'])]:text-foreground [&_svg]:opacity-100">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {COUNTRIES.map((c) => (
                                        <SelectItem key={c.value} value={c.value}>
                                            <div className="flex justify-between w-full min-w-[200px]">
                                                <span>{c.label.split(' (')[0]} ({c.value})</span>
                                            </div>
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        )}
                    />
                </div>
                <div className="relative bg-background">
                    <Label className="absolute top-2 left-3 text-xs text-muted-foreground z-10 pointer-events-none">Phone number</Label>
                    <Input
                        {...form.register('phoneNumber')}
                        type="tel"
                        className="border-0 rounded-none focus-visible:ring-0 shadow-none h-auto pt-6 pb-2 px-3 text-base"
                        placeholder=" "
                    />
                </div>
            </div>
            {form.formState.errors.phoneNumber && (
                <p className="text-red-500 text-xs flex items-center gap-1">
                    <CircleAlert className='size-3' />
                    {form.formState.errors.phoneNumber.message}
                </p>
            )}

            <div className="text-xs text-muted-foreground">
                We&apos;ll text you to confirm your number. <a href="#" className="underline">Privacy Policy</a>
            </div>

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
