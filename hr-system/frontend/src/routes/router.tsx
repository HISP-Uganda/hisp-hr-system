import React from "react";
import { createRootRouteWithContext, createRoute, createRouter, Outlet, redirect } from "@tanstack/react-router";
import { AuthContextValue } from "../auth/AuthContext";
import { LoginPage } from "../components/LoginPage";
import { AppShell } from "../components/AppShell";
import { SectionPage } from "../components/SectionPage";
import { UserRole, hasRole } from "./access";

type RouterContext = {
  auth: AuthContextValue;
};

function requireAuth(context: RouterContext) {
  if (!context.auth.isAuthenticated) {
    throw redirect({ to: "/login" });
  }
}

function requireRole(context: RouterContext, allowedRoles: UserRole[]) {
  requireAuth(context);
  if (!hasRole(context.auth, allowedRoles)) {
    throw redirect({ to: "/dashboard" });
  }
}

const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: () => <Outlet />,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  beforeLoad: ({ context }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({ to: "/dashboard" });
    }
  },
  component: LoginPage,
});

const shellRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  beforeLoad: ({ context, location }) => {
    requireAuth(context);
    if (location.pathname === "/") {
      throw redirect({ to: "/dashboard" });
    }
  },
  component: AppShell,
});

const dashboardRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "dashboard",
  component: () => <SectionPage title="Dashboard" description="Overview and high-level HR metrics will be displayed here." />,
});

const employeesRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "employees",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "HR Officer"]),
  component: () => <SectionPage title="Employees" description="Employee management UI will be implemented in Phase 5." />,
});

const departmentsRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "departments",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "HR Officer"]),
  component: () => <SectionPage title="Departments" description="Department management UI will be implemented in Phase 6." />,
});

const leaveRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "leave",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "HR Officer"]),
  component: () => <SectionPage title="Leave" description="Leave workflow pages will be implemented in Phase 7." />,
});

const payrollRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "payroll",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "Finance Officer"]),
  component: () => <SectionPage title="Payroll" description="Payroll workflow pages will be implemented in Phase 8." />,
});

const usersRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "users",
  beforeLoad: ({ context }) => requireRole(context, ["Admin"]),
  component: () => <SectionPage title="Users" description="Admin user management UI will be implemented in Phase 9." />,
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  shellRoute.addChildren([
    dashboardRoute,
    employeesRoute,
    departmentsRoute,
    leaveRoute,
    payrollRoute,
    usersRoute,
  ]),
]);

export const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  context: {
    auth: undefined as unknown as AuthContextValue,
  },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
