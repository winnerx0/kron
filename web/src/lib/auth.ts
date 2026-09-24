const API_URL = import.meta.env.VITE_API_URL || "http://localhost:5000/api";
const ROOT_URL = API_URL.replace(/\/api\/?$/, "");
export const LOGIN_URL = `${ROOT_URL}/oauth/google/login`;

const ACCESS_TOKEN_KEY = "kron.access_token";
const REFRESH_TOKEN_KEY = "kron.refresh_token";

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setTokens(tokens: TokenPair): void {
  localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token);
  localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token);
}

export function clearTokens(): void {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export function redirectToLogin(): void {
  clearTokens();
  if (window.location.pathname !== "/login") {
    window.location.href = "/login";
  }
}

let refreshInFlight: Promise<TokenPair | null> | null = null;

export async function refreshTokens(): Promise<TokenPair | null> {
  if (refreshInFlight) return refreshInFlight;

  const refresh_token = getRefreshToken();
  if (!refresh_token) return null;

  refreshInFlight = (async () => {
    try {
      const res = await fetch(`${API_URL}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token }),
      });
      if (!res.ok) {
        clearTokens();
        return null;
      }
      const tokens = (await res.json()) as TokenPair;
      setTokens(tokens);
      return tokens;
    } finally {
      refreshInFlight = null;
    }
  })();

  return refreshInFlight;
}
