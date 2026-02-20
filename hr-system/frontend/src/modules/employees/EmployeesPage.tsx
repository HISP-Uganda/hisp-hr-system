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
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Snackbar,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TablePagination,
  TextField,
  Typography,
} from "@mui/material";
import {
  CreateEmployee,
  DeleteEmployee,
  GetEmployee,
  ListEmployeeDepartments,
  ListEmployees,
  UpdateEmployee,
} from "../../../wailsjs/go/main/App";
import { useAuth } from "../../auth/AuthContext";
import {
  DepartmentListResponse,
  DepartmentOption,
  EmployeeInput,
  EmployeeListQuery,
  EmployeeListResponse,
  EmployeeResponse,
  EmployeeView,
} from "./types";

type SnackbarState = {
  open: boolean;
  severity: "success" | "error";
  message: string;
};

type EmployeeFormState = {
  first_name: string;
  last_name: string;
  other_name: string;
  gender: string;
  dob: string;
  phone: string;
  email: string;
  national_id: string;
  address: string;
  department_id: string;
  position: string;
  employment_status: string;
  hire_date: string;
  base_salary: string;
};

const defaultForm: EmployeeFormState = {
  first_name: "",
  last_name: "",
  other_name: "",
  gender: "",
  dob: "",
  phone: "",
  email: "",
  national_id: "",
  address: "",
  department_id: "",
  position: "",
  employment_status: "",
  hire_date: "",
  base_salary: "",
};

const employmentStatuses = ["Active", "Inactive", "Suspended", "Terminated"];
const genders = ["Male", "Female", "Other"];

function normalizeError(err: unknown): string {
  if (typeof err === "string") {
    return err;
  }
  if (err instanceof Error && err.message) {
    return err.message;
  }
  return "request failed";
}

function toFormState(employee: EmployeeView): EmployeeFormState {
  return {
    first_name: employee.first_name,
    last_name: employee.last_name,
    other_name: employee.other_name,
    gender: employee.gender,
    dob: employee.dob,
    phone: employee.phone,
    email: employee.email,
    national_id: employee.national_id,
    address: employee.address,
    department_id: employee.department_id ? String(employee.department_id) : "",
    position: employee.position,
    employment_status: employee.employment_status,
    hire_date: employee.hire_date,
    base_salary: String(employee.base_salary),
  };
}

function toInput(form: EmployeeFormState): EmployeeInput {
  const departmentID = form.department_id ? Number(form.department_id) : undefined;

  return {
    first_name: form.first_name.trim(),
    last_name: form.last_name.trim(),
    other_name: form.other_name.trim(),
    gender: form.gender,
    dob: form.dob,
    phone: form.phone.trim(),
    email: form.email.trim(),
    national_id: form.national_id.trim(),
    address: form.address.trim(),
    department_id: Number.isFinite(departmentID) ? departmentID : undefined,
    position: form.position.trim(),
    employment_status: form.employment_status,
    hire_date: form.hire_date,
    base_salary: Number(form.base_salary),
  };
}

function validateForm(form: EmployeeFormState): string | null {
  if (!form.first_name.trim() || !form.last_name.trim()) return "First and last name are required";
  if (!form.gender) return "Gender is required";
  if (!form.dob) return "Date of birth is required";
  if (!form.phone.trim()) return "Phone is required";
  if (!form.position.trim()) return "Position is required";
  if (!form.employment_status) return "Employment status is required";
  if (!form.hire_date) return "Hire date is required";
  if (!form.base_salary.trim()) return "Base salary is required";

  const baseSalary = Number(form.base_salary);
  if (!Number.isFinite(baseSalary) || baseSalary < 0) return "Base salary must be a non-negative number";

  return null;
}

