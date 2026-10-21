"use client";

import { useEffect, useState } from "react";

export function useDebouncedValue<T>(value: T, delay = 80): T {
  const [debounced, setDebounced] = useState(value);

  useEffect(() => {
    const timeout = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(timeout);
  }, [value, delay]);

  return debounced;
}
