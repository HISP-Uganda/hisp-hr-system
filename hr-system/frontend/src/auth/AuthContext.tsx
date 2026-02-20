import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { Login, Logout, Me, Refresh } from "../../wailsjs/go/main/App";
import { AuthResult, AuthUser, AuthApiResponse, MeApiResponse } from "./types";
import { tokenStorage } from "./storage";

export type AuthContextValue = {
  isLoading: boolean;
  isAuthenticated: boolean;
  user: AuthUser | null;
  accessToken: string | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshSession: () => Promise<boolean>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

function normalizeError(err: unknown): string {
  if (typeof err === "string") {
    return err;
  }
  if (err instanceof Error && err.message) {
    return err.message;
  }
  return "authentication error";
}

function applyAuthResult(result: AuthResult, setUser: (user: AuthUser | null) => void, setAccessToken: (token: string | null) => void): void {
  tokenStorage.setAccessToken(result.access_token);
  tokenStorage.setRefreshToken(result.refresh_token);
  setUser(result.user);
  setAccessToken(result.access_token);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState<AuthUser | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(tokenStorage.getAccessToken());

  const clearAuth = () => {
    tokenStorage.clearAll();
    setUser(null);
    setAccessToken(null);
  };

  const refreshSession = async (): Promise<boolean> => {
    const refreshToken = tokenStorage.getRefreshToken();
    if (!refreshToken) {
      clearAuth();
      return false;
    }

    try {
      const response = (await Refresh(refreshToken)) as AuthApiResponse;
      applyAuthResult(response.data, setUser, setAccessToken);
      return true;
    } catch {
      clearAuth();
      return false;
    }
  };

  const loadMe = async (token: string): Promise<boolean> => {
    try {
      const meResponse = (await Me(token)) as MeApiResponse;
      setUser(meResponse.data);
      return true;
    } catch {
      return false;
    }
  };

  useEffect(() => {
    const init = async () => {
      const storedAccess = tokenStorage.getAccessToken();
      if (!storedAccess) {
        setIsLoading(false);
        return;
      }

      const valid = await loadMe(storedAccess);
      if (valid) {
        setAccessToken(storedAccess);
        setIsLoading(false);
        return;
      }

      await refreshSession();
      setIsLoading(false);
    };

    void init();
  }, []);

  const login = async (username: string, password: string): Promise<void> => {
    try {
      const response = (await Login(username, password)) as AuthApiResponse;
      applyAuthResult(response.data, setUser, setAccessToken);
    } catch (err) {
      throw new Error(normalizeError(err));
    }
  };

  const logout = async (): Promise<void> => {
    const refreshToken = tokenStorage.getRefreshToken();
    try {
      if (refreshToken) {
        await Logout(refreshToken);
      }
    } finally {
      clearAuth();
    }
  };

  const value = useMemo<AuthContextValue>(
    () => ({
      isLoading,
      isAuthenticated: Boolean(user && accessToken),
      user,
      accessToken,
      login,
      logout,
      refreshSession,
    }),
    [isLoading, user, accessToken],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return context;
}
