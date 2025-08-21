import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Icons } from "@/icons"
import { Link, useNavigate } from "react-router-dom"
import { FormEvent, useCallback, useState } from "react"
import axios from "axios"
import { useAuth } from "@/auth"
import { usePageTitle } from "@/hooks"

export function Register({
    className,
    ...props
}: React.ComponentProps<"div">) {
    usePageTitle("Signup");
    const [loading, setLoading] = useState<boolean>();
    const [error, setError] = useState<string>()
    const {loadUserInfo} = useAuth();

    const navigate = useNavigate();
    const handleSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        setLoading(true);
        axios.postForm("/api/auth/register", formData)
            .then((data) => {
                setLoading(false);
                navigate("/"); // frontend-reroute
                loadUserInfo?.();
            }).catch((error) => {
                setLoading(false);
                setError(error.rsponse?.data);
            });
    }, [loadUserInfo, navigate]);

    return (
        <div className={cn("flex flex-col gap-6", className)} {...props}>
            <Card className="bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm">
                <CardHeader className="text-center">
                    <CardTitle>Create an account</CardTitle>
                </CardHeader>
                <CardContent className="px-6 flex flex-col gap-4">
                    {/* <div className="flex flex-col gap-4">
                        <div className="text-muted-foreground text-sm">
                            Creat account using
                        </div>
                        <div className="grid grid-cols-2 gap-6">
                            <Button variant="outline" className="">
                                    <Icons.GitHub />
                                    Github
                            </Button>
                            <Button variant="outline" className="">
                                    <Icons.Google />
                                    Google
                            </Button>
                        </div>
                    </div>
                    <div className="relative">
                        <div className="absolute inset-0 flex items-center"><span className="w-full border-t"></span></div>
                        <div className="relative flex justify-center text-xs uppercase">
                            <span className="bg-card text-muted-foreground px-2">or continue with email</span>
                        </div>
                    </div> */}
                    <form onSubmit={handleSubmit}>
                        <div className="flex flex-col gap-6">
                            <div className="grid gap-3">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    name="email"
                                    type="email"
                                    placeholder="m@example.com"
                                    required
                                />
                            </div>
                            <div className="grid gap-3">
                                <div className="flex items-center">
                                    <Label htmlFor="password">Password</Label>
                                </div>
                                <Input id="password" name="password" type="password" autoComplete="password" required />
                            </div>
                            {error ? <p className="text-sm">{error}</p> : null}
                            <Button type="submit" variant="default" className="w-full">
                                Create account
                            </Button>
                        </div>
                        <div className="mt-4 text-center text-sm">
                            Already have an account?{" "}
                            <Link to={"/login"} className="underline underline-offset-4">Log in</Link>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    )
}
