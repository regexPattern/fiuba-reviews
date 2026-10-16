<svelte:options css="injected" />

<script lang="ts">
  interface Props {
    materia: { codigo: string; nombre: string };
    catedra: { nombre: string };
  }

  let { materia, catedra }: Props = $props();
</script>

<div
  style="display: flex; flex-direction: column; justify-content: space-between; width: 100%; height: 100%; background-color: #342F37; padding: 60px; font-family: 'Source Serif 4';"
>
  <div style="color: #f4dbd8; font-size: 32px; letter-spacing: -0.025em;">
    {materia.codigo} • {materia.nombre}
  </div>

  <div
    style="color: #ececec; font-size: 64px; letter-spacing: -0.025em; overflow: hidden; max-height: 240px; line-height: 1.25;"
  >
    Cátedra {catedra.nombre}
  </div>

  <div
    style="display: flex; justify-content: center; align-items: center; font-size: 32px; letter-spacing: -0.025em; gap: 8px; font-weight: 600;"
  >
    <span style="color: #4facd4;">FIUBA</span>
    <span style="color: #ececec;">Reviews</span>
  </div>
</div>
