import { useEffect, useState } from "react";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { usePostHog } from "posthog-js/react";
import { setTokens } from "@/lib/auth";

export function CallbackPage() {
  const navigate = useNavigate();
  const posthog = usePostHog();
  const search = useSearch({ strict: false }) as {
    access?: string;
    refresh?: string;
  };
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!search.access || !search.refresh) {
      setError("Missing tokens in callback");
      return;
    }
    setTokens({
      access_token: search.access,
      refresh_token: search.refresh,
    });
    posthog.capture('login_completed', { method: 'google' });
    navigate({ to: "/dashboard", replace: true });
  }, [search.access, search.refresh, navigate, posthog]);

  return (
    <div className="flex h-screen items-center justify-center">
      {error ? (
        <p className="text-red-500">{error}</p>
      ) : (
        <p className="text-muted-foreground">Signing you in…</p>
      )}
    </div>
  );
}
