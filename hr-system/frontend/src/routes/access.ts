import { AuthContextValue } from "../auth/AuthContext";

export type UserRole = "Admin" | "HR Officer" | "Finance Officer" | "Viewer";

export function hasRole(auth: AuthContextValue, roles: UserRole[]): boolean {
  const role = auth.user?.role;
  if (!role) {
    return false;
  }
  return roles.includes(role as UserRole);
}
