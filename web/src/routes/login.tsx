import { useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { FcGoogle } from "react-icons/fc";
import { getAccessToken, LOGIN_URL } from "@/lib/auth";
import { KronMark } from "@/components/kron-mark";

export function LoginPage() {
  const navigate = useNavigate();

  useEffect(() => {
    if (getAccessToken()) navigate({ to: "/dashboard", replace: true });
  }, [navigate]);

  return (
    <div className="relative min-h-screen overflow-hidden bg-background text-foreground">
      {/* Concentric orbit lines, decorative */}
      <div aria-hidden className="pointer-events-none absolute inset-0 grid place-items-center">
        <div className="relative h-[820px] w-[820px] opacity-[0.08]">
          <div className="absolute inset-0 rounded-full border border-foreground" />
          <div className="absolute inset-12 rounded-full border border-foreground" />
          <div className="absolute inset-28 rounded-full border border-foreground" />
          <div className="absolute inset-48 rounded-full border border-dashed border-foreground animate-orbit" />
        </div>
      </div>

      <main className="relative z-10 mx-auto flex min-h-screen w-full max-w-md flex-col items-center justify-center px-6">
        <div className="flex flex-col items-center text-center animate-fade-in">
          <KronMark className="h-10 w-10 text-foreground" />
          <span className="mt-3 text-[15px] font-semibold tracking-tight">Kron</span>

          <h1 className="font-display mt-10 text-[44px] leading-[1.05]">
            Welcome <span className="italic">back</span>
          </h1>
          <p className="mt-3 text-sm text-muted-foreground">
            Sign in to your account to continue
          </p>
        </div>

        <div className="mt-10 flex w-full flex-col gap-3">
          <button
            type="button"
            onClick={() => {
              window.location.href = LOGIN_URL;
            }}
            className="group flex h-12 w-full items-center justify-center gap-3 rounded-full border border-border bg-card px-4 text-sm font-medium text-foreground transition-all hover:bg-muted active:scale-[0.99]"
          >
            <FcGoogle className="h-[18px] w-[18px]" />
            Continue with Google
          </button>
        </div>

        <p className="mt-10 max-w-[280px] text-center text-[11px] leading-relaxed text-muted-foreground">
          By continuing, you agree to Kron&rsquo;s{" "}
          <a href="#" className="text-foreground underline decoration-border underline-offset-2 hover:decoration-foreground">
            Terms of Service
          </a>{" "}
          and{" "}
          <a href="#" className="text-foreground underline decoration-border underline-offset-2 hover:decoration-foreground">
            Privacy Policy
          </a>
          .
        </p>
      </main>

      {/* Bottom hairline */}
      <div className="pointer-events-none absolute bottom-0 left-0 right-0 h-px bg-border" />
    </div>
  );
}
