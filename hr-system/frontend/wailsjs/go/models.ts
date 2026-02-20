export namespace bootstrap {
	
	export class AuthUser {
	    id: number;
	    username: string;
	    role: string;
	    is_active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AuthUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.role = source["role"];
	        this.is_active = source["is_active"];
	    }
	}
	export class AuthResult {
	    access_token: string;
	    refresh_token: string;
	    access_expiry: string;
	    user: AuthUser;
	
	    static createFrom(source: any = {}) {
	        return new AuthResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.access_token = source["access_token"];
	        this.refresh_token = source["refresh_token"];
	        this.access_expiry = source["access_expiry"];
	        this.user = this.convertValues(source["user"], AuthUser);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace employees {
	
	export class DepartmentOption {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new DepartmentOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class EmployeeListFilter {
	    search: string;
	    department_id?: number;
	    status: string;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new EmployeeListFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search = source["search"];
	        this.department_id = source["department_id"];
	        this.status = source["status"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	}
	export class EmployeeView {
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
	
	    static createFrom(source: any = {}) {
	        return new EmployeeView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.first_name = source["first_name"];
	        this.last_name = source["last_name"];
	        this.other_name = source["other_name"];
	        this.gender = source["gender"];
	        this.dob = source["dob"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.national_id = source["national_id"];
	        this.address = source["address"];
	        this.department_id = source["department_id"];
	        this.department_name = source["department_name"];
	        this.position = source["position"];
	        this.employment_status = source["employment_status"];
	        this.hire_date = source["hire_date"];
	        this.base_salary = source["base_salary"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class EmployeeListResult {
	    items: EmployeeView[];
	    total: number;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new EmployeeListResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], EmployeeView);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class UpsertEmployeeInput {
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
	
	    static createFrom(source: any = {}) {
	        return new UpsertEmployeeInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.first_name = source["first_name"];
	        this.last_name = source["last_name"];
	        this.other_name = source["other_name"];
	        this.gender = source["gender"];
	        this.dob = source["dob"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.national_id = source["national_id"];
	        this.address = source["address"];
	        this.department_id = source["department_id"];
	        this.position = source["position"];
	        this.employment_status = source["employment_status"];
	        this.hire_date = source["hire_date"];
	        this.base_salary = source["base_salary"];
	    }
	}

}

export namespace leave {
	
	export class ApplyInput {
	    employee_id?: number;
	    leave_type_id: number;
	    start_date: string;
	    end_date: string;
	    comment: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplyInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employee_id = source["employee_id"];
	        this.leave_type_id = source["leave_type_id"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.comment = source["comment"];
	    }
	}
	export class Balance {
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
	
	    static createFrom(source: any = {}) {
	        return new Balance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employee_id = source["employee_id"];
	        this.year = source["year"];
	        this.leave_type_id = source["leave_type_id"];
	        this.type_name = source["type_name"];
	        this.total = source["total"];
	        this.reserved = source["reserved"];
	        this.pending = source["pending"];
	        this.approved = source["approved"];
	        this.available = source["available"];
	        this.used_percent = source["used_percent"];
	    }
	}
	export class BalanceSummary {
	    employee_id: number;
	    year: number;
	    items: Balance[];
	
	    static createFrom(source: any = {}) {
	        return new BalanceSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.employee_id = source["employee_id"];
	        this.year = source["year"];
	        this.items = this.convertValues(source["items"], Balance);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DecisionInput {
	    comment: string;
	
	    static createFrom(source: any = {}) {
	        return new DecisionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.comment = source["comment"];
	    }
	}
	export class LeaveRequest {
	    id: number;
	    employee_id: number;
	    leave_type_id: number;
	    // Go type: time
	    start_date: any;
	    // Go type: time
	    end_date: any;
	    working_days: number;
	    status: string;
	    requested_by?: number;
	    approved_by?: number;
	    // Go type: time
	    approved_at?: any;
	    rejected_by?: number;
	    // Go type: time
	    rejected_at?: any;
	    cancelled_by?: number;
	    // Go type: time
	    cancelled_at?: any;
	    comment: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    employee_name: string;
	    department_id?: number;
	    department_name: string;
	    type_name: string;
	
	    static createFrom(source: any = {}) {
	        return new LeaveRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.employee_id = source["employee_id"];
	        this.leave_type_id = source["leave_type_id"];
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.end_date = this.convertValues(source["end_date"], null);
	        this.working_days = source["working_days"];
	        this.status = source["status"];
	        this.requested_by = source["requested_by"];
	        this.approved_by = source["approved_by"];
	        this.approved_at = this.convertValues(source["approved_at"], null);
	        this.rejected_by = source["rejected_by"];
	        this.rejected_at = this.convertValues(source["rejected_at"], null);
	        this.cancelled_by = source["cancelled_by"];
	        this.cancelled_at = this.convertValues(source["cancelled_at"], null);
	        this.comment = source["comment"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.employee_name = source["employee_name"];
	        this.department_id = source["department_id"];
	        this.department_name = source["department_name"];
	        this.type_name = source["type_name"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveType {
	    id: number;
	    name: string;
	    annual_entitlement_days: number;
	    is_paid: boolean;
	    requires_attachment: boolean;
	    requires_approval: boolean;
	    counts_toward_entitlement: boolean;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new LeaveType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.annual_entitlement_days = source["annual_entitlement_days"];
	        this.is_paid = source["is_paid"];
	        this.requires_attachment = source["requires_attachment"];
	        this.requires_approval = source["requires_approval"];
	        this.counts_toward_entitlement = source["counts_toward_entitlement"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveTypeInput {
	    name: string;
	    annual_entitlement_days: number;
	    is_paid: boolean;
	    requires_attachment: boolean;
	    requires_approval: boolean;
	    counts_toward_entitlement: boolean;
	    is_active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LeaveTypeInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.annual_entitlement_days = source["annual_entitlement_days"];
	        this.is_paid = source["is_paid"];
	        this.requires_attachment = source["requires_attachment"];
	        this.requires_approval = source["requires_approval"];
	        this.counts_toward_entitlement = source["counts_toward_entitlement"];
	        this.is_active = source["is_active"];
	    }
	}
	export class LockDateInput {
	    date: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new LockDateInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.reason = source["reason"];
	    }
	}
	export class LockedDate {
	    id: number;
	    // Go type: time
	    lock_date: any;
	    reason: string;
	    created_by?: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new LockedDate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.lock_date = this.convertValues(source["lock_date"], null);
	        this.reason = source["reason"];
	        this.created_by = source["created_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RequestFilter {
	    from_date: string;
	    to_date: string;
	    department_id?: number;
	    employee_id?: number;
	    leave_type_id?: number;
	    status: string;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new RequestFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from_date = source["from_date"];
	        this.to_date = source["to_date"];
	        this.department_id = source["department_id"];
	        this.employee_id = source["employee_id"];
	        this.leave_type_id = source["leave_type_id"];
	        this.status = source["status"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	}
	export class RequestList {
	    items: LeaveRequest[];
	    total: number;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new RequestList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], LeaveRequest);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class AuthResponse {
	    success: boolean;
	    message: string;
	    data: bootstrap.AuthResult;
	
	    static createFrom(source: any = {}) {
	        return new AuthResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], bootstrap.AuthResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DepartmentListResponse {
	    success: boolean;
	    message: string;
	    data: employees.DepartmentOption[];
	
	    static createFrom(source: any = {}) {
	        return new DepartmentListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], employees.DepartmentOption);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EmployeeListResponse {
	    success: boolean;
	    message: string;
	    data: employees.EmployeeListResult;
	
	    static createFrom(source: any = {}) {
	        return new EmployeeListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], employees.EmployeeListResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EmployeeResponse {
	    success: boolean;
	    message: string;
	    data: employees.EmployeeView;
	
	    static createFrom(source: any = {}) {
	        return new EmployeeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], employees.EmployeeView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveBalanceResponse {
	    success: boolean;
	    message: string;
	    data: leave.BalanceSummary;
	
	    static createFrom(source: any = {}) {
	        return new LeaveBalanceResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.BalanceSummary);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveRequestListResponse {
	    success: boolean;
	    message: string;
	    data: leave.RequestList;
	
	    static createFrom(source: any = {}) {
	        return new LeaveRequestListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.RequestList);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveRequestResponse {
	    success: boolean;
	    message: string;
	    data: leave.LeaveRequest;
	
	    static createFrom(source: any = {}) {
	        return new LeaveRequestResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.LeaveRequest);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveTypeListResponse {
	    success: boolean;
	    message: string;
	    data: leave.LeaveType[];
	
	    static createFrom(source: any = {}) {
	        return new LeaveTypeListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.LeaveType);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LeaveTypeResponse {
	    success: boolean;
	    message: string;
	    data: leave.LeaveType;
	
	    static createFrom(source: any = {}) {
	        return new LeaveTypeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.LeaveType);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LockedDateListResponse {
	    success: boolean;
	    message: string;
	    data: leave.LockedDate[];
	
	    static createFrom(source: any = {}) {
	        return new LockedDateListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.LockedDate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LockedDateResponse {
	    success: boolean;
	    message: string;
	    data: leave.LockedDate;
	
	    static createFrom(source: any = {}) {
	        return new LockedDateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], leave.LockedDate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MeResponse {
	    success: boolean;
	    message: string;
	    data: bootstrap.AuthUser;
	
	    static createFrom(source: any = {}) {
	        return new MeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], bootstrap.AuthUser);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PayrollBatchDetailResponse {
	    success: boolean;
	    message: string;
	    data: payroll.BatchDetail;
	
	    static createFrom(source: any = {}) {
	        return new PayrollBatchDetailResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], payroll.BatchDetail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PayrollBatchListResponse {
	    success: boolean;
	    message: string;
	    data: payroll.BatchListResult;
	
	    static createFrom(source: any = {}) {
	        return new PayrollBatchListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], payroll.BatchListResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PayrollBatchResponse {
	    success: boolean;
	    message: string;
	    data: payroll.Batch;
	
	    static createFrom(source: any = {}) {
	        return new PayrollBatchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], payroll.Batch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PayrollCSVResponse {
	    success: boolean;
	    message: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new PayrollCSVResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}
	export class PayrollEntryResponse {
	    success: boolean;
	    message: string;
	    data: payroll.Entry;
	
	    static createFrom(source: any = {}) {
	        return new PayrollEntryResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], payroll.Entry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace payroll {
	
	export class Batch {
	    id: number;
	    month: string;
	    status: string;
	    created_by: number;
	    // Go type: time
	    created_at: any;
	    approved_by?: number;
	    // Go type: time
	    approved_at?: any;
	    // Go type: time
	    locked_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new Batch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.month = source["month"];
	        this.status = source["status"];
	        this.created_by = source["created_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.approved_by = source["approved_by"];
	        this.approved_at = this.convertValues(source["approved_at"], null);
	        this.locked_at = this.convertValues(source["locked_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Entry {
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
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.batch_id = source["batch_id"];
	        this.employee_id = source["employee_id"];
	        this.employee_name = source["employee_name"];
	        this.base_salary = source["base_salary"];
	        this.allowances_total = source["allowances_total"];
	        this.deductions_total = source["deductions_total"];
	        this.tax_total = source["tax_total"];
	        this.gross_pay = source["gross_pay"];
	        this.net_pay = source["net_pay"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BatchDetail {
	    batch: Batch;
	    entries: Entry[];
	
	    static createFrom(source: any = {}) {
	        return new BatchDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.batch = this.convertValues(source["batch"], Batch);
	        this.entries = this.convertValues(source["entries"], Entry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BatchFilter {
	    month: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.status = source["status"];
	    }
	}
	export class BatchListResult {
	    items: Batch[];
	
	    static createFrom(source: any = {}) {
	        return new BatchListResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], Batch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateBatchInput {
	    month: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateBatchInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	    }
	}
	
	export class UpdateEntryAmountsInput {
	    allowances_total: number;
	    deductions_total: number;
	    tax_total: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateEntryAmountsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.allowances_total = source["allowances_total"];
	        this.deductions_total = source["deductions_total"];
	        this.tax_total = source["tax_total"];
	    }
	}

}

