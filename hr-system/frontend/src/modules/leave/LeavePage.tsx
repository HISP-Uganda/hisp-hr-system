import React, { useCallback, useEffect, useMemo, useState } from "react";
import {
  Alert,
  Box,
  Button,
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
  Tab,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Tabs,
  TextField,
  Typography,
} from "@mui/material";
import {
  AdminLeaveBalance,
  ApplyLeave,
  ApproveLeave,
  CancelLeave,
  ConvertAbsenceToLeave,
  CreateLeaveType,
  DeactivateLeaveType,
  ListEmployees,
  ListLeaveRequests,
  ListLeaveTypes,
  ListLockedLeaveDates,
  LockLeaveDate,
  MeLeaveBalance,
  MasterDeleteLeave,
  RejectLeave,
  UnlockLeaveDate,
} from "../../../wailsjs/go/main/App";
import { useAuth } from "../../auth/AuthContext";
import { EmployeeListResponse } from "../employees/types";
import {
  LeaveBalanceResponse,
  LeaveRequest,
  LeaveRequestListResponse,
  LeaveType,
  LeaveTypeListResponse,
  LockedDate,
  LockedDateListResponse,
} from "./types";

function normalizeError(err: unknown): string {
  if (typeof err === "string") return err;
  if (err instanceof Error) return err.message;
  return "request failed";
}

function computeWorkingDays(startDate: string, endDate: string): number {
  if (!startDate || !endDate) return 0;
  const start = new Date(`${startDate}T00:00:00Z`);
  const end = new Date(`${endDate}T00:00:00Z`);
  if (end < start) return 0;

  let days = 0;
  for (let d = new Date(start); d <= end; d.setUTCDate(d.getUTCDate() + 1)) {
    const wd = d.getUTCDay();
    if (wd !== 0 && wd !== 6) days += 1;
  }
  return days;
}

