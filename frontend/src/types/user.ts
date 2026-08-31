export type UserRole = {
  role: string;
  merchant_role?: string;
  merchant_id?: number;
};

export type User = {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
  status: string;
  created_at: string;
  user_role?: UserRole;
};
