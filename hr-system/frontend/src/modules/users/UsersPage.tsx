import React, { useCallback, useEffect, useMemo, useState } from "react";
import {
  Alert,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Paper,
  Snackbar,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import {
  CreateUser,
  GetUser,
  ListUsers,
  ResetUserPassword,
  SetUserStatus,
  UpdateUser,
} from "../../../wailsjs/go/main/App";
import { useAuth } from "../../auth/AuthContext";
import {
  CreateUserInput,
  ResetPasswordInput,
  UpdateUserInput,
  UserListQuery,
  UserListResponse,
  UserResponse,
  UserView,
} from "./types";

type SnackbarState = {
  open: boolean;
  severity: "success" | "error";
  message: string;
};

type CreateForm = {
  username: string;
  password: string;
  role: string;
};

type EditForm = {
  username: string;
  role: string;
};

type ResetForm = {
  new_password: string;
  confirm_password: string;
};

const defaultCreateForm: CreateForm = { username: "", password: "", role: "admin" };
const defaultEditForm: EditForm = { username: "", role: "admin" };
const defaultResetForm: ResetForm = { new_password: "", confirm_password: "" };

function normalizeError(err: unknown): string {
  if (typeof err === "string") return err;
  if (err instanceof Error && err.message) return err.message;
  return "request failed";
}

function formatDate(value?: string | null): string {
  if (!value) return "-";
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return value;
  return parsed.toLocaleString();
}

export function UsersPage() {
  const auth = useAuth();
  const accessToken = auth.accessToken;

  const [items, setItems] = useState<UserView[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [searchInput, setSearchInput] = useState("");
  const [appliedSearch, setAppliedSearch] = useState("");
  const [loading, setLoading] = useState(false);

  const [createOpen, setCreateOpen] = useState(false);
  const [createForm, setCreateForm] = useState<CreateForm>(defaultCreateForm);
  const [createSaving, setCreateSaving] = useState(false);

  const [editOpen, setEditOpen] = useState(false);
  const [editTargetID, setEditTargetID] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<EditForm>(defaultEditForm);
  const [editSaving, setEditSaving] = useState(false);

  const [resetOpen, setResetOpen] = useState(false);
  const [resetTargetID, setResetTargetID] = useState<number | null>(null);
  const [resetForm, setResetForm] = useState<ResetForm>(defaultResetForm);
  const [resetSaving, setResetSaving] = useState(false);

  const [statusTarget, setStatusTarget] = useState<UserView | null>(null);
  const [statusSaving, setStatusSaving] = useState(false);

  const [snackbar, setSnackbar] = useState<SnackbarState>({
    open: false,
    severity: "success",
    message: "",
  });

  const showError = useCallback((message: string) => {
    setSnackbar({ open: true, severity: "error", message });
  }, []);

  const showSuccess = useCallback((message: string) => {
    setSnackbar({ open: true, severity: "success", message });
  }, []);

  const query = useMemo<UserListQuery>(
    () => ({ page: page + 1, page_size: pageSize, q: appliedSearch }),
    [page, pageSize, appliedSearch],
  );

  const loadUsers = useCallback(async () => {
    if (!accessToken) return;

    setLoading(true);
    try {
      const response = (await ListUsers(accessToken, query)) as UserListResponse;
      setItems(response.data.items);
      setTotal(response.data.total);
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setLoading(false);
    }
  }, [accessToken, query, showError]);

  useEffect(() => {
    void loadUsers();
  }, [loadUsers]);

  const onCreate = async () => {
    if (!accessToken) return;
    if (!createForm.username.trim()) {
      showError("Username is required");
      return;
    }
    if (createForm.password.trim().length < 8) {
      showError("Password must be at least 8 characters");
      return;
    }
    if (!createForm.role.trim()) {
      showError("Role is required");
      return;
    }

    setCreateSaving(true);
    try {
      const payload: CreateUserInput = {
        username: createForm.username.trim(),
        password: createForm.password,
        role: createForm.role.trim(),
      };
      await CreateUser(accessToken, payload);
      setCreateOpen(false);
      setCreateForm(defaultCreateForm);
      showSuccess("User created");
      await loadUsers();
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setCreateSaving(false);
    }
  };

  const openEdit = async (id: number) => {
    if (!accessToken) return;
    try {
      const response = (await GetUser(accessToken, id)) as UserResponse;
      setEditTargetID(id);
      setEditForm({ username: response.data.username, role: response.data.role });
      setEditOpen(true);
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onSaveEdit = async () => {
    if (!accessToken || editTargetID == null) return;
    if (!editForm.username.trim() || !editForm.role.trim()) {
      showError("Username and role are required");
      return;
    }

    setEditSaving(true);
    try {
      const payload: UpdateUserInput = { username: editForm.username.trim(), role: editForm.role.trim() };
      await UpdateUser(accessToken, editTargetID, payload);
      setEditOpen(false);
      setEditTargetID(null);
      showSuccess("User updated");
      await loadUsers();
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setEditSaving(false);
    }
  };

  const onResetPassword = async () => {
    if (!accessToken || resetTargetID == null) return;
    if (resetForm.new_password.trim().length < 8) {
      showError("Password must be at least 8 characters");
      return;
    }
    if (resetForm.new_password !== resetForm.confirm_password) {
      showError("Passwords do not match");
      return;
    }

    setResetSaving(true);
    try {
      const payload: ResetPasswordInput = { new_password: resetForm.new_password };
      await ResetUserPassword(accessToken, resetTargetID, payload);
      setResetOpen(false);
      setResetTargetID(null);
      setResetForm(defaultResetForm);
      showSuccess("Password reset successful");
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setResetSaving(false);
    }
  };

  const onConfirmStatusChange = async () => {
    if (!accessToken || !statusTarget) return;
    setStatusSaving(true);
    try {
      await SetUserStatus(accessToken, statusTarget.id, { is_active: !statusTarget.is_active });
      setStatusTarget(null);
      showSuccess("User status updated");
      await loadUsers();
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setStatusSaving(false);
    }
  };

  return (
    <Paper elevation={0} sx={{ p: { xs: 2, md: 3 }, borderRadius: 2.5, border: "1px solid #dbe3ef", backgroundColor: "rgba(255,255,255,0.92)" }}>
      <Stack spacing={2}>
        <Stack direction={{ xs: "column", md: "row" }} spacing={1.2} justifyContent="space-between" alignItems={{ md: "center" }}>
          <Box>
            <Typography variant="h5" sx={{ fontWeight: 800 }}>Users</Typography>
            <Typography variant="body2" color="text.secondary">Admin user management: create, update role, reset password, and activate/deactivate.</Typography>
          </Box>
          <Stack direction={{ xs: "column", sm: "row" }} spacing={1}>
            <TextField
              size="small"
              label="Search username"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
            <Button variant="outlined" onClick={() => { setPage(0); setAppliedSearch(searchInput.trim()); }}>Search</Button>
            <Button variant="contained" onClick={() => setCreateOpen(true)}>Create User</Button>
          </Stack>
        </Stack>

        <Box sx={{ overflowX: "auto" }}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Username</TableCell>
                <TableCell>Role</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Last Login</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {items.map((item) => {
                const isSelf = auth.user?.id === item.id;
                return (
                  <TableRow key={item.id} hover>
                    <TableCell>{item.username}</TableCell>
                    <TableCell>{item.role}</TableCell>
                    <TableCell>
                      <Chip size="small" color={item.is_active ? "success" : "default"} label={item.is_active ? "Active" : "Inactive"} />
                    </TableCell>
                    <TableCell>{formatDate(item.created_at)}</TableCell>
                    <TableCell>{formatDate(item.last_login_at)}</TableCell>
                    <TableCell align="right">
                      <Stack direction="row" spacing={1} justifyContent="flex-end">
                        <Button size="small" onClick={() => void openEdit(item.id)}>Edit</Button>
                        <Button size="small" onClick={() => { setResetTargetID(item.id); setResetOpen(true); }}>Reset Password</Button>
                        <Button
                          size="small"
                          color={item.is_active ? "warning" : "success"}
                          disabled={isSelf && item.is_active}
                          onClick={() => setStatusTarget(item)}
                        >
                          {item.is_active ? "Deactivate" : "Activate"}
                        </Button>
                      </Stack>
                    </TableCell>
                  </TableRow>
                );
              })}
              {!loading && items.length === 0 && (
                <TableRow>
                  <TableCell colSpan={6}>
                    <Typography variant="body2" color="text.secondary">No users found.</Typography>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </Box>

        <TablePagination
          component="div"
          count={total}
          page={page}
          rowsPerPage={pageSize}
          onPageChange={(_, next) => setPage(next)}
          onRowsPerPageChange={(event) => {
            const nextSize = Number(event.target.value);
            setPageSize(nextSize);
            setPage(0);
          }}
          rowsPerPageOptions={[5, 10, 20, 50]}
        />
      </Stack>

      <Dialog open={createOpen} onClose={() => setCreateOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Create User</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 0.5 }}>
            <TextField label="Username" value={createForm.username} onChange={(e) => setCreateForm((prev) => ({ ...prev, username: e.target.value }))} fullWidth />
            <TextField label="Password" type="password" value={createForm.password} onChange={(e) => setCreateForm((prev) => ({ ...prev, password: e.target.value }))} fullWidth />
            <TextField label="Role" value={createForm.role} onChange={(e) => setCreateForm((prev) => ({ ...prev, role: e.target.value }))} fullWidth />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => void onCreate()} disabled={createSaving}>Create</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={editOpen} onClose={() => setEditOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Edit User</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 0.5 }}>
            <TextField label="Username" value={editForm.username} onChange={(e) => setEditForm((prev) => ({ ...prev, username: e.target.value }))} fullWidth />
            <TextField label="Role" value={editForm.role} onChange={(e) => setEditForm((prev) => ({ ...prev, role: e.target.value }))} fullWidth />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => void onSaveEdit()} disabled={editSaving}>Save</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={resetOpen} onClose={() => setResetOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Reset Password</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 0.5 }}>
            <TextField label="New Password" type="password" value={resetForm.new_password} onChange={(e) => setResetForm((prev) => ({ ...prev, new_password: e.target.value }))} fullWidth />
            <TextField label="Confirm Password" type="password" value={resetForm.confirm_password} onChange={(e) => setResetForm((prev) => ({ ...prev, confirm_password: e.target.value }))} fullWidth />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setResetOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => void onResetPassword()} disabled={resetSaving}>Reset</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={Boolean(statusTarget)} onClose={() => setStatusTarget(null)} fullWidth maxWidth="xs">
        <DialogTitle>Confirm Status Change</DialogTitle>
        <DialogContent>
          <Typography variant="body2">
            {statusTarget?.is_active
              ? `Deactivate ${statusTarget?.username}?`
              : `Activate ${statusTarget?.username}?`}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setStatusTarget(null)}>Cancel</Button>
          <Button variant="contained" onClick={() => void onConfirmStatusChange()} disabled={statusSaving}>Confirm</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={3500}
        onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      >
        <Alert
          severity={snackbar.severity}
          onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}
          sx={{ width: "100%" }}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Paper>
  );
}
