export type UserAddress = {
  id: number;
  user_id: number;
  label?: string;
  phone?: string;
  line1: string;
  line2?: string;
  city: string;
  region?: string;
  postal_code: string;
  country_code: string;
  is_default: boolean;
  created_at: string;
  updated_at: string;
};
