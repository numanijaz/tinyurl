import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";


export function NotFound() {
  const navigate = useNavigate();
  return (
    <div className="p-6 flex flex-col items-center gap-6">
      <h1 className="text-4xl font-semibold tracking-tight text-balance">
        404
      </h1>
      <p>Page not found.</p>
      <Button onClick={() => navigate("/")}>
        Return to website
      </Button>
    </div>
  );
}
