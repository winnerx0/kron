import { useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Moon, Sun } from "lucide-react";
import { FcGoogle } from "react-icons/fc";
import { Button } from "@/components/ui/button";
import { getAccessToken, LOGIN_URL } from "@/lib/auth";
import { useTheme } from "@/lib/theme";

export function LoginPage() {
  const navigate = useNavigate();
  const { theme, toggle } = useTheme();

  useEffect(() => {
    if (getAccessToken()) navigate({ to: "/", replace: true });
  }, [navigate]);

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="h-14 border-b border-border bg-card">
        <div className="mx-auto flex h-full w-full max-w-5xl items-center justify-between px-6">
          <div className="flex items-center gap-2.5">
            <LogoMark />
            <span className="text-[15px] font-semibold tracking-tight">
              Kron
            </span>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="gap-2"
            onClick={toggle}
          >
            {theme === "dark" ? (
              <Sun className="h-3.5 w-3.5" />
            ) : (
              <Moon className="h-3.5 w-3.5" />
            )}
            {theme === "dark" ? "Light" : "Dark"}
          </Button>
        </div>
      </header>

      <main className="mx-auto flex flex-col min-h-[calc(100vh-3.5rem)] w-full max-w-5xl items-center justify-center px-6 py-10">
          <div className="px-6 py-4 text-center">
            <h2 className="text-2xl font-semibold tracking-tight">
              Welcome back
            </h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Sign in to your account to continue
            </p>
          </div>

          <div className="p-6">
            <Button
              size="lg"
              variant="outline"
              className="h-11 w-full justify-start gap-3 px-4 rounded-full"
              onClick={() => {
                window.location.href = LOGIN_URL;
              }}
            >
              <span className="grid h-7 w-7 place-items-center rounded-md border border-border bg-background">
                <FcGoogle className="h-4 w-4" />
              </span>
              Continue with Google
            </Button>
          </div>
      </main>
    </div>
  );
}

function LogoMark() {
  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 20 20"
      fill="none"
      className="shrink-0"
    >
      <circle cx="10" cy="10" r="8" stroke="currentColor" strokeWidth="1.5" />
      <line
        x1="10"
        y1="10"
        x2="10"
        y2="4.5"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <line
        x1="10"
        y1="10"
        x2="13.5"
        y2="12"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <circle cx="10" cy="10" r="1" fill="currentColor" />
    </svg>
  );
}
