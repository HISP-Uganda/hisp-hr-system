export type EmployeeView = {
  id: number;
  first_name: string;
  last_name: string;
  other_name: string;
  gender: string;
  dob: string;
  phone: string;
  email: string;
  national_id: string;
  address: string;
  department_id?: number;
  department_name: string;
  position: string;
  employment_status: string;
  hire_date: string;
  base_salary: number;
  created_at: string;
  updated_at: string;
};

export type EmployeeInput = {
  first_name: string;
  last_name: string;
  other_name: string;
  gender: string;
  dob: string;
  phone: string;
  email: string;
  national_id: string;
  address: string;
  department_id?: number;
  position: string;
  employment_status: string;
  hire_date: string;
  base_salary: number;
};

export type EmployeeListQuery = {
  search: string;
  department_id?: number;
  status: string;
  page: number;
  page_size: number;
};

export type EmployeeListResult = {
  items: EmployeeView[];
  total: number;
  page: number;
  page_size: number;
};

export type EmployeeResponse = {
  success: boolean;
  message: string;
  data: EmployeeView;
};

export type EmployeeListResponse = {
  success: boolean;
  message: string;
  data: EmployeeListResult;
};

export type DepartmentOption = {
  id: number;
  name: string;
};

export type DepartmentListResponse = {
  success: boolean;
  message: string;
  data: DepartmentOption[];
};
