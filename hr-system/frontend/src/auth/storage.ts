const ACCESS_TOKEN_KEY = "hr.auth.access_token";
const REFRESH_TOKEN_KEY = "hr.auth.refresh_token";

export const tokenStorage = {
  getAccessToken(): string | null {
    return window.localStorage.getItem(ACCESS_TOKEN_KEY);
  },
  setAccessToken(token: string): void {
    window.localStorage.setItem(ACCESS_TOKEN_KEY, token);
  },
  clearAccessToken(): void {
    window.localStorage.removeItem(ACCESS_TOKEN_KEY);
  },
  getRefreshToken(): string | null {
    return window.localStorage.getItem(REFRESH_TOKEN_KEY);
  },
  setRefreshToken(token: string): void {
    window.localStorage.setItem(REFRESH_TOKEN_KEY, token);
  },
  clearRefreshToken(): void {
    window.localStorage.removeItem(REFRESH_TOKEN_KEY);
  },
  clearAll(): void {
    this.clearAccessToken();
    this.clearRefreshToken();
  },
};
