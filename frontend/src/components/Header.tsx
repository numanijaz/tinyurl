import { Link } from "react-router-dom";
import {Button} from "@/components/ui/button";
import { useAuth } from "@/auth";

export function Header() {
  const {userInfo, logout} = useAuth();
  return (
    <header className="w-full border-b py-3 px-6 flex items-center justify-between">
      <div className="flex items-center gap-4">
        <Link to="/" className="font-bold">TinyURL</Link>
        <nav className="flex gap-3">
          <Link to="/">Home</Link>
        </nav>
      </div>
      <div>
        {userInfo ? (
          <div className="flex items-center gap-3">
            Hi {userInfo.name}
            <Button size={"sm"} variant={"link"} onClick={logout}>Logout</Button>
          </div>
        ) : (
          <div className="flex">
            <Link to="/login">
              <Button size={"sm"} variant={"link"}>Login</Button>
            </Link>
            <Link to="/signup">
              <Button size={"sm"} variant={"link"}>Signup</Button>
            </Link>
          </div>
        )}
      </div>
    </header>
  );
}
