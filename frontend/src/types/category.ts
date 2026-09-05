/** Home carousel category tile (static marketing data). */
export type Category = {
  title: string;
  id: number;
  img: string;
};

/** Hierarchical category from GET /categories. */
export type ApiCategory = {
  categoryId: number;
  name: string;
  slug: string;
  parentId?: number;
  isActive: boolean;
  children?: ApiCategory[];
};
