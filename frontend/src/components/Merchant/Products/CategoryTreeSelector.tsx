"use client";

import {
  useEffect,
  useMemo,
  useRef,
  useState,
  type KeyboardEvent,
  type MouseEvent,
} from "react";
import type { ApiCategory } from "@/types/category";

type CategoryTreeSelectorProps = {
  categories: ApiCategory[];
  selectedIds: number[];
  onChange: (ids: number[]) => void;
  disabled?: boolean;
  error?: string | null;
};

const depthPadding = [
  "pl-2",
  "pl-6",
  "pl-10",
  "pl-14",
  "pl-16",
  "pl-[4.5rem]",
] as const;

const triggerClassName =
  "w-full rounded-md border border-gray-3 bg-gray-1 py-2.5 px-5 text-left outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20 disabled:opacity-50 disabled:cursor-not-allowed";

const searchInputClassName =
  "rounded-md border border-gray-3 bg-white placeholder:text-dark-5 w-full py-2 px-3 text-custom-sm outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20";

function collectCategoryMap(
  categories: ApiCategory[],
  map = new Map<number, ApiCategory>()
): Map<number, ApiCategory> {
  for (const category of categories) {
    map.set(category.categoryId, category);
    if (category.children?.length) {
      collectCategoryMap(category.children, map);
    }
  }
  return map;
}

/** Ancestor ids that must be expanded to reveal currently selected categories. */
function collectAncestorIds(
  categories: ApiCategory[],
  selectedIds: number[]
): Set<number> {
  const selected = new Set(selectedIds);
  const ancestors = new Set<number>();

  function walk(nodes: ApiCategory[], path: number[]): boolean {
    let foundInBranch = false;
    for (const node of nodes) {
      const childMatch = node.children?.length
        ? walk(node.children, [...path, node.categoryId])
        : false;
      const selfMatch = selected.has(node.categoryId);
      if (selfMatch || childMatch) {
        foundInBranch = true;
        for (const id of path) {
          ancestors.add(id);
        }
        if (childMatch) {
          ancestors.add(node.categoryId);
        }
      }
    }
    return foundInBranch;
  }

  walk(categories, []);
  return ancestors;
}

/**
 * Prune tree to nodes that match the query or have a matching descendant.
 * Preserves hierarchy (not a flat list).
 */
function filterCategoryTree(
  categories: ApiCategory[],
  query: string
): ApiCategory[] {
  const normalized = query.trim().toLowerCase();
  if (!normalized) return categories;

  const result: ApiCategory[] = [];
  for (const category of categories) {
    const nameMatches = category.name.toLowerCase().includes(normalized);
    const filteredChildren = category.children?.length
      ? filterCategoryTree(category.children, query)
      : [];
    if (nameMatches || filteredChildren.length > 0) {
      result.push({
        ...category,
        children:
          filteredChildren.length > 0 ? filteredChildren : undefined,
      });
    }
  }
  return result;
}

/** Expand parents that have matching descendants so search hits are visible. */
function collectExpandIdsForFilteredTree(
  categories: ApiCategory[]
): Set<number> {
  const ids = new Set<number>();

  function walk(nodes: ApiCategory[]): void {
    for (const node of nodes) {
      if (node.children?.length) {
        ids.add(node.categoryId);
        walk(node.children);
      }
    }
  }

  walk(categories);
  return ids;
}

