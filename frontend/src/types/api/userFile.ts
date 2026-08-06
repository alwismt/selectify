export type UserFile = {
  id: string;
  user_id: number;
  content_type: string;
};

export function isUserFile(value: unknown): value is UserFile {
  if (!value || typeof value !== "object") return false;
  const record = value as Record<string, unknown>;
  return typeof record.id === "string" && record.id.length > 0;
}
