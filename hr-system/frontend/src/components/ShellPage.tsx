import React from "react";
import { Button, Chip, Paper, Stack, Typography } from "@mui/material";
import { useNavigate } from "@tanstack/react-router";
import { useAuth } from "../auth/AuthContext";

export function ShellPage() {
  const auth = useAuth();
  const navigate = useNavigate();

  const onLogout = async () => {
    await auth.logout();
    await navigate({ to: "/login" });
  };

  return (
    <Stack sx={{ minHeight: "100vh", p: 4, background: "radial-gradient(circle at top right, #e8f0ff 0%, #f8fbff 40%, #eef4f1 100%)" }}>
      <Paper elevation={1} sx={{ p: 2.5, display: "flex", justifyContent: "space-between", alignItems: "center", borderRadius: 2.5 }}>
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Typography variant="h6" sx={{ fontWeight: 700 }}>HISP HR System</Typography>
          <Chip size="small" color="primary" label={auth.user?.role ?? "Unknown"} />
        </Stack>
        <Stack direction="row" spacing={2} alignItems="center">
          <Typography variant="body2">{auth.user?.username}</Typography>
          <Button variant="outlined" onClick={onLogout}>Logout</Button>
        </Stack>
      </Paper>

      <Paper elevation={0} sx={{ mt: 3, p: 4, borderRadius: 2.5, border: "1px solid #d7dfeb", backgroundColor: "rgba(255,255,255,0.9)" }}>
        <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>Welcome</Typography>
        <Typography variant="body1" color="text.secondary">
          Authenticated shell is active. Module navigation will be added in the next phase.
        </Typography>
      </Paper>
    </Stack>
  );
}