function CategoryTreeNode({
  category,
  depth,
  selectedIds,
  expandedIds,
  onToggleSelect,
  onToggleExpand,
  disabled,
}: {
  category: ApiCategory;
  depth: number;
  selectedIds: number[];
  expandedIds: Set<number>;
  onToggleSelect: (categoryId: number, checked: boolean) => void;
  onToggleExpand: (categoryId: number) => void;
  disabled?: boolean;
}) {
  const isInactive = !category.isActive;
  const isChecked = selectedIds.includes(category.categoryId);
  const canSelect = !isInactive && !disabled;
  const hasChildren = Boolean(category.children?.length);
  const isExpanded = expandedIds.has(category.categoryId);
  const padClass =
    depthPadding[Math.min(depth, depthPadding.length - 1)] ?? "pl-[4.5rem]";

  return (
    <li>
      <div
        className={`flex items-center gap-1.5 py-1.5 text-custom-sm ${padClass} ${
          isInactive ? "text-dark-4 opacity-60" : "text-dark"
        }`}
      >
        {hasChildren ? (
          <button
            type="button"
            aria-label={isExpanded ? "Collapse category" : "Expand category"}
            aria-expanded={isExpanded}
            onClick={(e: MouseEvent) => {
              e.preventDefault();
              e.stopPropagation();
              onToggleExpand(category.categoryId);
            }}
            className="inline-flex h-5 w-5 shrink-0 items-center justify-center text-dark-4 hover:text-dark"
          >
            <span className="text-xs leading-none" aria-hidden="true">
              {isExpanded ? "▼" : "▶"}
            </span>
          </button>
        ) : (
          <span className="inline-block h-5 w-5 shrink-0" aria-hidden="true" />
        )}

        <label
          className={`inline-flex min-w-0 items-center gap-2 ${
            canSelect ? "cursor-pointer" : "cursor-not-allowed"
          }`}
        >
          <input
            type="checkbox"
            checked={isChecked && !isInactive}
            disabled={!canSelect}
            onChange={(e) =>
              onToggleSelect(category.categoryId, e.target.checked)
            }
            onClick={(e) => e.stopPropagation()}
          />
          <span className="truncate">{category.name}</span>
          {isInactive ? (
            <span className="shrink-0 text-custom-xs text-dark-4">
              (inactive)
            </span>
          ) : null}
        </label>
      </div>

      {hasChildren && isExpanded ? (
        <ul className="list-none">
          {category.children!.map((child) => (
            <CategoryTreeNode
              key={child.categoryId}
              category={child}
              depth={depth + 1}
              selectedIds={selectedIds}
              expandedIds={expandedIds}
              onToggleSelect={onToggleSelect}
              onToggleExpand={onToggleExpand}
              disabled={disabled}
            />
          ))}
        </ul>
      ) : null}
    </li>
  );
}

