export type PayrollBatch = {
  id: number;
  month: string;
  status: "Draft" | "Approved" | "Locked";
  created_by: number;
  created_at: string;
  approved_by?: number;
  approved_at?: string;
  locked_at?: string;
};

export type PayrollEntry = {
  id: number;
  batch_id: number;
  employee_id: number;
  employee_name: string;
  base_salary: number;
  allowances_total: number;
  deductions_total: number;
  tax_total: number;
  gross_pay: number;
  net_pay: number;
  created_at: string;
  updated_at: string;
};

export type PayrollBatchDetail = {
  batch: PayrollBatch;
  entries: PayrollEntry[];
};

export type PayrollBatchListResult = {
  items: PayrollBatch[];
};

export type PayrollBatchFilter = {
  month: string;
  status: string;
};

export type PayrollCreateBatchInput = {
  month: string;
};

export type PayrollUpdateEntryAmountsInput = {
  allowances_total: number;
  deductions_total: number;
  tax_total: number;
};

export type PayrollBatchListResponse = {
  success: boolean;
  message: string;
  data: PayrollBatchListResult;
};

export type PayrollBatchDetailResponse = {
  success: boolean;
  message: string;
  data: PayrollBatchDetail;
};

export type PayrollBatchResponse = {
  success: boolean;
  message: string;
  data: PayrollBatch;
};

export type PayrollEntryResponse = {
  success: boolean;
  message: string;
  data: PayrollEntry;
};

export type PayrollCSVResponse = {
  success: boolean;
  message: string;
  data: string;
};
