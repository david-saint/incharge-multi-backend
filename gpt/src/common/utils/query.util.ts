export function parseCsvList(value?: string | string[]): string[] {
  if (!value) {
    return [];
  }
  if (Array.isArray(value)) {
    return value.flatMap((item) => parseCsvList(item));
  }
  return value
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean);
}

export function parseBooleanQuery(value?: string | boolean): boolean {
  if (typeof value === 'boolean') {
    return value;
  }
  return ['1', 'true', 'yes', 'on'].includes((value ?? '').toLowerCase());
}

export function toNumberOrNull(value?: string | number | null): number | null {
  if (value === null || value === undefined || value === '') {
    return null;
  }
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : null;
}

export type SortDirective = { column: string; direction: 'ASC' | 'DESC' };

export function parseSortDirectives(
  value: string | undefined,
  allowedColumns: string[],
): SortDirective[] {
  return parseCsvList(value)
    .map((directive) => {
      const [column, direction = 'asc'] = directive.split('|');
      return {
        column,
        direction: direction.toUpperCase() === 'DESC' ? 'DESC' : 'ASC',
      } as SortDirective;
    })
    .filter((directive) => allowedColumns.includes(directive.column));
}

export function buildPaginationMeta(
  basePath: string,
  page: number,
  perPage: number,
  total: number,
) {
  const lastPage = Math.max(1, Math.ceil(total / perPage));
  const from = total === 0 ? null : (page - 1) * perPage + 1;
  const to = total === 0 ? null : Math.min(page * perPage, total);
  const query = (targetPage: number) => `${basePath}?page=${targetPage}&per_page=${perPage}`;
  return {
    current_page: page,
    per_page: perPage,
    total,
    last_page: lastPage,
    from,
    to,
    first_page_url: query(1),
    last_page_url: query(lastPage),
    next_page_url: page < lastPage ? query(page + 1) : null,
    prev_page_url: page > 1 ? query(page - 1) : null,
    path: basePath,
  };
}
