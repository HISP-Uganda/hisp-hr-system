export type UserView = {
  id: number;
  username: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string | null;
};

export type UserListQuery = {
  page: number;
  page_size: number;
  q: string;
};

export type CreateUserInput = {
  username: string;
  password: string;
  role: string;
};

export type UpdateUserInput = {
  username: string;
  role: string;
};

export type ResetPasswordInput = {
  new_password: string;
};

export type UserStatusInput = {
  is_active: boolean;
};

export type UserResponse = {
  success: boolean;
  message: string;
  data: UserView;
};

export type UserListResponse = {
  success: boolean;
  message: string;
  data: {
    items: UserView[];
    total: number;
    page: number;
    page_size: number;
  };
};
