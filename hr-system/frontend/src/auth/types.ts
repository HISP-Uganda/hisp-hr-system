export type AuthUser = {
  id: number;
  username: string;
  role: string;
  is_active: boolean;
};

export type AuthResult = {
  access_token: string;
  refresh_token: string;
  access_expiry: string;
  user: AuthUser;
};

export type AuthApiResponse = {
  success: boolean;
  message: string;
  data: AuthResult;
};

export type MeApiResponse = {
  success: boolean;
  message: string;
  data: AuthUser;
};
