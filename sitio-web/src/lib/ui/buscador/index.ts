import Dialog from "./Dialog.svelte";
import state from "./state.svelte";
import Trigger from "./Trigger.svelte";

const Buscador = {
  Dialog,
  Trigger,
  state
};

export default Buscador;

export type { MateriaBuscador } from "./state.svelte";
