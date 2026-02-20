import React, { useMemo, useState } from "react";
import {
  AppBar,
  Box,
  Button,
  Chip,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemText,
  Stack,
  Toolbar,
  Typography,
} from "@mui/material";
import { Outlet, useNavigate, useRouterState } from "@tanstack/react-router";
import { useAuth } from "../auth/AuthContext";
import { UserRole } from "../routes/access";

const drawerWidth = 240;

type NavItem = {
  label: string;
  to: "/dashboard" | "/employees" | "/departments" | "/leave" | "/payroll" | "/users";
  roles: UserRole[];
};

const navItems: NavItem[] = [
  { label: "Dashboard", to: "/dashboard", roles: ["Admin", "HR Officer", "Finance Officer", "Viewer"] },
  { label: "Employees", to: "/employees", roles: ["Admin", "HR Officer"] },
  { label: "Departments", to: "/departments", roles: ["Admin", "HR Officer"] },
  { label: "Leave", to: "/leave", roles: ["Admin", "HR Officer", "Finance Officer", "Viewer"] },
  { label: "Payroll", to: "/payroll", roles: ["Admin", "Finance Officer"] },
  { label: "Users", to: "/users", roles: ["Admin"] },
];

export function AppShell() {
  const auth = useAuth();
  const navigate = useNavigate();
  const pathname = useRouterState({ select: (state) => state.location.pathname });
  const [mobileOpen, setMobileOpen] = useState(false);

  const allowedItems = useMemo(() => {
    const role = auth.user?.role as UserRole | undefined;
    if (!role) {
      return [];
    }
    return navItems.filter((item) => item.roles.includes(role));
  }, [auth.user?.role]);

  const onNavigate = async (to: NavItem["to"]) => {
    setMobileOpen(false);
    await navigate({ to });
  };

  const onLogout = async () => {
    await auth.logout();
    await navigate({ to: "/login" });
  };

  const drawerContent = (
    <Box sx={{ height: "100%", background: "linear-gradient(180deg, #113f78 0%, #0a2c53 100%)", color: "#f3f8ff" }}>
      <Toolbar>
        <Stack spacing={0.1}>
          <Typography variant="h6" sx={{ fontWeight: 800, lineHeight: 1.1 }}>HISP HR</Typography>
          <Typography variant="caption" sx={{ opacity: 0.85 }}>Management System</Typography>
        </Stack>
      </Toolbar>
      <Divider sx={{ borderColor: "rgba(255,255,255,0.2)" }} />
      <List sx={{ pt: 1.2 }}>
        {allowedItems.map((item) => (
          <ListItemButton
            key={item.to}
            selected={pathname === item.to}
            onClick={() => void onNavigate(item.to)}
            sx={{
              mx: 1,
              mb: 0.5,
              borderRadius: 1.5,
              color: "inherit",
              "&.Mui-selected": { backgroundColor: "rgba(255,255,255,0.2)" },
              "&.Mui-selected:hover": { backgroundColor: "rgba(255,255,255,0.25)" },
            }}
          >
            <ListItemText primary={item.label} />
          </ListItemButton>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: "flex", minHeight: "100vh", background: "radial-gradient(circle at top right, #e9f1ff 0%, #f7fbff 55%, #eef6f4 100%)" }}>
      <AppBar
        position="fixed"
        color="transparent"
        elevation={0}
        sx={{
          backdropFilter: "blur(8px)",
          borderBottom: "1px solid #d5deec",
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
        }}
      >
        <Toolbar sx={{ justifyContent: "space-between" }}>
          <Stack direction="row" spacing={1.4} alignItems="center">
            <IconButton color="primary" onClick={() => setMobileOpen((prev) => !prev)} sx={{ display: { md: "none" } }}>
              <Typography sx={{ fontWeight: 900, fontSize: 18 }}>≡</Typography>
            </IconButton>
            <Typography variant="h6" sx={{ fontWeight: 700 }}>HISP HR System</Typography>
          </Stack>

          <Stack direction="row" spacing={1.2} alignItems="center">
            <Typography variant="body2">{auth.user?.username}</Typography>
            <Chip size="small" color="primary" label={auth.user?.role ?? "Unknown"} />
            <Button variant="outlined" onClick={onLogout}>Logout</Button>
          </Stack>
        </Toolbar>
      </AppBar>

      <Box component="nav" sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}>
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={() => setMobileOpen(false)}
          ModalProps={{ keepMounted: true }}
          sx={{ display: { xs: "block", md: "none" }, "& .MuiDrawer-paper": { width: drawerWidth, boxSizing: "border-box" } }}
        >
          {drawerContent}
        </Drawer>
        <Drawer
          variant="permanent"
          open
          sx={{ display: { xs: "none", md: "block" }, "& .MuiDrawer-paper": { width: drawerWidth, boxSizing: "border-box", borderRight: "none" } }}
        >
          {drawerContent}
        </Drawer>
      </Box>

      <Box component="main" sx={{ flexGrow: 1, p: { xs: 2, md: 3 }, width: { md: `calc(100% - ${drawerWidth}px)` } }}>
        <Toolbar />
        <Outlet />
      </Box>
    </Box>
  );
}
