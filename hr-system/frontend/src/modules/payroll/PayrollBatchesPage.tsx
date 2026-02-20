import React, { useEffect, useMemo, useState } from "react";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import { useNavigate } from "@tanstack/react-router";
import { CreatePayrollBatch, ListPayrollBatches } from "../../../wailsjs/go/main/App";
import { useAuth } from "../../auth/AuthContext";
import { PayrollBatchListResponse, PayrollBatchResponse } from "./types";

export function PayrollBatchesPage() {
  const { accessToken } = useAuth();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState("");
  const [createMonth, setCreateMonth] = useState("");
  const [batches, setBatches] = useState<PayrollBatchListResponse["data"]["items"]>([]);

  const filters = useMemo(() => ({ month: "", status: statusFilter }), [statusFilter]);

  const loadBatches = async () => {
    if (!accessToken) return;
    setLoading(true);
    setError(null);
    try {
      const response = (await ListPayrollBatches(accessToken, filters)) as PayrollBatchListResponse;
      setBatches(response.data.items);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load payroll batches");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadBatches();
  }, [accessToken, statusFilter]);

  const onCreateBatch = async () => {
    if (!accessToken || !createMonth) return;
    setError(null);
    try {
      await CreatePayrollBatch(accessToken, { month: createMonth }) as PayrollBatchResponse;
      setCreateMonth("");
      await loadBatches();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create payroll batch");
    }
  };

  const goToBatch = async (batchID: number) => {
    await navigate({ to: "/payroll/$batchId", params: { batchId: String(batchID) } });
  };

  return (
    <Stack spacing={2.5}>
      <Box>
        <Typography variant="h5" sx={{ fontWeight: 800 }}>Payroll Batches</Typography>
        <Typography variant="body2" color="text.secondary">Create monthly payroll batches and open details for generation, review, approval, and lock.</Typography>
      </Box>

      {error && <Alert severity="error">{error}</Alert>}

      <Card variant="outlined">
        <CardContent>
          <Stack direction={{ xs: "column", md: "row" }} spacing={1.2} alignItems={{ xs: "stretch", md: "center" }}>
            <TextField
              type="month"
              label="Batch Month"
              value={createMonth}
              onChange={(e) => setCreateMonth(e.target.value)}
              InputLabelProps={{ shrink: true }}
            />
            <Button variant="contained" onClick={() => void onCreateBatch()} disabled={!createMonth}>Create Batch</Button>
            <Box sx={{ flexGrow: 1 }} />
            <FormControl size="small" sx={{ minWidth: 180 }}>
              <InputLabel>Status Filter</InputLabel>
              <Select value={statusFilter} label="Status Filter" onChange={(e) => setStatusFilter(e.target.value)}>
                <MenuItem value="">All</MenuItem>
                <MenuItem value="Draft">Draft</MenuItem>
                <MenuItem value="Approved">Approved</MenuItem>
                <MenuItem value="Locked">Locked</MenuItem>
              </Select>
            </FormControl>
          </Stack>
        </CardContent>
      </Card>

      <Card variant="outlined">
        <CardContent sx={{ p: 0 }}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Month</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Created</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {batches.length === 0 && (
                <TableRow>
                  <TableCell colSpan={4}>{loading ? "Loading..." : "No payroll batches found"}</TableCell>
                </TableRow>
              )}
              {batches.map((batch) => (
                <TableRow key={batch.id} hover>
                  <TableCell>{batch.month}</TableCell>
                  <TableCell>
                    <Chip
                      size="small"
                      label={batch.status}
                      color={batch.status === "Draft" ? "warning" : batch.status === "Approved" ? "info" : "success"}
                    />
                  </TableCell>
                  <TableCell>{new Date(batch.created_at).toLocaleString()}</TableCell>
                  <TableCell align="right">
                    <Button size="small" onClick={() => void goToBatch(batch.id)}>Open</Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </Stack>
  );
}
