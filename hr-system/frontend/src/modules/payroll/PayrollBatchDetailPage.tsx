import React, { useEffect, useMemo, useState } from "react";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import { Link, useNavigate, useParams } from "@tanstack/react-router";
import {
  ApprovePayrollBatch,
  ExportPayrollBatchCSV,
  GeneratePayrollEntries,
  GetPayrollBatch,
  LockPayrollBatch,
  UpdatePayrollEntryAmounts,
} from "../../../wailsjs/go/main/App";
import { useAuth } from "../../auth/AuthContext";
import {
  PayrollBatchDetailResponse,
  PayrollCSVResponse,
  PayrollEntry,
  PayrollEntryResponse,
} from "./types";

type EntryDraftValues = {
  allowances_total: string;
  deductions_total: string;
  tax_total: string;
};

function toDraftValues(entry: PayrollEntry): EntryDraftValues {
  return {
    allowances_total: String(entry.allowances_total),
    deductions_total: String(entry.deductions_total),
    tax_total: String(entry.tax_total),
  };
}

export function PayrollBatchDetailPage() {
  const { accessToken } = useAuth();
  const navigate = useNavigate();
  const params = useParams({ from: "/payroll/$batchId" });
  const batchID = Number(params.batchId);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [detail, setDetail] = useState<PayrollBatchDetailResponse["data"] | null>(null);
  const [draftValues, setDraftValues] = useState<Record<number, EntryDraftValues>>({});

  const status = detail?.batch.status ?? "Draft";
  const canEdit = status === "Draft";

  const loadBatch = async () => {
    if (!accessToken || !Number.isFinite(batchID) || batchID <= 0) return;
    setLoading(true);
    setError(null);
    try {
      const response = (await GetPayrollBatch(accessToken, batchID)) as PayrollBatchDetailResponse;
      setDetail(response.data);
      const next: Record<number, EntryDraftValues> = {};
      response.data.entries.forEach((entry) => {
        next[entry.id] = toDraftValues(entry);
      });
      setDraftValues(next);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load payroll batch");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadBatch();
  }, [accessToken, batchID]);

  const onGenerate = async () => {
    if (!accessToken) return;
    setError(null);
    try {
      await GeneratePayrollEntries(accessToken, batchID);
      await loadBatch();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to generate payroll entries");
    }
  };

  const onApprove = async () => {
    if (!accessToken) return;
    setError(null);
    try {
      await ApprovePayrollBatch(accessToken, batchID);
      await loadBatch();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to approve payroll batch");
    }
  };

  const onLock = async () => {
    if (!accessToken) return;
    setError(null);
    try {
      await LockPayrollBatch(accessToken, batchID);
      await loadBatch();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to lock payroll batch");
    }
  };

  const onSaveRow = async (entry: PayrollEntry) => {
    if (!accessToken || !canEdit) return;
    const values = draftValues[entry.id] ?? toDraftValues(entry);

    const allowances = Number(values.allowances_total);
    const deductions = Number(values.deductions_total);
    const tax = Number(values.tax_total);
    if (!Number.isFinite(allowances) || !Number.isFinite(deductions) || !Number.isFinite(tax)) {
      setError("invalid numeric values");
      return;
    }

    setError(null);
    try {
      await UpdatePayrollEntryAmounts(accessToken, entry.id, {
        allowances_total: allowances,
        deductions_total: deductions,
        tax_total: tax,
      }) as PayrollEntryResponse;
      await loadBatch();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to update payroll entry");
    }
  };

  const onExportCSV = async () => {
    if (!accessToken) return;
    setError(null);
    try {
      const response = (await ExportPayrollBatchCSV(accessToken, batchID)) as PayrollCSVResponse;
      const blob = new Blob([response.data], { type: "text/csv;charset=utf-8" });
      const url = window.URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = url;
      anchor.download = `payroll-${detail?.batch.month ?? batchID}.csv`;
      anchor.click();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to export payroll csv");
    }
  };

  const headerMeta = useMemo(() => {
    if (!detail) return null;
    return {
      createdAt: new Date(detail.batch.created_at).toLocaleString(),
      approvedAt: detail.batch.approved_at ? new Date(detail.batch.approved_at).toLocaleString() : "-",
      lockedAt: detail.batch.locked_at ? new Date(detail.batch.locked_at).toLocaleString() : "-",
    };
  }, [detail]);

  return (
    <Stack spacing={2.5}>
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Box>
          <Typography variant="h5" sx={{ fontWeight: 800 }}>Payroll Batch Details</Typography>
          <Typography variant="body2" color="text.secondary">
            <Link to="/payroll">Back to batches</Link>
          </Typography>
        </Box>
        <Button variant="outlined" onClick={() => void navigate({ to: "/payroll" })}>Back</Button>
      </Stack>

      {error && <Alert severity="error">{error}</Alert>}

      {!detail && !loading && <Alert severity="info">No batch data available.</Alert>}

      {detail && (
        <>
          <Card variant="outlined">
            <CardContent>
              <Stack spacing={1.2}>
                <Stack direction="row" spacing={1.1} alignItems="center">
                  <Typography variant="h6" sx={{ fontWeight: 700 }}>{detail.batch.month}</Typography>
                  <Chip
                    size="small"
                    label={detail.batch.status}
                    color={detail.batch.status === "Draft" ? "warning" : detail.batch.status === "Approved" ? "info" : "success"}
                  />
                </Stack>
                <Typography variant="body2" color="text.secondary">Created: {headerMeta?.createdAt}</Typography>
                <Typography variant="body2" color="text.secondary">Approved: {headerMeta?.approvedAt}</Typography>
                <Typography variant="body2" color="text.secondary">Locked: {headerMeta?.lockedAt}</Typography>
                <Stack direction={{ xs: "column", md: "row" }} spacing={1.1}>
                  {status === "Draft" && <Button variant="contained" onClick={() => void onGenerate()}>Generate Entries</Button>}
                  {status === "Draft" && <Button variant="outlined" onClick={() => void onApprove()}>Approve Batch</Button>}
                  {status === "Approved" && <Button variant="outlined" color="warning" onClick={() => void onLock()}>Lock Batch</Button>}
                  {(status === "Approved" || status === "Locked") && <Button variant="outlined" onClick={() => void onExportCSV()}>Export CSV</Button>}
                </Stack>
              </Stack>
            </CardContent>
          </Card>

          <Card variant="outlined">
            <CardContent sx={{ p: 0 }}>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Employee</TableCell>
                    <TableCell>Base</TableCell>
                    <TableCell>Allowances</TableCell>
                    <TableCell>Deductions</TableCell>
                    <TableCell>Tax</TableCell>
                    <TableCell>Gross</TableCell>
                    <TableCell>Net</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {detail.entries.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={8}>{loading ? "Loading..." : "No entries generated yet."}</TableCell>
                    </TableRow>
                  )}
                  {detail.entries.map((entry) => {
                    const values = draftValues[entry.id] ?? toDraftValues(entry);
                    return (
                      <TableRow key={entry.id} hover>
                        <TableCell>{entry.employee_name}</TableCell>
                        <TableCell>{entry.base_salary.toFixed(2)}</TableCell>
                        <TableCell>
                          <TextField
                            size="small"
                            type="number"
                            value={values.allowances_total}
                            onChange={(e) => setDraftValues((prev) => ({ ...prev, [entry.id]: { ...values, allowances_total: e.target.value } }))}
                            disabled={!canEdit}
                            sx={{ width: 110 }}
                          />
                        </TableCell>
                        <TableCell>
                          <TextField
                            size="small"
                            type="number"
                            value={values.deductions_total}
                            onChange={(e) => setDraftValues((prev) => ({ ...prev, [entry.id]: { ...values, deductions_total: e.target.value } }))}
                            disabled={!canEdit}
                            sx={{ width: 110 }}
                          />
                        </TableCell>
                        <TableCell>
                          <TextField
                            size="small"
                            type="number"
                            value={values.tax_total}
                            onChange={(e) => setDraftValues((prev) => ({ ...prev, [entry.id]: { ...values, tax_total: e.target.value } }))}
                            disabled={!canEdit}
                            sx={{ width: 110 }}
                          />
                        </TableCell>
                        <TableCell>{entry.gross_pay.toFixed(2)}</TableCell>
                        <TableCell>{entry.net_pay.toFixed(2)}</TableCell>
                        <TableCell align="right">
                          {canEdit ? <Button size="small" onClick={() => void onSaveRow(entry)}>Save</Button> : "-"}
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </>
      )}
    </Stack>
  );
}
