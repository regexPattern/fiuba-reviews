import type { LayoutLoad } from "./$types";
import { browser } from "$app/environment";
import { posthogInit } from "$lib/posthog";

export const load: LayoutLoad = async ({ data }) => {
  if (browser) {
    posthogInit();
  }

  return data;
};
