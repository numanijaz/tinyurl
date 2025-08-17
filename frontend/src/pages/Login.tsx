import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Link, useNavigate } from "react-router-dom"
import { Icons } from "@/icons"
import { FormEvent, useCallback, useState } from "react"
import axios from "axios"
import { useAuth } from "@/auth"
import { usePageTitle } from "@/hooks"

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  usePageTitle("Login");
  const [loading, setLoading] = useState<boolean>();
  const [error, setError] = useState<string>();
  const {loadUserInfo} = useAuth();

  const navigate = useNavigate();

  const handleLoginSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    setLoading(true);
    axios.postForm("/api/auth/login", formData)
      .then(() => {
        setLoading(false);
        navigate("/"); // frontend-reroute
        loadUserInfo?.();
      })
      .catch((error) => {
        setLoading(false);
        setError(error?.response?.data?.error);
      });
  }, [navigate]);

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle>Login to your account</CardTitle>
          <CardDescription>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4">
            <div className="text-muted-foreground text-sm">
              Enter your email below to login to your account
            </div>
            <form onSubmit={handleLoginSubmit}>
              <div className="flex flex-col gap-4">
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    name="email"
                    placeholder="m@example.com"
                    autoComplete="on"
                    required
                  />
                </div>
                <div className="grid gap-3">
                  <div className="flex items-center">
                    <Label htmlFor="password">Password</Label>
                    <a
                      href="#"
                      className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                    >
                      Forgot your password?
                    </a>
                  </div>
                  <Input id="password" type="password" name="password" required />
                </div>
                {error ? <p className="text-sm">{error}</p> : null}
                <div className="flex flex-col gap-3">
                  <Button type="submit" className="w-full">
                    Login
                  </Button>
                </div>
              </div>
            </form>
          </div>
          <div className="flex flex-col gap-4">
            <div className="text-muted-foreground text-sm pt-4">
              Or Login using
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
          <div className="mt-4 text-center text-sm">
            Don&apos;t have an account?{" "}
            <Link to={"/signup"} className="underline underline-offset-4">Sign up</Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
