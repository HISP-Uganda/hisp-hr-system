export type LeaveType = {
  id: number;
  name: string;
  annual_entitlement_days: number;
  is_paid: boolean;
  requires_attachment: boolean;
  requires_approval: boolean;
  counts_toward_entitlement: boolean;
  is_active: boolean;
};

export type LeaveRequest = {
  id: number;
  employee_id: number;
  leave_type_id: number;
  start_date: string;
  end_date: string;
  working_days: number;
  status: string;
  comment: string;
  employee_name?: string;
  type_name?: string;
  department_name?: string;
};

export type LockedDate = {
  id: number;
  lock_date: string;
  reason: string;
};

export type LeaveRequestList = {
  items: LeaveRequest[];
  total: number;
  page: number;
  page_size: number;
};

export type LeaveBalanceItem = {
  employee_id: number;
  year: number;
  leave_type_id: number;
  type_name: string;
  total: number;
  reserved: number;
  pending: number;
  approved: number;
  available: number;
  used_percent: number;
};

export type LeaveBalanceSummary = {
  employee_id: number;
  year: number;
  items: LeaveBalanceItem[];
};

export type LeaveTypeListResponse = { success: boolean; message: string; data: LeaveType[] };
export type LeaveTypeResponse = { success: boolean; message: string; data: LeaveType };
export type LeaveRequestResponse = { success: boolean; message: string; data: LeaveRequest };
export type LeaveRequestListResponse = { success: boolean; message: string; data: LeaveRequestList };
export type LeaveBalanceResponse = { success: boolean; message: string; data: LeaveBalanceSummary };
export type LockedDateListResponse = { success: boolean; message: string; data: LockedDate[] };
export type LockedDateResponse = { success: boolean; message: string; data: LockedDate };
