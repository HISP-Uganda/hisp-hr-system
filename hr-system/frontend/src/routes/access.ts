import { AuthContextValue } from "../auth/AuthContext";

export type UserRole = "Admin" | "HR Officer" | "Finance Officer" | "Viewer";

export function hasRole(auth: AuthContextValue, roles: UserRole[]): boolean {
  const role = auth.user?.role;
  if (!role) {
    return false;
  }
  const normalized = role.trim().toLowerCase();
  return roles.some((allowed) => allowed.toLowerCase() === normalized);
}
