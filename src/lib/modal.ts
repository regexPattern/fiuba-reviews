/**
 * Abre un `<dialog>` como modal compensando la barra de scroll del documento: al abrirlo,
 * `html:has(dialog:modal)` oculta el overflow y la barra desaparece. Para que el layout no salte,
 * se mide su ancho antes de abrir y `globals.css` lo usa como padding mientras el modal está abierto.
 */
export function abrirModal(dialog: HTMLDialogElement) {
  const html = document.documentElement;
  // Con barras overlay (o sin scroll) da 0 y no se compensa nada.
  html.style.setProperty("--ancho-scrollbar", `${window.innerWidth - html.clientWidth}px`);
  dialog.showModal();
}
