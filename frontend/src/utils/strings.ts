export const REQUIRED = 'Obrigatório.';

/*
 * Converts a string to a bool.
 *
 * This conversion will:
 *
 *  - match 'true', 'on', or '1' as true.
 *  - ignore all white-space padding
 *  - ignore capitalization (case).
 *
 * '  tRue  ','ON', and '1   ' will all evaluate as true.
 *
 */
export function strToBool(s: string): boolean {
  const regex = /^\s*(true|1|on)\s*$/i;

  return regex.test(s);
}