export default function CategoryTreeSelector({
  categories,
  selectedIds,
  onChange,
  disabled = false,
  error = null,
}: CategoryTreeSelectorProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [expandedIds, setExpandedIds] = useState<Set<number>>(() => new Set());

  const categoryMap = useMemo(
    () => collectCategoryMap(categories),
    [categories]
  );

  const filteredCategories = useMemo(
    () => filterCategoryTree(categories, searchQuery),
    [categories, searchQuery]
  );

  const selectedCategories = useMemo(() => {
    return selectedIds
      .map((id) => categoryMap.get(id))
      .filter((c): c is ApiCategory => c != null);
  }, [selectedIds, categoryMap]);

  useEffect(() => {
    if (!isOpen) return;

    const handlePointerDown = (event: Event) => {
      const target = event.target as Node | null;
      if (
        containerRef.current &&
        target &&
        !containerRef.current.contains(target)
      ) {
        setIsOpen(false);
      }
    };

    const handleKeyDown = (event: Event) => {
      if ((event as globalThis.KeyboardEvent).key === "Escape") {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("mousedown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) return;
    const frame = requestAnimationFrame(() => {
      searchInputRef.current?.focus();
    });
    return () => cancelAnimationFrame(frame);
  }, [isOpen]);

  // Sync expand state when opening or when the search query changes — not on every selection.
  useEffect(() => {
    if (!isOpen) return;
    if (searchQuery.trim()) {
      setExpandedIds(collectExpandIdsForFilteredTree(filteredCategories));
      return;
    }
    setExpandedIds(collectAncestorIds(categories, selectedIds));
    // selectedIds intentionally omitted: selecting a checkbox must not collapse the tree
    // eslint-disable-next-line react-hooks/exhaustive-deps -- open/search driven
  }, [isOpen, searchQuery, filteredCategories, categories]);

  const handleToggleSelect = (categoryId: number, checked: boolean) => {
    if (checked) {
      if (selectedIds.includes(categoryId)) return;
      onChange([...selectedIds, categoryId]);
      return;
    }
    onChange(selectedIds.filter((id) => id !== categoryId));
  };

  const handleToggleExpand = (categoryId: number) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(categoryId)) {
        next.delete(categoryId);
      } else {
        next.add(categoryId);
      }
      return next;
    });
  };

  const openDropdown = () => {
    if (disabled) return;
    setSearchQuery("");
    setExpandedIds(collectAncestorIds(categories, selectedIds));
    setIsOpen(true);
  };

  const handleTriggerClick = () => {
    if (disabled) return;
    if (isOpen) {
      setIsOpen(false);
      return;
    }
    openDropdown();
  };

  const handleTriggerKeyDown = (event: KeyboardEvent<HTMLButtonElement>) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      handleTriggerClick();
    } else if (event.key === "ArrowDown" && !isOpen) {
      event.preventDefault();
      openDropdown();
    }
  };

  if (categories.length === 0) {
    return (
      <div>
        <p className="text-custom-sm text-dark-4">
          No categories available. Please try again later.
        </p>
        {error ? (
          <p className="mt-2 text-custom-sm text-red">{error}</p>
        ) : null}
      </div>
    );
  }

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        className={`${triggerClassName} flex items-center justify-between gap-3 ${
          isOpen ? "ring-2 ring-blue/20" : ""
        }`}
        onClick={handleTriggerClick}
        onKeyDown={handleTriggerKeyDown}
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={isOpen}
      >
        <span
          className={selectedIds.length > 0 ? "text-dark" : "text-dark-5"}
        >
          {selectedIds.length > 0
            ? `${selectedIds.length} categor${
                selectedIds.length === 1 ? "y" : "ies"
              } selected`
            : "Select categories..."}
        </span>
        <span
          className={`text-dark-4 transition-transform ${
            isOpen ? "rotate-180" : ""
          }`}
          aria-hidden="true"
        >
          ▾
        </span>
      </button>

      {selectedCategories.length > 0 ? (
        <div className="mt-2.5 flex flex-wrap gap-2">
          {selectedCategories.map((category) => (
            <span
              key={category.categoryId}
              className="inline-flex items-center gap-1.5 rounded-md bg-gray-2 px-2.5 py-1 text-custom-sm text-dark"
            >
              {category.name}
              <button
                type="button"
                aria-label={`Remove ${category.name}`}
                disabled={disabled}
                onClick={() => handleToggleSelect(category.categoryId, false)}
                className="text-dark-4 hover:text-red disabled:opacity-50"
              >
                ×
              </button>
            </span>
          ))}
        </div>
      ) : null}

      {isOpen ? (
        <div className="absolute left-0 right-0 z-20 mt-1.5 rounded-md border border-gray-3 bg-white shadow-1">
          <div className="border-b border-gray-3 p-3">
            <input
              ref={searchInputRef}
              type="search"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search categories..."
              className={searchInputClassName}
              disabled={disabled}
              autoComplete="off"
            />
          </div>
          <div className="max-h-[320px] overflow-y-auto py-2 px-2">
            {filteredCategories.length === 0 ? (
              <p className="px-3 py-2 text-custom-sm text-dark-4">
                No categories match your search.
              </p>
            ) : (
              <ul className="list-none" role="listbox" aria-multiselectable>
                {filteredCategories.map((category) => (
                  <CategoryTreeNode
                    key={category.categoryId}
                    category={category}
                    depth={0}
                    selectedIds={selectedIds}
                    expandedIds={expandedIds}
                    onToggleSelect={handleToggleSelect}
                    onToggleExpand={handleToggleExpand}
                    disabled={disabled}
                  />
                ))}
              </ul>
            )}
          </div>
        </div>
      ) : null}

      {error ? (
        <p className="mt-2 text-custom-sm text-red">{error}</p>
      ) : null}
    </div>
  );
}
