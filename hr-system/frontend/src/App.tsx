import React from "react";
import { CssBaseline, ThemeProvider, createTheme } from "@mui/material";
import { RouterProvider } from "@tanstack/react-router";
import { AuthProvider, useAuth } from "./auth/AuthContext";
import { router } from "./routes/router";

const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#1e4e8c",
    },
    secondary: {
      main: "#0f766e",
    },
    background: {
      default: "#f4f7fb",
    },
  },
  shape: {
    borderRadius: 10,
  },
  typography: {
    fontFamily: "Nunito, sans-serif",
  },
});

function AppRouter() {
  const auth = useAuth();

  if (auth.isLoading) {
    return null;
  }

  return <RouterProvider router={router} context={{ auth }} />;
}

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <AppRouter />
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