export function EmployeesPage() {
  const auth = useAuth();

  const [items, setItems] = useState<EmployeeView[]>([]);
  const [departments, setDepartments] = useState<DepartmentOption[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [searchInput, setSearchInput] = useState("");
  const [appliedSearch, setAppliedSearch] = useState("");
  const [departmentFilter, setDepartmentFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [loading, setLoading] = useState(false);

  const [formOpen, setFormOpen] = useState(false);
  const [formSaving, setFormSaving] = useState(false);
  const [editingEmployeeID, setEditingEmployeeID] = useState<number | null>(null);
  const [form, setForm] = useState<EmployeeFormState>(defaultForm);

  const [deleteTarget, setDeleteTarget] = useState<EmployeeView | null>(null);
  const [deleteSubmitting, setDeleteSubmitting] = useState(false);

  const [snackbar, setSnackbar] = useState<SnackbarState>({
    open: false,
    severity: "success",
    message: "",
  });

  const accessToken = auth.accessToken;

  const showError = useCallback((message: string) => {
    setSnackbar({ open: true, severity: "error", message });
  }, []);

  const showSuccess = useCallback((message: string) => {
    setSnackbar({ open: true, severity: "success", message });
  }, []);

  const listQuery = useMemo<EmployeeListQuery>(() => {
    const query: EmployeeListQuery = {
      search: appliedSearch,
      status: statusFilter,
      page: page + 1,
      page_size: pageSize,
    };
    if (departmentFilter) {
      query.department_id = Number(departmentFilter);
    }
    return query;
  }, [appliedSearch, statusFilter, page, pageSize, departmentFilter]);

  const loadEmployees = useCallback(async () => {
    if (!accessToken) {
      return;
    }

    setLoading(true);
    try {
      const response = (await ListEmployees(accessToken, listQuery)) as EmployeeListResponse;
      setItems(response.data.items);
      setTotal(response.data.total);
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setLoading(false);
    }
  }, [accessToken, listQuery, showError]);

  const loadDepartments = useCallback(async () => {
    if (!accessToken) {
      return;
    }

    try {
      const response = (await ListEmployeeDepartments(accessToken)) as DepartmentListResponse;
      setDepartments(response.data);
    } catch (err) {
      showError(normalizeError(err));
    }
  }, [accessToken, showError]);

  useEffect(() => {
    void loadDepartments();
  }, [loadDepartments]);

  useEffect(() => {
    void loadEmployees();
  }, [loadEmployees]);

  const openCreateDialog = () => {
    setEditingEmployeeID(null);
    setForm(defaultForm);
    setFormOpen(true);
  };

  const openEditDialog = async (employeeID: number) => {
    if (!accessToken) {
      return;
    }

    try {
      const response = (await GetEmployee(accessToken, employeeID)) as EmployeeResponse;
      setEditingEmployeeID(employeeID);
      setForm(toFormState(response.data));
      setFormOpen(true);
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onSubmitForm = async () => {
    if (!accessToken) {
      return;
    }

    const validationError = validateForm(form);
    if (validationError) {
      showError(validationError);
      return;
    }

    setFormSaving(true);
    try {
      const payload = toInput(form);
      if (editingEmployeeID == null) {
        await CreateEmployee(accessToken, payload);
        showSuccess("Employee created");
      } else {
        await UpdateEmployee(accessToken, editingEmployeeID, payload);
        showSuccess("Employee updated");
      }
      setFormOpen(false);
      await loadEmployees();
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setFormSaving(false);
    }
  };

  const onConfirmDelete = async () => {
    if (!accessToken || !deleteTarget) {
      return;
    }

    setDeleteSubmitting(true);
    try {
      await DeleteEmployee(accessToken, deleteTarget.id);
      setDeleteTarget(null);
      showSuccess("Employee deleted");
      await loadEmployees();
    } catch (err) {
      showError(normalizeError(err));
    } finally {
      setDeleteSubmitting(false);
    }
  };

  return (
    <Paper elevation={0} sx={{ p: { xs: 2, md: 3 }, borderRadius: 2.5, border: "1px solid #dbe3ef", backgroundColor: "rgba(255,255,255,0.92)" }}>
      <Stack spacing={2}>
        <Stack direction={{ xs: "column", md: "row" }} spacing={1.2} justifyContent="space-between" alignItems={{ md: "center" }}>
          <Box>
            <Typography variant="h5" sx={{ fontWeight: 800 }}>Employees</Typography>
            <Typography variant="body2" color="text.secondary">Manage employee records, status, and departmental assignment.</Typography>
          </Box>
          <Button variant="contained" onClick={openCreateDialog}>New Employee</Button>
        </Stack>

        <Stack direction={{ xs: "column", md: "row" }} spacing={1.2}>
          <TextField
            label="Search name"
            size="small"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                setPage(0);
                setAppliedSearch(searchInput.trim());
              }
            }}
            sx={{ minWidth: 220 }}
          />
          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel>Department</InputLabel>
            <Select
              value={departmentFilter}
              label="Department"
              onChange={(e) => {
                setPage(0);
                setDepartmentFilter(e.target.value);
              }}
            >
              <MenuItem value="">All</MenuItem>
              {departments.map((department) => (
                <MenuItem key={department.id} value={String(department.id)}>{department.name}</MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel>Status</InputLabel>
            <Select
              value={statusFilter}
              label="Status"
              onChange={(e) => {
                setPage(0);
                setStatusFilter(e.target.value);
              }}
            >
              <MenuItem value="">All</MenuItem>
              {employmentStatuses.map((status) => (
                <MenuItem key={status} value={status}>{status}</MenuItem>
              ))}
            </Select>
          </FormControl>
          <Stack direction="row" spacing={1}>
            <Button variant="outlined" onClick={() => { setPage(0); setAppliedSearch(searchInput.trim()); }}>Search</Button>
            <Button
              variant="text"
              onClick={() => {
                setSearchInput("");
                setAppliedSearch("");
                setStatusFilter("");
                setDepartmentFilter("");
                setPage(0);
              }}
            >
              Reset
            </Button>
          </Stack>
        </Stack>

        <Box sx={{ overflowX: "auto", border: "1px solid #e4eaf3", borderRadius: 2 }}>
          <Table size="small">
            <TableHead>
              <TableRow sx={{ backgroundColor: "#f4f7fc" }}>
                <TableCell>Name</TableCell>
                <TableCell>Department</TableCell>
                <TableCell>Position</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Phone</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Base Salary</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {items.length === 0 && (
                <TableRow>
                  <TableCell colSpan={8}>
                    <Typography variant="body2" color="text.secondary" sx={{ py: 2, textAlign: "center" }}>
                      {loading ? "Loading..." : "No employees found"}
                    </Typography>
                  </TableCell>
                </TableRow>
              )}
              {items.map((employee) => (
                <TableRow key={employee.id} hover>
                  <TableCell>{employee.last_name}, {employee.first_name}</TableCell>
                  <TableCell>{employee.department_name || "-"}</TableCell>
                  <TableCell>{employee.position}</TableCell>
                  <TableCell>
                    <Chip size="small" label={employee.employment_status} color={employee.employment_status === "Active" ? "success" : "default"} />
                  </TableCell>
                  <TableCell>{employee.phone}</TableCell>
                  <TableCell>{employee.email || "-"}</TableCell>
                  <TableCell>{employee.base_salary.toFixed(2)}</TableCell>
                  <TableCell align="right">
                    <Stack direction="row" spacing={1} justifyContent="flex-end">
                      <Button size="small" onClick={() => void openEditDialog(employee.id)}>Edit</Button>
                      <Button size="small" color="error" onClick={() => setDeleteTarget(employee)}>Delete</Button>
                    </Stack>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Box>

        <TablePagination
          component="div"
          count={total}
          page={page}
          rowsPerPage={pageSize}
          onPageChange={(_event, nextPage) => setPage(nextPage)}
          onRowsPerPageChange={(event) => {
            setPage(0);
            setPageSize(Number(event.target.value));
          }}
          rowsPerPageOptions={[5, 10, 20, 50]}
        />
      </Stack>

      <Dialog open={formOpen} onClose={() => setFormOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>{editingEmployeeID == null ? "Create Employee" : "Edit Employee"}</DialogTitle>
        <DialogContent>
          <Stack spacing={1.2} sx={{ pt: 1 }}>
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.2}>
              <TextField label="First Name" value={form.first_name} onChange={(e) => setForm((prev) => ({ ...prev, first_name: e.target.value }))} fullWidth required />
              <TextField label="Last Name" value={form.last_name} onChange={(e) => setForm((prev) => ({ ...prev, last_name: e.target.value }))} fullWidth required />
            </Stack>
            <TextField label="Other Name" value={form.other_name} onChange={(e) => setForm((prev) => ({ ...prev, other_name: e.target.value }))} fullWidth />
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.2}>
              <FormControl fullWidth>
                <InputLabel>Gender</InputLabel>
                <Select value={form.gender} label="Gender" onChange={(e) => setForm((prev) => ({ ...prev, gender: e.target.value }))}>
                  {genders.map((gender) => (
                    <MenuItem key={gender} value={gender}>{gender}</MenuItem>
                  ))}
                </Select>
              </FormControl>
              <TextField label="Date of Birth" type="date" value={form.dob} onChange={(e) => setForm((prev) => ({ ...prev, dob: e.target.value }))} fullWidth InputLabelProps={{ shrink: true }} required />
            </Stack>
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.2}>
              <TextField label="Phone" value={form.phone} onChange={(e) => setForm((prev) => ({ ...prev, phone: e.target.value }))} fullWidth required />
              <TextField label="Email" value={form.email} onChange={(e) => setForm((prev) => ({ ...prev, email: e.target.value }))} fullWidth />
            </Stack>
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.2}>
              <TextField label="National ID" value={form.national_id} onChange={(e) => setForm((prev) => ({ ...prev, national_id: e.target.value }))} fullWidth />
              <FormControl fullWidth>
                <InputLabel>Department</InputLabel>
                <Select value={form.department_id} label="Department" onChange={(e) => setForm((prev) => ({ ...prev, department_id: e.target.value }))}>
                  <MenuItem value="">Unassigned</MenuItem>
                  {departments.map((department) => (
                    <MenuItem key={department.id} value={String(department.id)}>{department.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Stack>
            <TextField label="Address" value={form.address} onChange={(e) => setForm((prev) => ({ ...prev, address: e.target.value }))} fullWidth />
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.2}>
              <TextField label="Position" value={form.position} onChange={(e) => setForm((prev) => ({ ...prev, position: e.target.value }))} fullWidth required />
              <FormControl fullWidth>
                <InputLabel>Employment Status</InputLabel>
                <Select value={form.employment_status} label="Employment Status" onChange={(e) => setForm((prev) => ({ ...prev, employment_status: e.target.value }))}>
                  {employmentStatuses.map((status) => (
                    <MenuItem key={status} value={status}>{status}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Stack>
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.2}>
              <TextField label="Hire Date" type="date" value={form.hire_date} onChange={(e) => setForm((prev) => ({ ...prev, hire_date: e.target.value }))} fullWidth InputLabelProps={{ shrink: true }} required />
              <TextField label="Base Salary" type="number" value={form.base_salary} onChange={(e) => setForm((prev) => ({ ...prev, base_salary: e.target.value }))} fullWidth required />
            </Stack>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setFormOpen(false)} disabled={formSaving}>Cancel</Button>
          <Button onClick={() => void onSubmitForm()} variant="contained" disabled={formSaving}>Save</Button>
        </DialogActions>
      </Dialog>

      <Dialog open={Boolean(deleteTarget)} onClose={() => setDeleteTarget(null)}>
        <DialogTitle>Delete Employee</DialogTitle>
        <DialogContent>
          <Typography variant="body2">Delete {deleteTarget?.first_name} {deleteTarget?.last_name}? This action cannot be undone.</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteTarget(null)} disabled={deleteSubmitting}>Cancel</Button>
          <Button color="error" variant="contained" onClick={() => void onConfirmDelete()} disabled={deleteSubmitting}>Delete</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={3500}
        onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
      >
        <Alert severity={snackbar.severity} variant="filled" onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Paper>
  );
}
