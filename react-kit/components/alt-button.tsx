import { Button } from "@/components/ui/button";

interface AltButtonProps {
    icon: React.ReactNode;
    text: string;
    onClick: () => void;
}

export function AltButton({ icon, text, onClick }: AltButtonProps) {
    return (
        <Button variant="outline" className="w-full h-12 justify-start pl-4 relative text-foreground border-input hover:bg-secondary/50 rounded-lg" onClick={onClick}>
            <div className="absolute left-6 top-1/2 -translate-y-1/2 w-5 h-5">
                {icon}
            </div>
            <span className="w-full text-center font-medium">{text}</span>
        </Button>
    )
}