export function LeavePage() {
  const auth = useAuth();
  const accessToken = auth.accessToken;
  const role = auth.user?.role ?? "";
  const isAdmin = role === "Admin";
  const isHR = role === "HR Officer";
  const isMaster = role === "Master" || role === "Master Admin";
  const canManage = isAdmin || isHR || isMaster;

  const [tab, setTab] = useState(0);
  const [leaveTypes, setLeaveTypes] = useState<LeaveType[]>([]);
  const [lockedDates, setLockedDates] = useState<LockedDate[]>([]);
  const [requests, setRequests] = useState<LeaveRequest[]>([]);
  const [employees, setEmployees] = useState<{ id: number; first_name: string; last_name: string }[]>([]);
  const [year, setYear] = useState<number>(new Date().getUTCFullYear());
  const [balanceRows, setBalanceRows] = useState<LeaveBalanceResponse["data"]["items"]>([]);

  const [applyForm, setApplyForm] = useState({ employee_id: "", leave_type_id: "", start_date: "", end_date: "", comment: "" });
  const [filters, setFilters] = useState({ status: "", employee_id: "", leave_type_id: "", from_date: "", to_date: "" });

  const [lockDate, setLockDate] = useState("");
  const [lockReason, setLockReason] = useState("");

  const [newTypeOpen, setNewTypeOpen] = useState(false);
  const [newType, setNewType] = useState({ name: "", annual_entitlement_days: "0" });

  const [snackbar, setSnackbar] = useState({ open: false, severity: "success" as "success" | "error", message: "" });

  const showSuccess = (message: string) => setSnackbar({ open: true, severity: "success", message });
  const showError = (message: string) => setSnackbar({ open: true, severity: "error", message });

  const workingDaysPreview = useMemo(() => computeWorkingDays(applyForm.start_date, applyForm.end_date), [applyForm.start_date, applyForm.end_date]);

  const loadLeaveTypes = useCallback(async () => {
    if (!accessToken) return;
    try {
      const response = (await ListLeaveTypes(accessToken)) as LeaveTypeListResponse;
      setLeaveTypes(response.data);
    } catch (err) {
      showError(normalizeError(err));
    }
  }, [accessToken]);

  const loadLockedDates = useCallback(async () => {
    if (!accessToken) return;
    try {
      const response = (await ListLockedLeaveDates(accessToken, year)) as LockedDateListResponse;
      setLockedDates(response.data);
    } catch (err) {
      showError(normalizeError(err));
    }
  }, [accessToken, year]);

  const loadRequests = useCallback(async () => {
    if (!accessToken) return;
    try {
      const response = (await ListLeaveRequests(accessToken, {
        status: filters.status,
        employee_id: filters.employee_id ? Number(filters.employee_id) : undefined,
        leave_type_id: filters.leave_type_id ? Number(filters.leave_type_id) : undefined,
        from_date: filters.from_date,
        to_date: filters.to_date,
        page: 1,
        page_size: 100,
      })) as LeaveRequestListResponse;
      setRequests(response.data.items);
    } catch (err) {
      showError(normalizeError(err));
    }
  }, [accessToken, filters]);

  const loadEmployees = useCallback(async () => {
    if (!accessToken || !canManage) return;
    try {
      const response = (await ListEmployees(accessToken, { search: "", status: "", page: 1, page_size: 200 })) as EmployeeListResponse;
      setEmployees(response.data.items);
    } catch (err) {
      showError(normalizeError(err));
    }
  }, [accessToken, canManage]);

  const loadBalance = useCallback(async () => {
    if (!accessToken) return;
    try {
      const response = canManage
        ? ((await AdminLeaveBalance(accessToken, filters.employee_id ? Number(filters.employee_id) : (employees[0]?.id ?? 0), year)) as LeaveBalanceResponse)
        : ((await MeLeaveBalance(accessToken, year)) as LeaveBalanceResponse);
      setBalanceRows(response.data.items);
    } catch (err) {
      setBalanceRows([]);
      showError(normalizeError(err));
    }
  }, [accessToken, canManage, filters.employee_id, employees, year]);

  useEffect(() => { void loadLeaveTypes(); }, [loadLeaveTypes]);
  useEffect(() => { void loadLockedDates(); }, [loadLockedDates]);
  useEffect(() => { void loadRequests(); }, [loadRequests]);
  useEffect(() => { void loadEmployees(); }, [loadEmployees]);
  useEffect(() => { void loadBalance(); }, [loadBalance]);

  const onApply = async () => {
    if (!accessToken) return;
    if (!applyForm.leave_type_id || !applyForm.start_date || !applyForm.end_date) {
      showError("Leave type and date range are required");
      return;
    }
    if (workingDaysPreview <= 0) {
      showError("Selected range has no working days");
      return;
    }

    try {
      await ApplyLeave(accessToken, {
        employee_id: canManage && applyForm.employee_id ? Number(applyForm.employee_id) : undefined,
        leave_type_id: Number(applyForm.leave_type_id),
        start_date: applyForm.start_date,
        end_date: applyForm.end_date,
        comment: applyForm.comment,
      });
      showSuccess("Leave request submitted");
      setApplyForm({ employee_id: "", leave_type_id: "", start_date: "", end_date: "", comment: "" });
      await loadRequests();
      await loadBalance();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onApprove = async (id: number) => {
    if (!accessToken) return;
    try {
      await ApproveLeave(accessToken, id, { comment: "approved" });
      showSuccess("Leave approved");
      await loadRequests();
      await loadBalance();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onReject = async (id: number) => {
    if (!accessToken) return;
    try {
      await RejectLeave(accessToken, id, { comment: "rejected" });
      showSuccess("Leave rejected");
      await loadRequests();
      await loadBalance();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onCancel = async (id: number) => {
    if (!accessToken) return;
    try {
      await CancelLeave(accessToken, id, { comment: "cancelled" });
      showSuccess("Leave cancelled");
      await loadRequests();
      await loadBalance();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onLockDate = async () => {
    if (!accessToken || !lockDate) return;
    try {
      await LockLeaveDate(accessToken, { date: lockDate, reason: lockReason });
      setLockDate("");
      setLockReason("");
      showSuccess("Date locked");
      await loadLockedDates();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onUnlockDate = async (date: string) => {
    if (!accessToken) return;
    try {
      await UnlockLeaveDate(accessToken, date);
      showSuccess("Date unlocked");
      await loadLockedDates();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onCreateType = async () => {
    if (!accessToken) return;
    try {
      await CreateLeaveType(accessToken, {
        name: newType.name,
        annual_entitlement_days: Number(newType.annual_entitlement_days),
        is_paid: true,
        requires_attachment: false,
        requires_approval: true,
        counts_toward_entitlement: true,
        is_active: true,
      });
      setNewTypeOpen(false);
      setNewType({ name: "", annual_entitlement_days: "0" });
      showSuccess("Leave type created");
      await loadLeaveTypes();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onDeactivateType = async (id: number) => {
    if (!accessToken) return;
    try {
      await DeactivateLeaveType(accessToken, id);
      showSuccess("Leave type deactivated");
      await loadLeaveTypes();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onMasterDelete = async (id: number) => {
    if (!accessToken) return;
    try {
      await MasterDeleteLeave(accessToken, id);
      showSuccess("Leave deleted");
      await loadRequests();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  const onConvertAbsence = async () => {
    if (!accessToken || !filters.employee_id || !applyForm.leave_type_id || !applyForm.start_date) return;
    try {
      await ConvertAbsenceToLeave(accessToken, Number(filters.employee_id), applyForm.start_date, Number(applyForm.leave_type_id));
      showSuccess("Absence converted to leave");
      await loadRequests();
      await loadBalance();
    } catch (err) {
      showError(normalizeError(err));
    }
  };

  return (
    <Paper elevation={0} sx={{ p: { xs: 2, md: 3 }, borderRadius: 2.5, border: "1px solid #dbe3ef", backgroundColor: "rgba(255,255,255,0.92)" }}>
      <Stack spacing={2}>
        <Box>
          <Typography variant="h5" sx={{ fontWeight: 800 }}>Leave Management</Typography>
          <Typography variant="body2" color="text.secondary">Planner, applications, approvals, balances, and locked-date controls.</Typography>
        </Box>

        <Tabs value={tab} onChange={(_, v) => setTab(v)}>
          <Tab label="Planner" />
          <Tab label="Apply" />
          <Tab label="Requests" />
          <Tab label="Balances" />
        </Tabs>

        {tab === 0 && (
          <Stack spacing={1.2}>
            <Stack direction="row" spacing={1.2} alignItems="center">
              <TextField label="Year" type="number" size="small" value={year} onChange={(e) => setYear(Number(e.target.value))} sx={{ width: 120 }} />
              <Button variant="outlined" onClick={() => void loadLockedDates()}>Reload</Button>
            </Stack>
            {canManage && (
              <Stack direction={{ xs: "column", md: "row" }} spacing={1.2}>
                <TextField type="date" label="Lock Date" size="small" value={lockDate} onChange={(e) => setLockDate(e.target.value)} InputLabelProps={{ shrink: true }} />
                <TextField label="Reason" size="small" value={lockReason} onChange={(e) => setLockReason(e.target.value)} sx={{ minWidth: 240 }} />
                <Button variant="contained" onClick={onLockDate}>Lock</Button>
              </Stack>
            )}
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Locked Date</TableCell>
                  <TableCell>Reason</TableCell>
                  <TableCell align="right">Action</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {lockedDates.map((row) => (
                  <TableRow key={row.id}>
                    <TableCell>{row.lock_date.slice(0, 10)}</TableCell>
                    <TableCell>{row.reason || "-"}</TableCell>
                    <TableCell align="right">{canManage && <Button size="small" color="error" onClick={() => void onUnlockDate(row.lock_date.slice(0, 10))}>Unlock</Button>}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Stack>
        )}

        {tab === 1 && (
          <Stack spacing={1.2}>
            {canManage && (
              <FormControl size="small" sx={{ maxWidth: 260 }}>
                <InputLabel>Employee</InputLabel>
                <Select value={applyForm.employee_id} label="Employee" onChange={(e) => setApplyForm((prev) => ({ ...prev, employee_id: e.target.value }))}>
                  {employees.map((employee) => (
                    <MenuItem key={employee.id} value={String(employee.id)}>{employee.last_name}, {employee.first_name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}
            <Stack direction={{ xs: "column", md: "row" }} spacing={1.2}>
              <FormControl size="small" fullWidth>
                <InputLabel>Leave Type</InputLabel>
                <Select value={applyForm.leave_type_id} label="Leave Type" onChange={(e) => setApplyForm((prev) => ({ ...prev, leave_type_id: e.target.value }))}>
                  {leaveTypes.filter((t) => t.is_active).map((leaveType) => (
                    <MenuItem key={leaveType.id} value={String(leaveType.id)}>{leaveType.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
              <TextField type="date" label="Start" size="small" value={applyForm.start_date} onChange={(e) => setApplyForm((prev) => ({ ...prev, start_date: e.target.value }))} InputLabelProps={{ shrink: true }} fullWidth />
              <TextField type="date" label="End" size="small" value={applyForm.end_date} onChange={(e) => setApplyForm((prev) => ({ ...prev, end_date: e.target.value }))} InputLabelProps={{ shrink: true }} fullWidth />
            </Stack>
            <TextField label="Comment" multiline minRows={2} value={applyForm.comment} onChange={(e) => setApplyForm((prev) => ({ ...prev, comment: e.target.value }))} />
            <Typography variant="body2" color="text.secondary">Working days preview (excluding weekends): <strong>{workingDaysPreview}</strong></Typography>
            <Stack direction="row" spacing={1.2}>
              <Button variant="contained" onClick={() => void onApply()}>Submit Leave</Button>
              {canManage && <Button variant="outlined" onClick={() => void setNewTypeOpen(true)}>Add Leave Type</Button>}
              {canManage && <Button variant="outlined" onClick={() => void onConvertAbsence()}>Convert Absence (1 day)</Button>}
            </Stack>
            {canManage && (
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Type</TableCell>
                    <TableCell>Entitlement</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell align="right">Action</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {leaveTypes.map((type) => (
                    <TableRow key={type.id}>
                      <TableCell>{type.name}</TableCell>
                      <TableCell>{type.annual_entitlement_days}</TableCell>
                      <TableCell>{type.is_active ? "Active" : "Inactive"}</TableCell>
                      <TableCell align="right">{type.is_active && <Button size="small" color="error" onClick={() => void onDeactivateType(type.id)}>Deactivate</Button>}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </Stack>
        )}

        {tab === 2 && (
          <Stack spacing={1.2}>
            <Stack direction={{ xs: "column", md: "row" }} spacing={1.2}>
              <FormControl size="small" sx={{ minWidth: 160 }}>
                <InputLabel>Status</InputLabel>
                <Select value={filters.status} label="Status" onChange={(e) => setFilters((prev) => ({ ...prev, status: e.target.value }))}>
                  <MenuItem value="">All</MenuItem>
                  <MenuItem value="Pending">Pending</MenuItem>
                  <MenuItem value="Approved">Approved</MenuItem>
                  <MenuItem value="Rejected">Rejected</MenuItem>
                  <MenuItem value="Cancelled">Cancelled</MenuItem>
                </Select>
              </FormControl>
              {canManage && (
                <FormControl size="small" sx={{ minWidth: 220 }}>
                  <InputLabel>Employee</InputLabel>
                  <Select value={filters.employee_id} label="Employee" onChange={(e) => setFilters((prev) => ({ ...prev, employee_id: e.target.value }))}>
                    <MenuItem value="">All</MenuItem>
                    {employees.map((employee) => (
                      <MenuItem key={employee.id} value={String(employee.id)}>{employee.last_name}, {employee.first_name}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
              <Button variant="outlined" onClick={() => void loadRequests()}>Apply Filters</Button>
            </Stack>

            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Employee</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>Period</TableCell>
                  <TableCell>Days</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {requests.map((request) => (
                  <TableRow key={request.id}>
                    <TableCell>{request.employee_name || request.employee_id}</TableCell>
                    <TableCell>{request.type_name || request.leave_type_id}</TableCell>
                    <TableCell>{request.start_date.slice(0, 10)} - {request.end_date.slice(0, 10)}</TableCell>
                    <TableCell>{request.working_days}</TableCell>
                    <TableCell>{request.status}</TableCell>
                    <TableCell align="right">
                      <Stack direction="row" spacing={1} justifyContent="flex-end">
                        {canManage && request.status === "Pending" && <Button size="small" onClick={() => void onApprove(request.id)}>Approve</Button>}
                        {canManage && request.status === "Pending" && <Button size="small" color="warning" onClick={() => void onReject(request.id)}>Reject</Button>}
                        {(request.status === "Pending" || (canManage && request.status === "Approved")) && <Button size="small" color="error" onClick={() => void onCancel(request.id)}>Cancel</Button>}
                        {isMaster && <Button size="small" color="error" onClick={() => void onMasterDelete(request.id)}>Delete</Button>}
                      </Stack>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Stack>
        )}

        {tab === 3 && (
          <Stack spacing={1.2}>
            <Stack direction="row" spacing={1.2}>
              <TextField type="number" size="small" label="Year" value={year} onChange={(e) => setYear(Number(e.target.value))} sx={{ width: 140 }} />
              {canManage && (
                <FormControl size="small" sx={{ minWidth: 220 }}>
                  <InputLabel>Employee</InputLabel>
                  <Select value={filters.employee_id} label="Employee" onChange={(e) => setFilters((prev) => ({ ...prev, employee_id: e.target.value }))}>
                    {employees.map((employee) => (
                      <MenuItem key={employee.id} value={String(employee.id)}>{employee.last_name}, {employee.first_name}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
              <Button variant="outlined" onClick={() => void loadBalance()}>Refresh</Button>
            </Stack>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Leave Type</TableCell>
                  <TableCell>Total</TableCell>
                  <TableCell>Reserved</TableCell>
                  <TableCell>Pending</TableCell>
                  <TableCell>Approved</TableCell>
                  <TableCell>Available</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {balanceRows.map((row) => (
                  <TableRow key={`${row.leave_type_id}-${row.type_name}`}>
                    <TableCell>{row.type_name}</TableCell>
                    <TableCell>{row.total}</TableCell>
                    <TableCell>{row.reserved}</TableCell>
                    <TableCell>{row.pending}</TableCell>
                    <TableCell>{row.approved}</TableCell>
                    <TableCell>{row.available}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Stack>
        )}
      </Stack>

      <Dialog open={newTypeOpen} onClose={() => setNewTypeOpen(false)}>
        <DialogTitle>Create Leave Type</DialogTitle>
        <DialogContent>
          <Stack spacing={1.2} sx={{ pt: 1 }}>
            <TextField label="Name" value={newType.name} onChange={(e) => setNewType((prev) => ({ ...prev, name: e.target.value }))} />
            <TextField label="Annual Entitlement" type="number" value={newType.annual_entitlement_days} onChange={(e) => setNewType((prev) => ({ ...prev, annual_entitlement_days: e.target.value }))} />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setNewTypeOpen(false)}>Cancel</Button>
          <Button variant="contained" onClick={() => void onCreateType()}>Create</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={3500}
        onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
      >
        <Alert severity={snackbar.severity} variant="filled" onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}>{snackbar.message}</Alert>
      </Snackbar>
    </Paper>
  );
}
