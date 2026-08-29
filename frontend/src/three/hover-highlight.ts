import { Mesh, Raycaster, Vector2, type Camera } from 'three';

import { cellIdentityOf, setCellHighlighted } from './cube-geometry';

export interface HoverHighlight {
  setTargets(targets: Mesh[]): void;
  update(): void;
  clear(): void;
  dispose(): void;
}

export function createHoverHighlight(
  canvas: HTMLCanvasElement,
  camera: Camera,
): HoverHighlight {
  const raycaster = new Raycaster();
  const pointer = new Vector2();

  let targets: Mesh[] = [];
  let pointerInside = false;
  let highlighted: Mesh | null = null;

  const clearHighlight = () => {
    if (highlighted) {
      setCellHighlighted(highlighted, false);
      highlighted = null;
    }
  };

  const onPointerMove = (event: PointerEvent) => {
    const rect = canvas.getBoundingClientRect();
    pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
    pointer.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
    pointerInside = true;
  };

  const onPointerLeave = () => {
    pointerInside = false;
    clearHighlight();
  };

  canvas.addEventListener('pointermove', onPointerMove);
  canvas.addEventListener('pointerleave', onPointerLeave);

  return {
    setTargets(next) {
      clearHighlight();
      targets = next;
    },
    update() {
      if (!pointerInside || targets.length === 0) {
        return;
      }

      raycaster.setFromCamera(pointer, camera);
      const hits = raycaster.intersectObjects(targets, false);
      const nearest = hits[0]?.object;
      const candidate =
        nearest instanceof Mesh && cellIdentityOf(nearest) ? nearest : null;

      if (candidate === highlighted) {
        return;
      }

      clearHighlight();
      if (candidate) {
        setCellHighlighted(candidate, true);
        highlighted = candidate;
      }
    },
    clear: clearHighlight,
    dispose() {
      clearHighlight();
      canvas.removeEventListener('pointermove', onPointerMove);
      canvas.removeEventListener('pointerleave', onPointerLeave);
    },
  };
}
