import type { ValidationError } from '@nestjs/common';

export function formatValidationErrors(
  errors: ValidationError[],
  parentPath?: string,
): Record<string, string[]> {
  return errors.reduce<Record<string, string[]>>((accumulator, error) => {
    const currentPath = parentPath ? `${parentPath}.${error.property}` : error.property;
    if (error.constraints) {
      accumulator[currentPath] = Object.values(error.constraints);
    }
    if (error.children?.length) {
      Object.assign(accumulator, formatValidationErrors(error.children, currentPath));
    }
    return accumulator;
  }, {});
}

export function isValidPhone(value?: string | null): boolean {
  if (!value) {
    return true;
  }

  const normalized = value.replace(/[\s()-]/g, '');
  const ng = /^(?:\+?234|0)[789][01]\d{8}$/;
  const us = /^(?:\+?1)?\d{10}$/;
  return ng.test(normalized) || us.test(normalized);
}
