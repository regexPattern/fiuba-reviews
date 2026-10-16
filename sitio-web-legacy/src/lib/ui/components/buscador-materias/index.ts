import Dialog from "./Dialog.svelte";
import state from "./state.svelte";
import Trigger from "./Trigger.svelte";

const BuscadorMaterias = {
  Dialog,
  Trigger,
  state
};

export default BuscadorMaterias;

export type { MateriaBuscador } from "./state.svelte";
