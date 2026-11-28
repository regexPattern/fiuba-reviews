import { useEffect, useState } from "react";

/** Hook que devuelve `value` recién cuando pasan `delayMs` sin que cambie. */
export function useDebouncedValue<T>(value: T, delayMs = 80): T {
  const [debounced, setDebounced] = useState(value);

  useEffect(() => {
    const timeout = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(timeout);
  }, [value, delayMs]);

  return debounced;
}
