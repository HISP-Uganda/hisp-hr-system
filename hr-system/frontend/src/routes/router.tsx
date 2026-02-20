import React from "react";
import { createRootRouteWithContext, createRoute, createRouter, Outlet, redirect } from "@tanstack/react-router";
import { AuthContextValue } from "../auth/AuthContext";
import { LoginPage } from "../components/LoginPage";
import { AppShell } from "../components/AppShell";
import { SectionPage } from "../components/SectionPage";
import { EmployeesPage } from "../modules/employees/EmployeesPage";
import { LeavePage } from "../modules/leave/LeavePage";
import { PayrollBatchesPage } from "../modules/payroll/PayrollBatchesPage";
import { PayrollBatchDetailPage } from "../modules/payroll/PayrollBatchDetailPage";
import { UsersPage } from "../modules/users/UsersPage";
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
    throw redirect({ to: "/access-denied" });
  }
}

const rootRoute = createRootRouteWithContext<RouterContext>()({
  notFoundComponent: () => <SectionPage title="Not Found" description="The requested page was not found." />,
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
  path: "/dashboard",
  component: AppShell,
  // component: () => <SectionPage title="Dashboard" description="Overview and high-level HR metrics will be displayed here." />,
});

const dashboardTypoRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/dashboad",
  beforeLoad: () => {
    throw redirect({ to: "/dashboard" });
  },
  component: () => null,
});

const accessDeniedRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/access-denied",
  component: () => <SectionPage title="Access Denied" description="You do not have permission to view this page." />,
});

const employeesRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/employees",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "HR Officer"]),
  component: EmployeesPage,
});

const departmentsRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/departments",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "HR Officer"]),
  component: () => <SectionPage title="Departments" description="Department management UI will be implemented in Phase 6." />,
});

const leaveRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/leave",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "HR Officer", "Finance Officer", "Viewer"]),
  component: LeavePage,
});

const payrollRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/payroll",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "Finance Officer"]),
  component: PayrollBatchesPage,
});

const payrollDetailRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/payroll/$batchId",
  beforeLoad: ({ context }) => requireRole(context, ["Admin", "Finance Officer"]),
  component: PayrollBatchDetailPage,
});

const usersRoute = createRoute({
  getParentRoute: () => shellRoute,
  path: "/users",
  beforeLoad: ({ context }) => requireRole(context, ["Admin"]),
  component: UsersPage,
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  shellRoute.addChildren([
    dashboardRoute,
    dashboardTypoRoute,
    accessDeniedRoute,
    employeesRoute,
    departmentsRoute,
    leaveRoute,
    payrollRoute,
    payrollDetailRoute,
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
